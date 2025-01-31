package main

import (
	"fmt"
	"os"

	"example.com/cjon/tokenizer"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if command != "tokenize" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	filename := os.Args[2]
	lines, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	errCh := make(chan error)
	tokenCh := make(chan tokenizer.TokenStruct)

	go tokenizer.Tokenize(tokenCh, errCh, lines)
	Parse(tokenCh)

	err = <-errCh
	if err != nil {
		os.Exit(65)
	}
	os.Exit(0)
}

func Parse(tokens chan tokenizer.TokenStruct) {
	for t := range tokens {
		fmt.Println(t)
	}
}
