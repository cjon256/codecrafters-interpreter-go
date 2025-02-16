package tokenizer

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"unicode"

	"example.com/cjon/interpreter-starter-go/pkg/token"
)

func Tokenize(tokens chan token.Struct, line []byte) {
	var err error = nil
	i := 0
	lineNumber := 1
	peek := func() (byte, bool) {
		if i+1 == len(line) {
			return ' ', false
		}
		return line[i+1], true
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
		tokens <- token.Struct{token.STRING, qstr, str, lineNumber}
		return nil
	}

	tokenizeNumber := func() {
		ts := []byte{}
		dotSeen := false
		for ; i < len(line); i++ {
			// fmt.Fprintf(os.Stderr, "num (%d): %v\n", i, line[i])
			ts = append(ts, line[i])
			nextByte, ok := peek()
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
		// restLn := line[i:]
		// fmt.Fprintf(os.Stderr, "num (%v): %s,%s ... %s\n", dotSeen, str, nstr, restLn)
		tokens <- token.Struct{Type: token.NUMBER, Lexeme: str, Literal: nstr}
	}

	isIdentifierByte := func(c byte) bool {
		// Check if the byte value falls within the range of alphanumeric characters
		return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c == '_')
	}

	tokenizeIdentifier := func() {
		ts := []byte{}
		for ; i < len(line); i++ {
			ts = append(ts, line[i])
			nextByte, ok := peek()
			if !ok {
				// at end of line
				break
			}
			if !(unicode.IsDigit(rune(nextByte)) || isIdentifierByte(nextByte)) {
				break
			}
		}
		str := string(ts)
		lexeme, ok := token.KEYWORDS[str]
		if !ok {
			tokens <- token.Struct{token.IDENTIFIER, str, "null", lineNumber}
		} else {
			tokens <- token.Struct{lexeme, str, "null", lineNumber}
		}
	}

	for ; i < len(line); i++ {
		switch line[i] {
		case '(':
			tokens <- token.Struct{token.LEFT_PAREN, "(", "null", lineNumber}
		case ')':
			tokens <- token.Struct{token.RIGHT_PAREN, ")", "null", lineNumber}
		case '{':
			tokens <- token.Struct{token.LEFT_BRACE, "{", "null", lineNumber}
		case '}':
			tokens <- token.Struct{token.RIGHT_BRACE, "}", "null", lineNumber}
		case ';':
			tokens <- token.Struct{token.SEMICOLON, ";", "null", lineNumber}
		case ',':
			tokens <- token.Struct{token.COMMA, ",", "null", lineNumber}
		case '+':
			tokens <- token.Struct{token.PLUS, "+", "null", lineNumber}
		case '-':
			tokens <- token.Struct{token.MINUS, "-", "null", lineNumber}
		case '*':
			tokens <- token.Struct{token.STAR, "*", "null", lineNumber}
		case '!':
			if i+1 < len(line) && line[i+1] == '=' {
				i++
				tokens <- token.Struct{token.BANG_EQUAL, "!=", "null", lineNumber}
			} else {
				tokens <- token.Struct{token.BANG, "!", "null", lineNumber}
			}
		case '=':
			if i+1 < len(line) && line[i+1] == '=' {
				tokens <- token.Struct{token.EQUAL_EQUAL, "==", "null", lineNumber}
				i++
			} else {
				tokens <- token.Struct{token.EQUAL, "=", "null", lineNumber}
			}
		case '<':
			if i+1 < len(line) && line[i+1] == '=' {
				tokens <- token.Struct{token.LESS_EQUAL, "<=", "null", lineNumber}
				i++
			} else {
				tokens <- token.Struct{token.LESS, "<", "null", lineNumber}
			}
		case '>':
			if i+1 < len(line) && line[i+1] == '=' {
				tokens <- token.Struct{token.GREATER_EQUAL, ">=", "null", lineNumber}
				i++
			} else {
				tokens <- token.Struct{token.GREATER, ">", "null", lineNumber}
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
				tokens <- token.Struct{token.SLASH, "/", "null", lineNumber}
			}
		case '.':
			tokens <- token.Struct{token.DOT, ".", "null", lineNumber}
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
			errChar := string(line[i])
			errLinePrefix := fmt.Sprintf("[line %d] Error:", lineNumber)
			errStr := fmt.Sprintf("%s Unexpected character: %s\n", errLinePrefix, errChar)
			err = errors.New(errStr)
			)
		}
		if err != nil {
			tokens <- token.Struct{token.ERROR, "", err.Error(), lineNumber}
		}
		// if err != nil {
		// 	break
		// }
	}
	tokens <- token.Struct{token.EOF, "", "null", lineNumber}
	close(tokens)
}
