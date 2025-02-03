package main

import (
	"errors"
	"fmt"
	"os"

	"example.com/cjon/parser"
	"example.com/cjon/tokenizer"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]
	var err error

	switch command {
	case "tokenize":
		errCh := make(chan error)
		tokenCh := make(chan tokenizer.TokenStruct)
		lines := getLines(os.Args[2])
		go tokenizer.Tokenize(tokenCh, errCh, lines)
		printTokens(tokenCh)
		err = <-errCh
		close(errCh)
	case "parse":
		errCh := make(chan error)
		tokenCh := make(chan tokenizer.TokenStruct)
		lines := getLines(os.Args[2])
		go tokenizer.Tokenize(tokenCh, errCh, lines)
		err = parser.Parse(tokenCh)
		if err == nil {
			err = <-errCh
		}
		close(errCh)
	default:
		err = errors.New("argument_error")
	}

	if err == nil {
		os.Exit(0)
	}

	switch err.Error() {
	case "argument_error":
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	case "syntax_error":
		os.Exit(65)
	case "parse_error":
		os.Exit(56)
	default:
		fmt.Fprintf(os.Stderr, "Unknown error: %s\n", command)
		os.Exit(-1)
	}
}

func getLines(filename string) []byte {
	lines, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}
	return lines
}

func printTokens(tokens chan tokenizer.TokenStruct) {
	for t := range tokens {
		fmt.Println(t)
	}
}
