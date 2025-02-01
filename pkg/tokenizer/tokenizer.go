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
	IDENTIFIER
)

func Tokenize(tokens chan TokenStruct, errCh chan error, line []byte) {
	var err error = nil
	i := 0
	lineNumber := 1
	peek := func() (bool, byte) {
		if i+1 == len(line) {
			return false, ' '
		}
		return true, line[i+1]
	}

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
			ok, nextByte := peek()
			if !ok {
				// at end of line
				break
			}
			if nextByte == '.' {
				if dotSeen {
					break
				}
				dotSeen = true
			} else if !unicode.IsDigit(rune(nextByte)) {
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

	isIdentifierByte := func(c byte) bool {
		// Check if the byte value falls within the range of alphanumeric characters
		return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c == '_')
	}

	tokenizeIdentifier := func() {
		ts := []byte{}
		for ; i < len(line); i++ {
			ts = append(ts, line[i])
			ok, nextByte := peek()
			if !ok {
				// at end of line
				break
			}
			if !(unicode.IsDigit(rune(nextByte)) || isIdentifierByte(nextByte)) {
				break
			}
		}
		str := string(ts)
		tokens <- TokenStruct{IDENTIFIER, str, "null"}
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
			if isIdentifierByte(line[i]) {
				tokenizeIdentifier()
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
