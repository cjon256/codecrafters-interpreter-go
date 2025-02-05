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
	group = func() (ASTnode, error) {
		g := ASTgroup{}
		c := <-tokens
		switch c.Type {
		case token.EOF:
			return g, errors.New("parse_error")
		case token.RIGHT_PAREN:
			return g, errors.New("parse_error")
		case token.LEFT_PAREN:
			var err error
			g.Contents, err = group()
			if err != nil {
			}
		case token.STRING:
			g.Contents = ASTliteral{c.Literal}
		case token.NUMBER:
			g.Contents = ASTliteral{c.Literal}
		case token.IDENTIFIER:
			g.Contents = ASTliteral{c.Literal}
		case token.TRUE:
			g.Contents = ASTliteral{c.Literal}
		case token.FALSE:
			g.Contents = ASTliteral{c.Literal}
		case token.NIL:
			g.Contents = ASTliteral{c.Literal}

		default:
			return g, errors.New("parse_error")
		}

		close := <-tokens
		if close.Type != token.LEFT_PAREN {
			return g, errors.New("parse_error")
		}
		return g, nil
	}
	for t := range tokens {
		switch t.Type {
		case token.EOF:
			continue
		case token.RIGHT_PAREN:
			return errors.New("parse_error")
		case token.LEFT_PAREN:
			node, err := group()
			if err != nil {
				return err
			}
		case token.STRING:
			fmt.Print(t.Literal)
		case token.NUMBER:
			fmt.Print(t.Literal)
		case token.IDENTIFIER:
			fmt.Print(t.Lexeme)
		case token.TRUE:
			fmt.Print(t.Lexeme)
		case token.FALSE:
			fmt.Print(t.Lexeme)
		case token.NIL:
			fmt.Print(t.Lexeme)
		default:
			return errors.New("parse_error")
		}
	}
	return nil
}
