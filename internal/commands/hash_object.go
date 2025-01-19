package commands

import (
	"flag"
	"fmt"
	"log"

	"github.com/supersection/mygit/internal/core"
)

// git hash-object is used to compute the SHA hash of a Git object.
// When used with the -w flag, it also writes the object to the .git/objects directory.
func HandleHashObject(args []string) {

	// Parse flags
	flags := flag.NewFlagSet("hash-object", flag.ExitOnError)
	write := flags.Bool("w", false, "Write the object to .git/objects")
	if err := flags.Parse(args); err != nil {
		log.Fatal(err)
	}

	// Validate arguments
	files := flags.Args()
	if len(files) != 1 {
		log.Fatal("Usage: mygit hash-object [-w] <file>")
	}
	filePath := files[0]


	// Initialize repository
	repo, err := core.NewRepository(".")
	if err != nil {
		log.Fatalf("Could not find Git repository: %v", err)
	}

	shaHash, err := core.HashObject(repo, filePath, *write)
	if err != nil {
		log.Fatalf("Error in hash-object: %v", err)
	}

	fmt.Println(shaHash)
}
