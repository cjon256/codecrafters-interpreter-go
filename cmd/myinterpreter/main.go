package main

import (
	"errors"
	"fmt"
	"os"
)

type tokenStruct struct {
	Type      TokenType
	Str       string
	NullThing interface{}
}

func (t tokenStruct) String() string {
	nullStr := ""
	if t.NullThing == nil {
		nullStr = "null"
	}
	return fmt.Sprintf("%s %s %s", t.Type, t.Str, nullStr)
}

//go:generate stringer -type=TokenType
type TokenType int

const (
	EOF TokenType = iota
	LEFT_PAREN
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	SEMICOLON
	COMMA
	PLUS
	MINUS
	STAR
	BANG_EQUAL
	EQUAL_EQUAL
	LESS_EQUAL
	GREATER_EQUAL
	LESS
	GREATER
	SLASH
	DOT
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
	ch := make(chan tokenStruct)

	go Tokenize(ch, errCh, lines, os.Stderr)
	Parse(ch)

	err = <-errCh
	if err != nil {
		os.Exit(65)
	}
	os.Exit(0)
}

func Parse(tokens chan tokenStruct) {
	for t := range tokens {
		fmt.Println(t)
	}
}

func Tokenize(tokens chan tokenStruct, errCh chan error, line []byte, errout *os.File) {
	var err error = nil
	lineNumber := 1

loop:
	for i := 0; i < len(line); i++ {
		switch line[i] {
		case '(':
			tokens <- tokenStruct{LEFT_PAREN, "(", nil}
		case ')':
			tokens <- tokenStruct{RIGHT_PAREN, ")", nil}
		case '{':
			tokens <- tokenStruct{LEFT_BRACE, "{", nil}
		case '}':
			tokens <- tokenStruct{RIGHT_BRACE, "}", nil}
		case ';':
			tokens <- tokenStruct{SEMICOLON, ";", nil}
		case ',':
			tokens <- tokenStruct{COMMA, ",", nil}
		case '+':
			tokens <- tokenStruct{PLUS, "+", nil}
		case '-':
			tokens <- tokenStruct{MINUS, "-", nil}
		case '*':
			tokens <- tokenStruct{STAR, "*", nil}
		case '!':
			if i+1 < len(line) && line[i+1] == '=' {
				tokens <- tokenStruct{BANG_EQUAL, "!=", nil}
				i++
			} else {
				err = errors.New("oops")
				fmt.Fprintf(errout, "[line %d] Error: Unexpected character: %s\n", lineNumber, string(line[i]))
			}
		case '=':
			if i+1 < len(line) && line[i+1] == '=' {
				tokens <- tokenStruct{EQUAL_EQUAL, "==", nil}
				i++
			} else {
				err = errors.New("oops")
				fmt.Fprintf(errout, "[line %d] Error: Unexpected character: %s\n", lineNumber, string(line[i]))
			}
		case '<':
			if i+1 < len(line) && line[i+1] == '=' {
				tokens <- tokenStruct{LESS_EQUAL, "<=", nil}
				i++
			} else {
				tokens <- tokenStruct{LESS, "<", nil}
			}
		case '>':
			if i+1 < len(line) && line[i+1] == '=' {
				tokens <- tokenStruct{GREATER_EQUAL, ">=", nil}
				i++
			} else {
				tokens <- tokenStruct{GREATER, ">", nil}
			}
		case '/':
			if i+1 < len(line) && line[i+1] == '/' {
				// handle comments
				break loop
			} else {
				tokens <- tokenStruct{SLASH, "/", nil}
			}
		case '.':
			tokens <- tokenStruct{DOT, ".", nil}
		case ' ':
			// ignore
		default:
			err = errors.New("syntax_error")
			fmt.Fprintf(errout, "[line %d] Error: Unexpected character: %s\n", lineNumber, string(line[i]))
			continue
		}
	}
	tokens <- tokenStruct{EOF, "", nil}
	close(tokens)
	errCh <- err
	close(errCh)
}
