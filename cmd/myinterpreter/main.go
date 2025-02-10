package main

import (
	"errors"
	"fmt"
	"os"

	"example.com/cjon/interpreter-starter-go/pkg/parser"
	"example.com/cjon/interpreter-starter-go/pkg/token"
	"example.com/cjon/interpreter-starter-go/pkg/tokenizer"
)

type ArgumentError struct{}

func (a ArgumentError) Error() string {
	// not used
	return ""
}

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
		serrCh := make(chan error)
		perrCh := make(chan error)
		tokenCh := make(chan token.Struct)
		parserCh := make(chan parser.ASTnode)
		lines := getLines(os.Args[2])
		go tokenizer.Tokenize(tokenCh, serrCh, lines)
		go parser.Parse(tokenCh, parserCh, perrCh)
		err = printAST(parserCh, serrCh, perrCh)
	case "execute":
		fmt.Println("false")
	default:
		err = errors.New("argument_error")
	}

	if err == nil {
		os.Exit(0)
	}

	switch err.(type) {
	case ArgumentError:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	default:
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(65)
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

func printAST(astNodes chan parser.ASTnode, serrCh chan error, perrCh chan error) error {
	initial := true
	select {
	case node := <-astNodes:
		if initial {
			initial = false
		} else {
			fmt.Print(" ")
		}
		fmt.Println(node)
	case err := <-serrCh:
		if err != nil {
			return err
		}
	case err := <-perrCh:
		if err != nil {
			return err
		}
	}
	err := <-perrCh
	if err != nil {
		return err
	}
	err = <-serrCh
	if err != nil {
		return err
	}
	return nil
}
