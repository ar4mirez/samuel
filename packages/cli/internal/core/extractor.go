package core

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Extractor handles extracting framework files to a destination
type Extractor struct {
	sourcePath string
	destPath   string
}

// NewExtractor creates a new extractor
func NewExtractor(sourcePath, destPath string) *Extractor {
	return &Extractor{
		sourcePath: sourcePath,
		destPath:   destPath,
	}
}

// ExtractResult contains the result of an extraction
type ExtractResult struct {
	FilesCreated  []string
	DirsCreated   []string
	FilesSkipped  []string
	Errors        []error
}

// Extract copies specific files from source to destination
// The paths parameter contains destination paths (e.g., ".agent/language-guides/go.md")
// Source paths are calculated by prepending TemplatePrefix (e.g., "template/.agent/language-guides/go.md")
func (e *Extractor) Extract(paths []string, force bool) (*ExtractResult, error) {
	result := &ExtractResult{
		FilesCreated: make([]string, 0),
		DirsCreated:  make([]string, 0),
		FilesSkipped: make([]string, 0),
		Errors:       make([]error, 0),
	}

	// Create destination directory if it doesn't exist
	if err := os.MkdirAll(e.destPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create destination directory: %w", err)
	}

	for _, path := range paths {
		// Source path includes template/ prefix, destination path does not
		srcPath := filepath.Join(e.sourcePath, TemplatePrefix, path)
		dstPath := filepath.Join(e.destPath, path)

		// Check if source exists
		srcInfo, err := os.Stat(srcPath)
		if err != nil {
			if os.IsNotExist(err) {
				result.Errors = append(result.Errors, fmt.Errorf("source not found: %s", path))
				continue
			}
			result.Errors = append(result.Errors, err)
			continue
		}

		// Handle directories
		if srcInfo.IsDir() {
			if err := e.extractDir(srcPath, dstPath, result, force); err != nil {
				result.Errors = append(result.Errors, err)
			}
			continue
		}

		// Handle files
		if err := e.extractFile(srcPath, dstPath, result, force); err != nil {
			result.Errors = append(result.Errors, err)
		}
	}

	return result, nil
}

// extractFile copies a single file
func (e *Extractor) extractFile(srcPath, dstPath string, result *ExtractResult, force bool) error {
	// Check if destination exists
	if _, err := os.Stat(dstPath); err == nil {
		if !force {
			relPath, _ := filepath.Rel(e.destPath, dstPath)
			result.FilesSkipped = append(result.FilesSkipped, relPath)
			return nil
		}
	}

	// Ensure parent directory exists
	parentDir := filepath.Dir(dstPath)
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", parentDir, err)
	}

	// Copy file
	if err := copyFile(srcPath, dstPath); err != nil {
		return fmt.Errorf("failed to copy %s: %w", srcPath, err)
	}

	relPath, _ := filepath.Rel(e.destPath, dstPath)
	result.FilesCreated = append(result.FilesCreated, relPath)

	return nil
}

// extractDir recursively copies a directory
func (e *Extractor) extractDir(srcPath, dstPath string, result *ExtractResult, force bool) error {
	// Create destination directory
	if err := os.MkdirAll(dstPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dstPath, err)
	}

	relDir, _ := filepath.Rel(e.destPath, dstPath)
	result.DirsCreated = append(result.DirsCreated, relDir)

	return filepath.Walk(srcPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate relative path from source
		relPath, err := filepath.Rel(srcPath, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(dstPath, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		return e.extractFile(path, destPath, result, force)
	})
}

// ExtractAll extracts all framework files from the template/ directory
func (e *Extractor) ExtractAll(force bool) (*ExtractResult, error) {
	// Get all files in template/ subdirectory of source
	templateDir := filepath.Join(e.sourcePath, TemplatePrefix)
	var paths []string
	err := filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// If template directory doesn't exist, return empty result
			if os.IsNotExist(err) && path == templateDir {
				return filepath.SkipAll
			}
			return err
		}

		// Calculate path relative to template/ directory (this is the destination path)
		relPath, err := filepath.Rel(templateDir, path)
		if err != nil {
			return err
		}

		// Skip hidden files and certain directories
		if shouldSkip(relPath) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if !info.IsDir() {
			paths = append(paths, relPath)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return e.Extract(paths, force)
}

// shouldSkip returns true if the path should be skipped during extraction
// Since we now extract from template/ directory, we only need minimal filtering
func shouldSkip(path string) bool {
	// Skip the root path
	if path == "." {
		return false
	}

	// Skip git directory (shouldn't be in template, but just in case)
	if strings.HasPrefix(path, ".git") {
		return true
	}

	// Skip node_modules (shouldn't be in template, but just in case)
	if strings.Contains(path, "node_modules") {
		return true
	}

	return false
}

// ValidateExtraction checks if extracted files are valid
func (e *Extractor) ValidateExtraction(paths []string) []string {
	var missing []string

	for _, path := range paths {
		fullPath := filepath.Join(e.destPath, path)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			missing = append(missing, path)
		}
	}

	return missing
}

// FileExists checks if a file exists in the destination
func (e *Extractor) FileExists(path string) bool {
	fullPath := filepath.Join(e.destPath, path)
	_, err := os.Stat(fullPath)
	return err == nil
}

// ReadFile reads a file from the destination
func (e *Extractor) ReadFile(path string) ([]byte, error) {
	fullPath := filepath.Join(e.destPath, path)
	return os.ReadFile(fullPath)
}

// WriteFile writes content to a file in the destination
func (e *Extractor) WriteFile(path string, content []byte) error {
	fullPath := filepath.Join(e.destPath, path)

	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return err
	}

	return os.WriteFile(fullPath, content, 0644)
}

// RemoveFile removes a file from the destination
func (e *Extractor) RemoveFile(path string) error {
	fullPath := filepath.Join(e.destPath, path)
	return os.Remove(fullPath)
}

// BackupFile creates a backup of a file
func (e *Extractor) BackupFile(path, backupDir string) error {
	srcPath := filepath.Join(e.destPath, path)
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		return nil // Nothing to backup
	}

	dstPath := filepath.Join(backupDir, path)

	// Ensure backup directory exists
	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return err
	}

	// Read source
	content, err := os.ReadFile(srcPath)
	if err != nil {
		return err
	}

	// Write backup
	return os.WriteFile(dstPath, content, 0644)
}

// RestoreBackup restores files from a backup directory
func (e *Extractor) RestoreBackup(backupDir string) error {
	return filepath.Walk(backupDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(backupDir, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(e.destPath, relPath)

		// Read backup
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// Ensure destination directory exists
		if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
			return err
		}

		// Restore file
		return os.WriteFile(dstPath, content, 0644)
	})
}

// GetSourcePath returns the source path
func (e *Extractor) GetSourcePath() string {
	return e.sourcePath
}

// GetDestPath returns the destination path
func (e *Extractor) GetDestPath() string {
	return e.destPath
}

// CopyFromCache copies a single file from the cache source directly
// The filePath is the destination path; source is found in template/ subdirectory
func CopyFromCache(cachePath, destPath, filePath string) error {
	srcFile := filepath.Join(cachePath, TemplatePrefix, filePath)
	dstFile := filepath.Join(destPath, filePath)

	// Ensure destination directory exists
	if err := os.MkdirAll(filepath.Dir(dstFile), 0755); err != nil {
		return err
	}

	src, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(dstFile)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}
