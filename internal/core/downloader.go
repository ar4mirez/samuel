package core

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ar4mirez/samuel/internal/github"
)

// MaxExtractedFileSize is the maximum allowed size for a single file
// extracted from a tar archive (100 MB). Prevents decompression bombs.
var MaxExtractedFileSize int64 = 100 * 1024 * 1024

// Downloader handles downloading and extracting framework files
type Downloader struct {
	client    *github.Client
	cachePath string
}

// NewDownloader creates a new downloader
func NewDownloader() (*Downloader, error) {
	cachePath, err := EnsureCacheDir()
	if err != nil {
		return nil, err
	}

	return &Downloader{
		client:    github.NewClient(DefaultOwner, DefaultRepo),
		cachePath: cachePath,
	}, nil
}

// DownloadVersion downloads a specific version to the cache
// If version is "dev", downloads from main branch
func (d *Downloader) DownloadVersion(version string) (string, error) {
	// Check if already cached (skip cache for dev version)
	cacheDest := filepath.Join(d.cachePath, fmt.Sprintf("samuel-%s", version))
	if version != github.DevVersion {
		if _, err := os.Stat(cacheDest); err == nil {
			return cacheDest, nil
		}
	} else {
		// Clear dev cache to always get fresh copy
		if err := os.RemoveAll(cacheDest); err != nil {
			return "", fmt.Errorf("failed to clear dev cache: %w", err)
		}
	}

	// Download archive
	var reader io.ReadCloser
	var err error

	if version == github.DevVersion {
		reader, _, err = d.client.DownloadBranchArchive(github.DefaultBranch)
	} else {
		reader, _, err = d.client.DownloadArchive(version)
	}

	if err != nil {
		return "", err
	}
	defer reader.Close()

	// Create temp directory for extraction
	tempDir, err := os.MkdirTemp("", "samuel-download-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Extract archive
	if err := extractTarGz(reader, tempDir); err != nil {
		return "", fmt.Errorf("failed to extract archive: %w", err)
	}

	// Find the extracted directory (GitHub adds repo-version prefix)
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		return "", err
	}

	if len(entries) != 1 || !entries[0].IsDir() {
		return "", fmt.Errorf("unexpected archive structure")
	}

	extractedDir := filepath.Join(tempDir, entries[0].Name())

	// Move to cache
	if err := os.MkdirAll(filepath.Dir(cacheDest), 0755); err != nil {
		return "", err
	}

	if err := os.Rename(extractedDir, cacheDest); err != nil {
		// If rename fails (cross-device), copy instead
		if err := copyDir(extractedDir, cacheDest); err != nil {
			return "", fmt.Errorf("failed to cache download: %w", err)
		}
	}

	return cacheDest, nil
}

// GetLatestVersion fetches the latest version number
// Returns "dev" if no releases exist
func (d *Downloader) GetLatestVersion() (string, error) {
	version, _, err := d.client.GetLatestVersionOrBranch()
	return version, err
}

// DownloadFile downloads a single file from a version
func (d *Downloader) DownloadFile(version, path string) ([]byte, error) {
	return d.client.DownloadFile(version, path)
}

// CheckForUpdates checks if a newer version is available
func (d *Downloader) CheckForUpdates(currentVersion string) (*github.VersionInfo, error) {
	return d.client.CheckForUpdates(currentVersion)
}

// extractTarGz extracts a tar.gz archive to a destination directory
func extractTarGz(reader io.Reader, dest string) error {
	gzReader, err := gzip.NewReader(reader)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar entry: %w", err)
		}

		// Sanitize path to prevent directory traversal
		target := filepath.Join(dest, header.Name)
		if !strings.HasPrefix(target, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path: %s", header.Name)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}

		case tar.TypeReg:
			// Ensure parent directory exists
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return fmt.Errorf("failed to create parent directory: %w", err)
			}

			file, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("failed to create file: %w", err)
			}

			// Limit read size to prevent decompression bombs
			n, err := io.Copy(file, io.LimitReader(tarReader, MaxExtractedFileSize+1))
			if err != nil {
				file.Close()
				return fmt.Errorf("failed to write file: %w", err)
			}
			if err := file.Close(); err != nil {
				return fmt.Errorf("failed to close file %q: %w", header.Name, err)
			}
			if n > MaxExtractedFileSize {
				return fmt.Errorf("file %q exceeds maximum size limit (%d bytes)", header.Name, MaxExtractedFileSize)
			}

		case tar.TypeSymlink:
			// Validate symlink target to prevent traversal attacks
			if err := validateSymlinkTarget(dest, target, header.Linkname); err != nil {
				return err
			}
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return fmt.Errorf("failed to create parent directory: %w", err)
			}
			if err := os.Symlink(header.Linkname, target); err != nil {
				// Skip symlink errors on Windows
				continue
			}
		}
	}

	return nil
}

// validateSymlinkTarget checks that a symlink target resolves within the
// destination directory. This prevents symlink traversal attacks where a
// malicious archive creates symlinks pointing outside the extraction directory.
func validateSymlinkTarget(dest, symlinkPath, linkTarget string) error {
	// Reject absolute symlink targets — they always point outside dest
	if filepath.IsAbs(linkTarget) {
		return fmt.Errorf("invalid symlink target: absolute path %q", linkTarget)
	}

	// Resolve the symlink target relative to the symlink's parent directory
	resolvedTarget := filepath.Join(filepath.Dir(symlinkPath), linkTarget)
	resolvedTarget = filepath.Clean(resolvedTarget)

	// Ensure the resolved path stays within the destination directory
	destPrefix := filepath.Clean(dest) + string(os.PathSeparator)
	if !strings.HasPrefix(resolvedTarget, destPrefix) && resolvedTarget != filepath.Clean(dest) {
		return fmt.Errorf("invalid symlink target: %q escapes destination directory", linkTarget)
	}

	return nil
}

// copyDir recursively copies a directory
func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate destination path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		return copyFile(path, destPath)
	})
}

// copyFile copies a single file
func copyFile(src, dst string) (err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return err
	}
	defer func() {
		if cerr := dstFile.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// ClearCache removes all cached downloads
func (d *Downloader) ClearCache() error {
	entries, err := os.ReadDir(d.cachePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	for _, entry := range entries {
		path := filepath.Join(d.cachePath, entry.Name())
		if err := os.RemoveAll(path); err != nil {
			return err
		}
	}

	return nil
}

// GetCacheSize returns the total size of the cache in bytes
func (d *Downloader) GetCacheSize() (int64, error) {
	var size int64

	err := filepath.Walk(d.cachePath, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})

	return size, err
}
