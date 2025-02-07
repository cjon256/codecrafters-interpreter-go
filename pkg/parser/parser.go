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
	var peek token.Struct
	var group func() (ASTnode, error)
	var primary func(token.Struct) (ASTnode, error)
	var expression func(token.Struct) (ASTnode, error)
	var equality func(token.Struct) (ASTnode, error)
	var comparison func(token.Struct) (ASTnode, error)
	var term func(token.Struct) (ASTnode, error)
	var factor func(token.Struct) (ASTnode, error)
	var unary func(token.Struct) (ASTnode, error)

	// expression     → equality ;
	expression = func(t token.Struct) (ASTnode, error) {
		expr, err := equality(t)
		return expr, err
	}

	// equality       → comparison ( ( "!=" | "==" ) comparison )* ;
	equality = func(t token.Struct) (ASTnode, error) {
		comp, err := comparison(t)
		return comp, err
	}

	// comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
	comparison = func(t token.Struct) (ASTnode, error) {
		trm, err := term(t)
		return trm, err
	}

	// term           → factor ( ( "-" | "+" ) factor )* ;
	term = func(t token.Struct) (ASTnode, error) {
		fact, err := factor(t)
		return fact, err
	}

	// factor         → unary ( ( "/" | "*" ) unary )* ;
	factor = func(t token.Struct) (ASTnode, error) {
		una, err := unary(t)
		return una, err
	}

	// unary          → ( "!" | "-" ) unary
	//                | primary ;
	unary = func(t token.Struct) (ASTnode, error) {
		prim, err := primary(t)
		return prim, err
	}

	group = func() (ASTnode, error) {
		g := ASTgroup{}

		// XXX hack to let me keep peek around for now
		peek = <-tokens
		c := peek
		switch c.Type {
		case token.EOF:
			return g, errors.New("parse_error")
		case token.RIGHT_PAREN:
			return g, errors.New("parse_error")
		default:
			node, err := primary(c)
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

	// primary        → NUMBER | STRING | "true" | "false" | "nil"
	//                | "(" expression ")" ;
	primary = func(t token.Struct) (ASTnode, error) {
		switch t.Type {
		case token.EOF:
			return ASTliteral{}, errors.New("EOF")
		case token.LEFT_PAREN:
			node, err := group()
			if err != nil {
				return ASTliteral{}, err
			}
			return node, nil
		case token.RIGHT_PAREN:
			return ASTliteral{}, errors.New("parse_error: unexpected ')' in input")
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
			return ASTliteral{}, errors.New("parse_error: upexpected character")
		}
	}

	lastStr := ""
	for t := range tokens {
		node, err := expression(t)
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
