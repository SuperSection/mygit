package main

import (
	"fmt"
	"os"

	"github.com/supersection/mygit/internal/commands"
)

// Usage: your_program.sh <command> <arg1> <arg2> ...
func main() {

	if len(os.Args) < 2 {
		showUsageAndExit()
	}

	switch command := os.Args[1]; command {
	case "init":
		if err := commands.Init(); err != nil {
			exitWithError(err)
		}

	case "cat-file":
		if len(os.Args) != 4 {
			fmt.Println("Usage: mygit cat-file < -p|-t|-s > <object-hash>")
			os.Exit(1)
		}
		flag, hash := os.Args[2], os.Args[3]
		if err := commands.CatFile(flag, hash); err != nil {
			exitWithError(err)
		}

	case "hash-object":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: mygit hash-object [-w] <file>")
			os.Exit(1)
		}
		commands.HandleHashObject(os.Args[2:])

	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		showUsageAndExit()
	}
}



func showUsageAndExit() {
	fmt.Fprintln(os.Stderr, "Usage: mygit <command> [<args>...]")
	fmt.Fprintln(os.Stderr, "Commands:")
	fmt.Fprintln(os.Stderr, "  init         Initialize a new Git repository")
	fmt.Fprintln(os.Stderr, "  hash-object  Compute SHA hash of a file or write it to .git/objects")
	fmt.Fprintln(os.Stderr, "  cat-file     Display information about a Git object")
	os.Exit(1)
}

func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}
