package parser

import (
	"errors"
	"fmt"

	"example.com/cjon/token"
)

type ASTnode interface {
	fmt.Stringer
}

type ASTgroup struct {
	Contents ASTnode
}

func (g ASTgroup) String() string {
	return fmt.Sprintf("(group %s)", g.Contents)
}

type ASTliteral struct {
	Contents string
}

func (l ASTliteral) String() string {
	return l.Contents
}

func ParseLines(tokens chan token.Struct) error {
	return nil
}

func parseOne() {
}

func Parse(tokens chan token.Struct) error {
	var group func() (ASTnode, error)
	var literal func(token.Struct) (ASTnode, error)

	group = func() (ASTnode, error) {
		g := ASTgroup{}
		c := <-tokens
		switch c.Type {
		case token.EOF:
			return g, errors.New("parse_error")
		case token.RIGHT_PAREN:
			return g, errors.New("parse_error")
		default:
			node, err := literal(c)
			if err != nil {
				return g, err
			}
			g.Contents = node
		}

		close := <-tokens
		if close.Type != token.RIGHT_PAREN {
			return g, errors.New("parse_error")
		}
		return g, nil
	}

	literal = func(t token.Struct) (ASTnode, error) {
		switch t.Type {
		case token.EOF:
			return ASTliteral{}, errors.New("EOF")
		case token.RIGHT_PAREN:
			return ASTliteral{}, errors.New("parse_error")
		case token.LEFT_PAREN:
			node, err := group()
			if err != nil {
				return ASTliteral{}, err
			}
			return node, nil
		case token.STRING:
			return ASTliteral{t.Literal}, nil
		case token.NUMBER:
			return ASTliteral{t.Literal}, nil
		case token.IDENTIFIER:
			return ASTliteral{t.Lexeme}, nil
		case token.TRUE:
			return ASTliteral{t.Lexeme}, nil
		case token.FALSE:
			return ASTliteral{t.Lexeme}, nil
		case token.NIL:
			return ASTliteral{t.Lexeme}, nil
		default:
			fmt.Print("errorrrrrr")
			return ASTliteral{}, errors.New("parse_error")
		}
	}

	lastStr := ""
	for t := range tokens {
		node, err := literal(t)
		if err != nil {
			switch err.Error() {
			case "EOF":
				fmt.Print(lastStr)
				return nil
			default:
				return err
			}
		}
		if lastStr != "" {
			fmt.Print(lastStr, " ")
		}
		lastStr = node.String()
	}
	return nil
}
