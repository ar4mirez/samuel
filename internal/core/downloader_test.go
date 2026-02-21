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

func TestExtractTarGz_InvalidGzip(t *testing.T) {
	buf := bytes.NewBufferString("this is not gzip data")

	err := extractTarGz(buf, t.TempDir())
	if err == nil {
		t.Fatal("expected error for invalid gzip data, got nil")
	}
	if !contains(err.Error(), "gzip") {
		t.Errorf("expected gzip-related error, got: %v", err)
	}
}

func TestExtractTarGz_EmptyArchive(t *testing.T) {
	dest := t.TempDir()

	// Create a valid tar.gz with no entries
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)
	tw.Close()
	gw.Close()

	err := extractTarGz(&buf, dest)
	if err != nil {
		t.Fatalf("unexpected error for empty archive: %v", err)
	}
}

func TestExtractTarGz_NestedDirsWithoutExplicitEntries(t *testing.T) {
	dest := t.TempDir()

	// Create archive with a deeply nested file but no explicit dir entries
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)

	content := "deep file content"
	if err := tw.WriteHeader(&tar.Header{
		Name:     "repo/a/b/c/deep.txt",
		Typeflag: tar.TypeReg,
		Mode:     0644,
		Size:     int64(len(content)),
	}); err != nil {
		t.Fatalf("failed to write header: %v", err)
	}
	if _, err := tw.Write([]byte(content)); err != nil {
		t.Fatalf("failed to write content: %v", err)
	}

	tw.Close()
	gw.Close()

	err := extractTarGz(&buf, dest)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := os.ReadFile(filepath.Join(dest, "repo", "a", "b", "c", "deep.txt"))
	if err != nil {
		t.Fatalf("failed to read deep file: %v", err)
	}
	if string(got) != content {
		t.Errorf("expected %q, got %q", content, string(got))
	}
}

func TestExtractTarGz_FilePermissions(t *testing.T) {
	dest := t.TempDir()

	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)

	content := "#!/bin/sh\necho hello"
	if err := tw.WriteHeader(&tar.Header{
		Name:     "repo/script.sh",
		Typeflag: tar.TypeReg,
		Mode:     0755,
		Size:     int64(len(content)),
	}); err != nil {
		t.Fatalf("failed to write header: %v", err)
	}
	if _, err := tw.Write([]byte(content)); err != nil {
		t.Fatalf("failed to write content: %v", err)
	}

	tw.Close()
	gw.Close()

	err := extractTarGz(&buf, dest)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	info, err := os.Stat(filepath.Join(dest, "repo", "script.sh"))
	if err != nil {
		t.Fatalf("failed to stat file: %v", err)
	}
	if info.Mode().Perm()&0111 == 0 {
		t.Errorf("expected executable permission, got %v", info.Mode().Perm())
	}
}

func TestExtractTarGz_DirectoryTraversal(t *testing.T) {
	dest := t.TempDir()

	// Create archive with a directory entry that traverses
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)

	if err := tw.WriteHeader(&tar.Header{
		Name:     "../../evil-dir/",
		Typeflag: tar.TypeDir,
		Mode:     0755,
	}); err != nil {
		t.Fatalf("failed to write header: %v", err)
	}

	tw.Close()
	gw.Close()

	err := extractTarGz(&buf, dest)
	if err == nil {
		t.Fatal("expected error for directory path traversal, got nil")
	}
	if !contains(err.Error(), "invalid file path") {
		t.Errorf("expected 'invalid file path' error, got: %v", err)
	}
}

