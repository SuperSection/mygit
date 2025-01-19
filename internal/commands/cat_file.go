package commands

import (
	"bytes"
	"compress/zlib"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// CatFile handles the "cat-file" command, supporting flags like -p, -t, and -s.
func CatFile(flag, hash string) error {
	repoPath, err := os.Getwd() // Get the current working directory
	if err != nil {
		return fmt.Errorf("could not get current directory: %w", err)
	}

	// Ensure the .git directory exists
	gitPath := filepath.Join(repoPath, ".git")
	if _, err := os.Stat(gitPath); os.IsNotExist(err) {
		return errors.New("no .git repository found in the current directory")
	}

	if len(hash) < 2 {
		return fmt.Errorf("invalid hash")
	}

	// Compute the object path
	objectPath := filepath.Join(gitPath, "objects", hash[:2], hash[2:])

	// Reads and decompress the object
	file, err := os.Open(objectPath)
	if err != nil {
		return fmt.Errorf("could not open object file: %w", err)
	}
	defer file.Close()

	// Read and decompress the object
	decompressed, err := decompressZlib(file)
	if err != nil {
		return fmt.Errorf("error decompressing object file: %w", err)
	}

	// Parse the blob
	nullIndex := bytes.IndexByte(decompressed, 0)
	if nullIndex == -1 {
		return fmt.Errorf("invalid blob format")
	}

		// Extract type, size, and content
	header := string(decompressed[:nullIndex])
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid object header")
	}
	objectType := parts[0]
	objectSize := parts[1]
	content := decompressed[nullIndex+1:]

	// Handle flags
	switch flag {
	case "-p":
		if objectType != "blob" {
			return fmt.Errorf("object is not a blob: %s", objectType)
		}
		fmt.Print(string(content)) // Print without newline
	case "-t":
		fmt.Println(objectType)
	case "-s":
		fmt.Println((objectSize))
	default:
		return fmt.Errorf("unknown flag: %s", flag)
	}

	return nil
}


// decompressZlib decompresses a zlib-compressed reader.
func decompressZlib(r io.Reader) ([]byte, error) {
	var out bytes.Buffer
	zr, err := zlib.NewReader(r)
	if err != nil {
		return nil, err
	}
	defer zr.Close()

	_, err = io.Copy(&out, zr)
	return out.Bytes(), err
}
