package tokenizer

import (
	"errors"
	"fmt"
	"os"
)

type TokenStruct struct {
	Type      TokenType
	Str       string
	NullThing interface{}
}

func (t TokenStruct) String() string {
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
	EQUAL
	STAR
	BANG_EQUAL
	EQUAL_EQUAL
	LESS_EQUAL
	GREATER_EQUAL
	LESS
	GREATER
	SLASH
	DOT
	BANG
)

// func main() {
// 	if len(os.Args) < 3 {
// 		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
// 		os.Exit(1)
// 	}
//
// 	command := os.Args[1]
//
// 	if command != "tokenize" {
// 		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
// 		os.Exit(1)
// 	}
//
// 	filename := os.Args[2]
// 	lines, err := os.ReadFile(filename)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
// 		os.Exit(1)
// 	}
//
// 	errCh := make(chan error)
// 	tokenCh := make(chan tokenStruct)
//
// 	go Tokenize(tokenCh, errCh, lines)
// 	Parse(tokenCh)
//
// 	err = <-errCh
// 	if err != nil {
// 		os.Exit(65)
// 	}
// 	os.Exit(0)
// }

func Parse(tokens chan TokenStruct) {
	for t := range tokens {
		fmt.Println(t)
	}
}

func Tokenize(tokens chan TokenStruct, errCh chan error, line []byte) {
	var err error = nil
	lineNumber := 1

	for i := 0; i < len(line); i++ {
		switch line[i] {
		case '(':
			tokens <- TokenStruct{LEFT_PAREN, "(", nil}
		case ')':
			tokens <- TokenStruct{RIGHT_PAREN, ")", nil}
		case '{':
			tokens <- TokenStruct{LEFT_BRACE, "{", nil}
		case '}':
			tokens <- TokenStruct{RIGHT_BRACE, "}", nil}
		case ';':
			tokens <- TokenStruct{SEMICOLON, ";", nil}
		case ',':
			tokens <- TokenStruct{COMMA, ",", nil}
		case '+':
			tokens <- TokenStruct{PLUS, "+", nil}
		case '-':
			tokens <- TokenStruct{MINUS, "-", nil}
		case '*':
			tokens <- TokenStruct{STAR, "*", nil}
		case '!':
			if i+1 < len(line) && line[i+1] == '=' {
				i++
				tokens <- TokenStruct{BANG_EQUAL, "!=", nil}
			} else {
				tokens <- TokenStruct{BANG, "!", nil}
			}
		case '=':
			if i+1 < len(line) && line[i+1] == '=' {
				tokens <- TokenStruct{EQUAL_EQUAL, "==", nil}
				i++
			} else {
				tokens <- TokenStruct{EQUAL, "=", nil}
			}
		case '<':
			if i+1 < len(line) && line[i+1] == '=' {
				tokens <- TokenStruct{LESS_EQUAL, "<=", nil}
				i++
			} else {
				tokens <- TokenStruct{LESS, "<", nil}
			}
		case '>':
			if i+1 < len(line) && line[i+1] == '=' {
				tokens <- TokenStruct{GREATER_EQUAL, ">=", nil}
				i++
			} else {
				tokens <- TokenStruct{GREATER, ">", nil}
			}
		case '/':
			if i+1 < len(line) && line[i+1] == '/' {
				// handle comments
				for i < len(line) {
					if line[i] == '\n' {
						lineNumber++
						break
					}
					i++
				}
			} else {
				tokens <- TokenStruct{SLASH, "/", nil}
			}
		case '.':
			tokens <- TokenStruct{DOT, ".", nil}
		case ' ':
			// ignore
		case '\t':
			// ignore
		case '\n':
			lineNumber++
		default:
			err = errors.New("syntax_error")
			fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: %s\n", lineNumber, string(line[i]))
			continue
		}
	}
	tokens <- TokenStruct{EOF, "", nil}
	close(tokens)
	errCh <- err
	close(errCh)
}
