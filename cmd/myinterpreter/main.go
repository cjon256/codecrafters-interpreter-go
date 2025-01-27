package main

import (
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
	tokens := Tokenize(lines)
	Parse(tokens)
}

func Parse(tokens []tokenStruct) {
	for _, t := range tokens {
		fmt.Println(t)
	}
}

func Tokenize(lines []byte) []tokenStruct {
	tokens := []tokenStruct{}

	for i := 0; i < len(lines); i++ {
		switch lines[i] {
		case '(':
			tokens = append(tokens, tokenStruct{LEFT_PAREN, "(", nil})
		case ')':
			tokens = append(tokens, tokenStruct{RIGHT_PAREN, ")", nil})
		case '{':
			tokens = append(tokens, tokenStruct{LEFT_BRACE, "{", nil})
		case '}':
			tokens = append(tokens, tokenStruct{RIGHT_BRACE, "}", nil})
		case ';':
			tokens = append(tokens, tokenStruct{SEMICOLON, ";", nil})
		case ',':
			tokens = append(tokens, tokenStruct{COMMA, ",", nil})
		case '+':
			tokens = append(tokens, tokenStruct{PLUS, "+", nil})
		case '-':
			tokens = append(tokens, tokenStruct{MINUS, "-", nil})
		case '*':
			tokens = append(tokens, tokenStruct{STAR, "*", nil})
		case '!':
			if i+1 < len(lines) && lines[i+1] == '=' {
				tokens = append(tokens, tokenStruct{BANG_EQUAL, "!=", nil})
				i++
			} else {
				// parse error?
			}
		case '=':
			if i+1 < len(lines) && lines[i+1] == '=' {
				tokens = append(tokens, tokenStruct{EQUAL_EQUAL, "==", nil})
				i++
			} else {
				// parse error?
			}
		case '<':
			if i+1 < len(lines) && lines[i+1] == '=' {
				tokens = append(tokens, tokenStruct{LESS_EQUAL, "<=", nil})
				i++
			} else {
				tokens = append(tokens, tokenStruct{LESS, "<", nil})
			}
		case '>':
			if i+1 < len(lines) && lines[i+1] == '=' {
				tokens = append(tokens, tokenStruct{GREATER_EQUAL, ">=", nil})
				i++
			} else {
				tokens = append(tokens, tokenStruct{GREATER, ">", nil})
			}
		case '/':
			tokens = append(tokens, tokenStruct{SLASH, "/", nil})
		case '.':
			tokens = append(tokens, tokenStruct{DOT, ".", nil})

		default:
			// ignore any whitespace or other unknown characters
		}
	}
	tokens = append(tokens, tokenStruct{EOF, "", nil})
	return tokens
}
