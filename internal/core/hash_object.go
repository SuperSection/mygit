package core

import (
	"crypto/sha1"
	"fmt"
	"os"
)


// HashObject computes the SHA hash of a file. If `write` is true, it also writes the blob to .git/objects.
func HashObject(repo *Repository, filePath string, write bool) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("could not read file: %w", err)
	}

	// Use the Blob implementation
	blob := NewBlob(content)
	data, err := blob.Serialize(repo)
	if err != nil {
		return "", fmt.Errorf("could not serialize blob: %w", err)
	}

	// Compute the SHA-1 hash
	hash := sha1.Sum(data)
	hashHex := fmt.Sprintf("%x", hash)

	if write {
		if err := repo.WriteObject(hashHex, data); err != nil {
			return "", fmt.Errorf("could not write blob to .git/objects: %w", err)
		}
	}

	return hashHex, nil
}
