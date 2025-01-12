package commands

import (
	"fmt"
	"os"
	"path/filepath"
)

// Initializes a new Git repository in the current directory.
func Init() error {
	repoPath, err := os.Getwd() // Get the current working directory
	if err != nil {
		return fmt.Errorf("could not get current directory: %w", err)
	}

	// Check if the .git directory exists
	gitPath := filepath.Join(repoPath, ".git")
	if _, err := os.Stat(gitPath); !os.IsNotExist(err) {
		return fmt.Errorf("repository already exists at %s", gitPath)
	}

	// Create .git directory structure
	dirs := []string{
		"objects",
		"refs/heads",
		"refs/tags",
	}
	for _, dir := range dirs {
		path := filepath.Join(gitPath, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("could not create directory %s: %w", path, err)
		}
	}

	// Create HEAD file
	headPath := filepath.Join(gitPath, "HEAD")
	headFileContents := []byte("ref: refs/heads/main\n")
	if err := os.WriteFile(headPath, headFileContents, 0644); err != nil {
		return fmt.Errorf("could not write HEAD file: %w", err)
	}

	fmt.Println("Initialized empty Git repository in", gitPath)
	return nil
}
