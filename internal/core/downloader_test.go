package core

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateSymlinkTarget(t *testing.T) {
	dest := "/tmp/extract"

	tests := []struct {
		name        string
		symlinkPath string
		linkTarget  string
		wantErr     bool
		errContains string
	}{
		{
			name:        "valid relative symlink within dest",
			symlinkPath: "/tmp/extract/repo/link",
			linkTarget:  "../otherfile",
			wantErr:     false,
		},
		{
			name:        "valid symlink to sibling file",
			symlinkPath: "/tmp/extract/repo/dir/link",
			linkTarget:  "target.txt",
			wantErr:     false,
		},
		{
			name:        "valid symlink to parent dir within dest",
			symlinkPath: "/tmp/extract/repo/dir/subdir/link",
			linkTarget:  "../../file.txt",
			wantErr:     false,
		},
		{
			name:        "absolute symlink target rejected",
			symlinkPath: "/tmp/extract/repo/link",
			linkTarget:  "/etc/passwd",
			wantErr:     true,
			errContains: "absolute path",
		},
		{
			name:        "relative symlink escaping dest",
			symlinkPath: "/tmp/extract/repo/link",
			linkTarget:  "../../../../../../etc/shadow",
			wantErr:     true,
			errContains: "escapes destination",
		},
		{
			name:        "relative symlink just barely escaping",
			symlinkPath: "/tmp/extract/link",
			linkTarget:  "../outside",
			wantErr:     true,
			errContains: "escapes destination",
		},
		{
			name:        "symlink targeting dest root is valid",
			symlinkPath: "/tmp/extract/repo/link",
			linkTarget:  "..",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateSymlinkTarget(dest, tt.symlinkPath, tt.linkTarget)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.errContains)
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("expected error containing %q, got %q", tt.errContains, err.Error())
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestExtractTarGz_SymlinkTraversal(t *testing.T) {
	dest := t.TempDir()

	// Build a tar.gz archive with a malicious symlink
	buf := createTarGzWithSymlink(t, "repo/evil-link", "../../etc/passwd")

	err := extractTarGz(buf, dest)
	if err == nil {
		t.Fatal("expected error for symlink traversal, got nil")
	}
	if !contains(err.Error(), "escapes destination") {
		t.Errorf("expected 'escapes destination' error, got: %v", err)
	}

	// Verify symlink was NOT created
	linkPath := filepath.Join(dest, "repo", "evil-link")
	if _, err := os.Lstat(linkPath); !os.IsNotExist(err) {
		t.Errorf("malicious symlink should not have been created")
	}
}

func TestExtractTarGz_AbsoluteSymlink(t *testing.T) {
	dest := t.TempDir()

	buf := createTarGzWithSymlink(t, "repo/abs-link", "/etc/passwd")

	err := extractTarGz(buf, dest)
	if err == nil {
		t.Fatal("expected error for absolute symlink target, got nil")
	}
	if !contains(err.Error(), "absolute path") {
		t.Errorf("expected 'absolute path' error, got: %v", err)
	}
}

func TestExtractTarGz_ValidSymlink(t *testing.T) {
	dest := t.TempDir()

	// Build a tar.gz with a valid symlink pointing to a sibling file
	buf := createTarGzWithFileAndSymlink(t,
		"repo/target.txt", "hello",
		"repo/link.txt", "target.txt",
	)

	err := extractTarGz(buf, dest)
	if err != nil {
		t.Fatalf("unexpected error for valid symlink: %v", err)
	}

	// Verify the symlink was created and points correctly
	linkPath := filepath.Join(dest, "repo", "link.txt")
	linkTarget, err := os.Readlink(linkPath)
	if err != nil {
		t.Fatalf("failed to read symlink: %v", err)
	}
	if linkTarget != "target.txt" {
		t.Errorf("expected symlink target 'target.txt', got %q", linkTarget)
	}
}

func TestExtractTarGz_BasicExtraction(t *testing.T) {
	dest := t.TempDir()

	buf := createTarGzWithFiles(t, map[string]string{
		"repo/README.md": "# Test",
		"repo/src/main.go": "package main",
	})

	err := extractTarGz(buf, dest)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify files were extracted
	content, err := os.ReadFile(filepath.Join(dest, "repo", "README.md"))
	if err != nil {
		t.Fatalf("failed to read extracted file: %v", err)
	}
	if string(content) != "# Test" {
		t.Errorf("expected '# Test', got %q", string(content))
	}
}

func TestExtractTarGz_PathTraversal(t *testing.T) {
	dest := t.TempDir()

	buf := createTarGzWithFiles(t, map[string]string{
		"../../etc/evil": "pwned",
	})

	err := extractTarGz(buf, dest)
	if err == nil {
		t.Fatal("expected error for path traversal, got nil")
	}
	if !contains(err.Error(), "invalid file path") {
		t.Errorf("expected 'invalid file path' error, got: %v", err)
	}
}

// Helper functions

func contains(s, substr string) bool {
	return bytes.Contains([]byte(s), []byte(substr))
}

func createTarGzWithSymlink(t *testing.T, name, target string) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)

	// Add a directory first
	dir := filepath.Dir(name)
	if dir != "." {
		if err := tw.WriteHeader(&tar.Header{
			Name:     dir + "/",
			Typeflag: tar.TypeDir,
			Mode:     0755,
		}); err != nil {
			t.Fatalf("failed to write dir header: %v", err)
		}
	}

	// Add the symlink
	if err := tw.WriteHeader(&tar.Header{
		Name:     name,
		Typeflag: tar.TypeSymlink,
		Linkname: target,
	}); err != nil {
		t.Fatalf("failed to write symlink header: %v", err)
	}

	tw.Close()
	gw.Close()
	return &buf
}