func TestExtractTarGz_FileSizeLimit(t *testing.T) {
	dest := t.TempDir()

	// Save and restore the original limit for test isolation
	origLimit := MaxExtractedFileSize
	MaxExtractedFileSize = 1024 // 1KB limit for testing
	defer func() { MaxExtractedFileSize = origLimit }()

	// Create archive with a file exceeding the limit
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)

	// Write a file larger than the 1KB test limit
	oversized := bytes.Repeat([]byte("x"), 2048)
	if err := tw.WriteHeader(&tar.Header{
		Name:     "repo/large-file.bin",
		Typeflag: tar.TypeReg,
		Mode:     0644,
		Size:     int64(len(oversized)),
	}); err != nil {
		t.Fatalf("failed to write header: %v", err)
	}
	if _, err := tw.Write(oversized); err != nil {
		t.Fatalf("failed to write content: %v", err)
	}

	tw.Close()
	gw.Close()

	err := extractTarGz(&buf, dest)
	if err == nil {
		t.Fatal("expected error for oversized file, got nil")
	}
	if !strings.Contains(err.Error(), "exceeds maximum size") {
		t.Errorf("expected 'exceeds maximum size' error, got: %v", err)
	}
}

func TestExtractTarGz_FileSizeAtLimit(t *testing.T) {
	dest := t.TempDir()

	origLimit := MaxExtractedFileSize
	MaxExtractedFileSize = 1024
	defer func() { MaxExtractedFileSize = origLimit }()

	// File exactly at the limit should succeed
	exactSize := bytes.Repeat([]byte("y"), 1024)
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)

	if err := tw.WriteHeader(&tar.Header{
		Name:     "repo/exact-file.bin",
		Typeflag: tar.TypeReg,
		Mode:     0644,
		Size:     int64(len(exactSize)),
	}); err != nil {
		t.Fatalf("failed to write header: %v", err)
	}
	if _, err := tw.Write(exactSize); err != nil {
		t.Fatalf("failed to write content: %v", err)
	}

	tw.Close()
	gw.Close()

	err := extractTarGz(&buf, dest)
	if err != nil {
		t.Fatalf("file at exact size limit should succeed, got: %v", err)
	}

	// Verify file was extracted
	data, err := os.ReadFile(filepath.Join(dest, "repo", "exact-file.bin"))
	if err != nil {
		t.Fatalf("failed to read extracted file: %v", err)
	}
	if len(data) != 1024 {
		t.Errorf("expected 1024 bytes, got %d", len(data))
	}
}

func TestCopyFile(t *testing.T) {
	t.Run("copies content and permissions", func(t *testing.T) {
		srcDir := t.TempDir()
		dstDir := t.TempDir()

		srcPath := filepath.Join(srcDir, "source.txt")
		dstPath := filepath.Join(dstDir, "dest.txt")

		content := "hello, world!"
		if err := os.WriteFile(srcPath, []byte(content), 0644); err != nil {
			t.Fatalf("failed to create source file: %v", err)
		}

		if err := copyFile(srcPath, dstPath); err != nil {
			t.Fatalf("copyFile failed: %v", err)
		}

		got, err := os.ReadFile(dstPath)
		if err != nil {
			t.Fatalf("failed to read destination: %v", err)
		}
		if string(got) != content {
			t.Errorf("expected %q, got %q", content, string(got))
		}
	})

	t.Run("source not found", func(t *testing.T) {
		err := copyFile("/nonexistent/file.txt", filepath.Join(t.TempDir(), "out.txt"))
		if err == nil {
			t.Fatal("expected error for missing source, got nil")
		}
	})

	t.Run("destination directory not found", func(t *testing.T) {
		srcDir := t.TempDir()
		srcPath := filepath.Join(srcDir, "src.txt")
		if err := os.WriteFile(srcPath, []byte("data"), 0644); err != nil {
			t.Fatalf("failed to create source: %v", err)
		}

		err := copyFile(srcPath, "/nonexistent/dir/out.txt")
		if err == nil {
			t.Fatal("expected error for missing dest dir, got nil")
		}
	})
}

