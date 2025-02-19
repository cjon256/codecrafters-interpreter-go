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
		tokenCh := make(chan token.Struct)
		lines := getLines(os.Args[2])
		go tokenizer.Tokenize(tokenCh, lines)
		err = printTokens(tokenCh)
	case "parse":
		perrCh := make(chan error)
		tokenCh := make(chan token.Struct)
		parserCh := make(chan parser.ASTnode)
		lines := getLines(os.Args[2])
		go tokenizer.Tokenize(tokenCh, lines)
		go parser.Parse(tokenCh, parserCh, perrCh)
		err = printAST(parserCh, perrCh)
	case "evaluate":
		errCh := make(chan error)
		tokenCh := make(chan token.Struct)
		parserCh := make(chan parser.ASTnode)
		lines := getLines(os.Args[2])
		go tokenizer.Tokenize(tokenCh, lines)
		go parser.Parse(tokenCh, parserCh, errCh)
		err = evaluateAST(parserCh, errCh)
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
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
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

func printTokens(tokens <-chan token.Struct) error {
	var err error = nil
	for t := range tokens {
		if t.Type == token.ERROR {
			err = errors.New(t.Literal)
			continue
		}
		fmt.Println(t)
	}
	return err
}

func printAST(astNodes <-chan parser.ASTnode, errCh <-chan error) error {
	initial := true
	for node := range astNodes {
		if initial {
			initial = false
		} else {
			fmt.Print(" ")
		}
		fmt.Print(node.String())
	}
	err := <-errCh
	if err != nil {
		return err
	}
	return nil
}

func evaluateAST(astNodes <-chan parser.ASTnode, errCh <-chan error) error {
	initial := true
	for node := range astNodes {
		if initial {
			initial = false
		} else {
			fmt.Print(" ")
		}
		fmt.Println(node.Execute())
	}
	err := <-errCh
	if err != nil {
		return err
	}
	return nil
}