func createTarGzWithFileAndSymlink(t *testing.T, fileName, fileContent, linkName, linkTarget string) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)

	// Add directory
	dir := filepath.Dir(fileName)
	if err := tw.WriteHeader(&tar.Header{
		Name:     dir + "/",
		Typeflag: tar.TypeDir,
		Mode:     0755,
	}); err != nil {
		t.Fatalf("failed to write dir header: %v", err)
	}

	// Add the regular file
	if err := tw.WriteHeader(&tar.Header{
		Name:     fileName,
		Typeflag: tar.TypeReg,
		Mode:     0644,
		Size:     int64(len(fileContent)),
	}); err != nil {
		t.Fatalf("failed to write file header: %v", err)
	}
	if _, err := tw.Write([]byte(fileContent)); err != nil {
		t.Fatalf("failed to write file content: %v", err)
	}

	// Add the symlink
	if err := tw.WriteHeader(&tar.Header{
		Name:     linkName,
		Typeflag: tar.TypeSymlink,
		Linkname: linkTarget,
	}); err != nil {
		t.Fatalf("failed to write symlink header: %v", err)
	}

	tw.Close()
	gw.Close()
	return &buf
}

func createTarGzWithFiles(t *testing.T, files map[string]string) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)

	// Track directories already added
	dirs := make(map[string]bool)

	for name, content := range files {
		// Add parent directories
		dir := filepath.Dir(name)
		parts := strings.Split(dir, "/")
		for i := range parts {
			d := strings.Join(parts[:i+1], "/") + "/"
			if !dirs[d] && d != "./" {
				dirs[d] = true
				if err := tw.WriteHeader(&tar.Header{
					Name:     d,
					Typeflag: tar.TypeDir,
					Mode:     0755,
				}); err != nil {
					t.Fatalf("failed to write dir header: %v", err)
				}
			}
		}

		// Add the file
		if err := tw.WriteHeader(&tar.Header{
			Name:     name,
			Typeflag: tar.TypeReg,
			Mode:     0644,
			Size:     int64(len(content)),
		}); err != nil {
			t.Fatalf("failed to write file header: %v", err)
		}
		if _, err := tw.Write([]byte(content)); err != nil {
			t.Fatalf("failed to write file content: %v", err)
		}
	}

	tw.Close()
	gw.Close()
	return &buf
}
