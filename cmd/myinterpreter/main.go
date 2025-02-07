package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"example.com/cjon/parser"
	"example.com/cjon/token"
	"example.com/cjon/tokenizer"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	var err error
	command := os.Args[1]

	switch command {
	case "tokenize":
		errCh := make(chan error)
		tokenCh := make(chan token.Struct)
		lines := getLines(os.Args[2])
		go tokenizer.Tokenize(tokenCh, errCh, lines)
		printTokens(tokenCh)
		err = <-errCh
	case "parse":
		errCh := make(chan error)
		tokenCh := make(chan token.Struct)
		lines := getLines(os.Args[2])
		go tokenizer.Tokenize(tokenCh, errCh, lines)
		err = parser.Parse(tokenCh)
		if err == nil {
			err = <-errCh
		}
	default:
		err = errors.New("argument_error")
	}

	if err == nil {
		os.Exit(0)
	}

	errStr := err.Error()
	switch {
	case errStr == "argument_error":
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	case errStr == "syntax_error":
		os.Exit(65)
	case strings.HasPrefix("parse_error", errStr):
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

func printTokens(tokens chan token.Struct) {
	for t := range tokens {
		fmt.Println(t)
	}
}
