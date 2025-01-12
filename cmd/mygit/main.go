package main

import (
	"fmt"
	"os"

	"github.com/supersection/mygit/internal/commands"
)

// Usage: your_program.sh <command> <arg1> <arg2> ...
func main() {

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
		os.Exit(1)
	}

	switch command := os.Args[1]; command {
	case "init":
		if err := commands.Init(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "cat-file":
		if len(os.Args) != 4 {
			fmt.Println("Usage: mygit cat-file < -p|-t|-s > <object-hash>")
			os.Exit(1)
		}
		flag, hash := os.Args[2], os.Args[3]
		if err := commands.CatFile(flag, hash); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}
