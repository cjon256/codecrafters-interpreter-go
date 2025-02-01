package tokenizer

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"unicode"
)

type TokenStruct struct {
	Type    TokenType
	Lexeme  string
	Literal string
}

func (t TokenStruct) String() string {
	return fmt.Sprintf("%s %s %s", t.Type, t.Lexeme, t.Literal)
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
	STRING
	NUMBER
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
	i := 0
	lineNumber := 1
	tokenizeString := func() error {
		terminated := false
		ln := lineNumber
		ts := []byte{}
		for ; i < len(line); i++ {
			if line[i] == '"' {
				terminated = true
				break
			}
			if line[i] == '\n' {
				lineNumber++
				break
			}
			ts = append(ts, line[i])
		}
		if !terminated {
			fmt.Fprintf(os.Stderr, "[line %d] Error: Unterminated string.\n", ln)
			return errors.New("syntax_error")
		}
		str := string(ts)
		qstr := "\"" + str + "\""
		tokens <- TokenStruct{STRING, qstr, str}
		return nil
	}

	tokenizeNumber := func() {
		ts := []byte{}
		dotSeen := false
		for ; i < len(line); i++ {
			// fmt.Fprintf(os.Stderr, "num (%d): %v\n", i, line[i])
			ts = append(ts, line[i])
			if i == len(line)-1 {
				// at end of line
				break
			}
			if line[i+1] == '.' {
				if dotSeen {
					break
				}
				dotSeen = true
			} else if !unicode.IsDigit(rune(line[i+1])) {
				break
			}
		}
		str := string(ts)
		nstr := str
		if !dotSeen {
			nstr = nstr + ".0"
		} else {
			// this bit just removes trailing zeros
			nstr = strings.TrimRight(nstr, "0")
			// and adds one back in if there were only zeros
			if strings.HasSuffix(nstr, ".") {
				nstr = nstr + "0"
			}
		}
		// fmt.Fprintf(os.Stderr, "num (%v): %s,%s ... %s\n", dotSeen, str, nstr, line[i:])
		tokens <- TokenStruct{NUMBER, str, nstr}
	}

	for ; i < len(line); i++ {
		switch line[i] {
		case '(':
			tokens <- TokenStruct{LEFT_PAREN, "(", "null"}
		case ')':
			tokens <- TokenStruct{RIGHT_PAREN, ")", "null"}
		case '{':
			tokens <- TokenStruct{LEFT_BRACE, "{", "null"}
		case '}':
			tokens <- TokenStruct{RIGHT_BRACE, "}", "null"}
		case ';':
			tokens <- TokenStruct{SEMICOLON, ";", "null"}
		case ',':
			tokens <- TokenStruct{COMMA, ",", "null"}
		case '+':
			tokens <- TokenStruct{PLUS, "+", "null"}
		case '-':
			tokens <- TokenStruct{MINUS, "-", "null"}
		case '*':
			tokens <- TokenStruct{STAR, "*", "null"}
		case '!':
			if i+1 < len(line) && line[i+1] == '=' {
				i++
				tokens <- TokenStruct{BANG_EQUAL, "!=", "null"}
			} else {
				tokens <- TokenStruct{BANG, "!", "null"}
			}
		case '=':
			if i+1 < len(line) && line[i+1] == '=' {
				tokens <- TokenStruct{EQUAL_EQUAL, "==", "null"}
				i++
			} else {
				tokens <- TokenStruct{EQUAL, "=", "null"}
			}
		case '<':
			if i+1 < len(line) && line[i+1] == '=' {
				tokens <- TokenStruct{LESS_EQUAL, "<=", "null"}
				i++
			} else {
				tokens <- TokenStruct{LESS, "<", "null"}
			}
		case '>':
			if i+1 < len(line) && line[i+1] == '=' {
				tokens <- TokenStruct{GREATER_EQUAL, ">=", "null"}
				i++
			} else {
				tokens <- TokenStruct{GREATER, ">", "null"}
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
				tokens <- TokenStruct{SLASH, "/", "null"}
			}
		case '.':
			tokens <- TokenStruct{DOT, ".", "null"}
		case ' ':
			// ignore
		case '\t':
			// ignore
		case '\n':
			lineNumber++
		case '"':
			i++
			err = tokenizeString()
		default:
			if unicode.IsDigit(rune(line[i])) {
				tokenizeNumber()
				continue
			}
			err = errors.New("syntax_error")
			fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: %s\n", lineNumber, string(line[i]))
		}
	}
	tokens <- TokenStruct{EOF, "", "null"}
	close(tokens)
	errCh <- err
	close(errCh)
}
