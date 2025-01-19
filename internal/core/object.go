package core

import "fmt"

// GitObject represents the common functionality of all Git objects.
type GitObject interface {
	Serialize(repo *Repository) ([]byte, error) // Convert the object into bytes to be stored
	Deserialize(data []byte) error              // Populate the object from raw bytes
	Init()                                      // Initialize a new, empty object
	Type() string                               // Returns the type of the Git object (e.g., "blob", "tree")
}


type BaseGitObject struct {
	// Add common fields here if needed
	// e.g., Metadata or type fields
}

// Default implementations for GitObject interface
func (o *BaseGitObject) Init() {
	// Default init: No-op
}

func (o *BaseGitObject) Serialize(repo *Repository) ([]byte, error) {
	return nil, fmt.Errorf("Serialize method not implemented")
}

func (o *BaseGitObject) Deserialize(data []byte) error {
	return fmt.Errorf("Deserialize method not implemented")
}

func (o *BaseGitObject) Type() string {
	return "base"
}
