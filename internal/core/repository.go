package core

import (
	"bytes"
	"compress/zlib"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Repository struct {
	// Path is the absolute path to the root of the repository
	// Example:
	// "/home/user/myproject"
	Path string
}

const (
	objectFileMode    = 0444 // Read-only for all
	objectDirMode     = 0755 // Readable and executable by all, writable by owner
	defaultGitDirName = ".git"
	objectsDirName    = "objects"
)

// NewRepository creates a new Repository instance, ensuring the `.git` directory exists
func NewRepository(repoPath string) (*Repository, error) {
	gitPath := filepath.Join(repoPath, defaultGitDirName)
	if _, err := os.Stat(gitPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("not a Git repository: %s", repoPath)
	}

	return &Repository{Path: repoPath}, nil
}

// GitDir returns the path to the `.git` directory
func (r *Repository) GitDir() string {
	return filepath.Join(r.Path, defaultGitDirName)
}

// ObjectsDir returns the path to the `objects` directory inside `.git`
func (r *Repository) ObjectsDir() string {
	return filepath.Join(r.GitDir(), objectsDirName)
}

// WriteObject writes a serialized Git object to the `objects` directory
// The hash parameter should be a valid SHA-1 hash string.
// The data parameter contains the uncompressed object data to be stored.
// Returns an error if the write operation fails.
//
// Example:
//
//	hash := "1234567890abcdef..."
//	data := []byte("blob 16\x00Hello, World!")
//	err := repo.WriteObject(hash, data)
func (r *Repository) WriteObject(hash string, data []byte) error {

	// Add hash validation
	if len(hash) != 40 {
		return fmt.Errorf("invalid hash length: expected 40 characters, got %d", len(hash))
	}

	// Compress the data using zlib
	var compressed bytes.Buffer
	writer := zlib.NewWriter(&compressed)
	if _, err := writer.Write(data); err != nil {
		return fmt.Errorf("could not compress data: %w", err)
	}
	writer.Close()

	// Determine the directory and file paths
	dir := filepath.Join(r.ObjectsDir(), hash[:2])
	file := filepath.Join(dir, hash[2:])

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("could not create object directory: %w", err)
	}

	// Write the compressed data to the file
	if err := os.WriteFile(file, data, 0444); err != nil {
		return fmt.Errorf("could not write object file: %w", err)
	}

	return nil
}

func (r *Repository) ReadObject(hash string) ([]byte, error) {
	dir := filepath.Join(r.ObjectsDir(), hash[:2])
	file := filepath.Join(dir, hash[2:])

	compressed, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("could not read object file: %w", err)
	}

	reader, err := zlib.NewReader(bytes.NewReader(compressed))
	if err != nil {
		return nil, fmt.Errorf("could not create zlib reader: %w", err)
	}
	defer reader.Close()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(reader); err != nil {
		return nil, fmt.Errorf("could not decompress data: %w", err)
	}

	return buf.Bytes(), nil
}

// ValidateObject checks for valid Git objects or verify object integrity
func (r *Repository) ValidateObject(hash string) error {
	if len(hash) != 40 {
		return fmt.Errorf("invalid object hash: %s", hash)
	}

	objectPath := filepath.Join(r.ObjectsDir(), hash[:2], hash[2:])
	if _, err := os.Stat(objectPath); os.IsNotExist(err) {
		return fmt.Errorf("object not found: %s", hash)
	}
	return nil
}


// FindGitDir locates the .git directory starting from the current directory.
func (r *Repository) FindGitDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("could not get current working directory: %w", err)
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, defaultGitDirName)); err == nil {
			return filepath.Join(dir, defaultGitDirName), nil
		}

		parent := filepath.Dir(dir)
		if parent == dir { // Reached the root directory
			break
		}
		dir = parent
	}

	return "", errors.New(".git directory not found")
}
