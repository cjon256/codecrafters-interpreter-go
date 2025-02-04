package parser

import (
	"errors"
	"fmt"

	"example.com/cjon/tokenizer"
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

func ParseLines(tokens chan tokenizer.TokenStruct) error {
}

func Parse(tokens chan tokenizer.TokenStruct) error {
	var group func() (ASTnode, error)
	group = func() (ASTnode, error) {
		g := ASTgroup{}
		c := <-tokens
		switch c.Type {
		case tokenizer.EOF:
			return g, errors.New("parse_error")
		case tokenizer.RIGHT_PAREN:
			return g, errors.New("parse_error")
		case tokenizer.LEFT_PAREN:
			var err error
			g.Contents, err = group()
			if err != nil {
			}
		case tokenizer.STRING:
			g.Contents = ASTliteral{c.Literal}
		case tokenizer.NUMBER:
			g.Contents = ASTliteral{c.Literal}
		case tokenizer.IDENTIFIER:
			g.Contents = ASTliteral{c.Literal}
		case tokenizer.TRUE:
			g.Contents = ASTliteral{c.Literal}
		case tokenizer.FALSE:
			g.Contents = ASTliteral{c.Literal}
		case tokenizer.NIL:
			g.Contents = ASTliteral{c.Literal}

		default:
			return g, errors.New("parse_error")
		}

		close := <-tokens
		if close.Type != tokenizer.LEFT_PAREN {
			return g, errors.New("parse_error")
		}
		return g, nil
	}
	for t := range tokens {
		switch t.Type {
		case tokenizer.EOF:
			continue
		case tokenizer.RIGHT_PAREN:
			return errors.New("parse_error")
		case tokenizer.LEFT_PAREN:
			err := group()
			if err != nil {
				return err
			}
		case tokenizer.STRING:
			fmt.Print(t.Literal)
		case tokenizer.NUMBER:
			fmt.Print(t.Literal)
		case tokenizer.IDENTIFIER:
			fmt.Print(t.Lexeme)
		case tokenizer.TRUE:
			fmt.Print(t.Lexeme)
		case tokenizer.FALSE:
			fmt.Print(t.Lexeme)
		case tokenizer.NIL:
			fmt.Print(t.Lexeme)
		default:
			return errors.New("parse_error")
		}
	}
	return nil
}
