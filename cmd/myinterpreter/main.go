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
		tokenCh := make(chan token.Struct)
		parserCh := make(chan parser.ASTnode)
		lines := getLines(os.Args[2])
		go tokenizer.Tokenize(tokenCh, lines)
		go parser.Parse(tokenCh, parserCh)
		err = printAST(parserCh)
	case "evaluate":
		tokenCh := make(chan token.Struct)
		parserCh := make(chan parser.ASTnode)
		lines := getLines(os.Args[2])
		go tokenizer.Tokenize(tokenCh, lines)
		go parser.Parse(tokenCh, parserCh)
		err = evaluateAST(parserCh)
	default:
		err = errors.New("argument_error")
	}

	if err == nil {
		os.Exit(0)
	}

	switch err.(type) {
	case ArgumentError:
		os.Exit(1)
	default:
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
	var err error
	for t := range tokens {
		if t.Type == token.ERROR {
			fmt.Fprint(os.Stderr, t.Literal)
			err = errors.New(t.Literal)
		} else {
			fmt.Println(t)
		}
	}
	return err
}

func printAST(astNodes <-chan parser.ASTnode) error {
	var err error
	initial := true
	for node := range astNodes {
		_, ok := node.(parser.ASTerror) // so ugly
		if ok {
			fmt.Fprintln(os.Stderr, node.String())
			err = errors.New(node.String())
			continue
		}
		if initial {
			initial = false
		} else {
			fmt.Print(" ")
		}
		fmt.Print(node.String())
	}
	return err
}

func evaluateAST(astNodes <-chan parser.ASTnode) error {
	var err error
	initial := true
	for node := range astNodes {
		_, ok := node.(parser.ASTerror)
		if ok {
			fmt.Fprintln(os.Stderr, node.Evaluate())
			err = errors.New(node.String())
			continue
		}
		if initial {
			initial = false
		} else {
			fmt.Print(" ")
		}
		fmt.Println(node.Evaluate())
	}
	return err
}
