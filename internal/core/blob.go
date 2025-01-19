package core

import (
	"fmt"
)

type Blob struct {
	BaseGitObject
	Content []byte
}


// NewBlob creates a new Blob object
func NewBlob(data []byte) *Blob {
	blob := &Blob{}
	if data != nil {
		_ = blob.Deserialize(data)
	} else {
		blob.Init()
	}
	return blob
}

// Type returns the type of the object
func (b *Blob) Type() string {
	return "blob"
}

// Init initializes a new empty Blob object
func (b *Blob) Init() {
	b.Content = []byte{}
}

// Serialize converts the Blob into its Git object format
func (b *Blob) Serialize(repo *Repository) ([]byte, error) {
	header := fmt.Sprintf("blob %d\000", len(b.Content))
	return append([]byte(header), b.Content...), nil
}

// Deserialize populates the Blob object from raw data
func (b *Blob) Deserialize(data []byte) error {
	b.Content = data // For blobs, raw data is just the content
	return nil
}