func TestCopyDir(t *testing.T) {
	t.Run("copies directory tree", func(t *testing.T) {
		src := t.TempDir()
		dst := filepath.Join(t.TempDir(), "dest")

		// Create source structure: src/a.txt, src/sub/b.txt
		if err := os.WriteFile(filepath.Join(src, "a.txt"), []byte("file-a"), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.MkdirAll(filepath.Join(src, "sub"), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(src, "sub", "b.txt"), []byte("file-b"), 0644); err != nil {
			t.Fatal(err)
		}

		if err := copyDir(src, dst); err != nil {
			t.Fatalf("copyDir failed: %v", err)
		}

		// Verify copied content
		gotA, err := os.ReadFile(filepath.Join(dst, "a.txt"))
		if err != nil {
			t.Fatalf("failed to read a.txt: %v", err)
		}
		if string(gotA) != "file-a" {
			t.Errorf("a.txt: expected 'file-a', got %q", string(gotA))
		}

		gotB, err := os.ReadFile(filepath.Join(dst, "sub", "b.txt"))
		if err != nil {
			t.Fatalf("failed to read sub/b.txt: %v", err)
		}
		if string(gotB) != "file-b" {
			t.Errorf("sub/b.txt: expected 'file-b', got %q", string(gotB))
		}
	})

	t.Run("copies empty directory", func(t *testing.T) {
		src := t.TempDir()
		dst := filepath.Join(t.TempDir(), "dest")

		if err := copyDir(src, dst); err != nil {
			t.Fatalf("copyDir failed on empty dir: %v", err)
		}

		info, err := os.Stat(dst)
		if err != nil {
			t.Fatalf("dest dir not created: %v", err)
		}
		if !info.IsDir() {
			t.Error("expected dest to be a directory")
		}
	})

	t.Run("source not found", func(t *testing.T) {
		err := copyDir("/nonexistent/source", filepath.Join(t.TempDir(), "dest"))
		if err == nil {
			t.Fatal("expected error for nonexistent source, got nil")
		}
	})
}

func TestClearCache(t *testing.T) {
	t.Run("clears populated cache", func(t *testing.T) {
		cacheDir := t.TempDir()
		d := &Downloader{cachePath: cacheDir}

		// Create some cached files
		if err := os.MkdirAll(filepath.Join(cacheDir, "samuel-v1.0.0"), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(cacheDir, "samuel-v1.0.0", "file.txt"), []byte("cached"), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(cacheDir, "other.txt"), []byte("other"), 0644); err != nil {
			t.Fatal(err)
		}

		if err := d.ClearCache(); err != nil {
			t.Fatalf("ClearCache failed: %v", err)
		}

		entries, err := os.ReadDir(cacheDir)
		if err != nil {
			t.Fatalf("failed to read cache dir: %v", err)
		}
		if len(entries) != 0 {
			t.Errorf("expected empty cache, got %d entries", len(entries))
		}
	})

	t.Run("nonexistent cache dir returns nil", func(t *testing.T) {
		d := &Downloader{cachePath: "/nonexistent/cache/path"}

		if err := d.ClearCache(); err != nil {
			t.Fatalf("expected nil for nonexistent cache, got: %v", err)
		}
	})

	t.Run("empty cache is no-op", func(t *testing.T) {
		cacheDir := t.TempDir()
		d := &Downloader{cachePath: cacheDir}

		if err := d.ClearCache(); err != nil {
			t.Fatalf("ClearCache failed on empty dir: %v", err)
		}
	})
}

func TestGetCacheSize(t *testing.T) {
	t.Run("calculates size of files", func(t *testing.T) {
		cacheDir := t.TempDir()
		d := &Downloader{cachePath: cacheDir}

		// Create files with known sizes
		data10 := bytes.Repeat([]byte("a"), 10)
		data20 := bytes.Repeat([]byte("b"), 20)

		if err := os.WriteFile(filepath.Join(cacheDir, "f1.txt"), data10, 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.MkdirAll(filepath.Join(cacheDir, "sub"), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(cacheDir, "sub", "f2.txt"), data20, 0644); err != nil {
			t.Fatal(err)
		}

		size, err := d.GetCacheSize()
		if err != nil {
			t.Fatalf("GetCacheSize failed: %v", err)
		}
		if size != 30 {
			t.Errorf("expected size 30, got %d", size)
		}
	})

	t.Run("empty cache returns zero", func(t *testing.T) {
		cacheDir := t.TempDir()
		d := &Downloader{cachePath: cacheDir}

		size, err := d.GetCacheSize()
		if err != nil {
			t.Fatalf("GetCacheSize failed: %v", err)
		}
		if size != 0 {
			t.Errorf("expected size 0, got %d", size)
		}
	})
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
