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

type lookaheadTokenStream struct {
	ch   chan token.Struct
	curr *token.Struct
}

func (lts *lookaheadTokenStream) peek() *token.Struct {
	if lts.curr == nil {
		t := <-lts.ch
		lts.curr = &t
	}
	return lts.curr
}

func (lts *lookaheadTokenStream) consume() *token.Struct {
	r := lts.peek()
	if lts.curr.Type != token.EOF {
		t := <-lts.ch
		lts.curr = &t
	}
	return r
}

func Parse(tokens chan token.Struct) error {
	lts := lookaheadTokenStream{ch: tokens}
	return parseWithLookahead(lts)
}

func parseWithLookahead(lts lookaheadTokenStream) error {
	var group func() (ASTnode, error)
	var primary func() (ASTnode, error)
	var expression func() (ASTnode, error)
	var equality func() (ASTnode, error)
	var comparison func() (ASTnode, error)
	var term func() (ASTnode, error)
	var factor func() (ASTnode, error)
	var unary func() (ASTnode, error)

	// expression     → equality ;
	expression = func() (ASTnode, error) {
		expr, err := equality()
		return expr, err
	}

	// equality       → comparison ( ( "!=" | "==" ) comparison )* ;
	equality = func() (ASTnode, error) {
		comp, err := comparison()
		return comp, err
	}

	// comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
	comparison = func() (ASTnode, error) {
		trm, err := term()
		return trm, err
	}

	// term           → factor ( ( "-" | "+" ) factor )* ;
	term = func() (ASTnode, error) {
		fact, err := factor()
		n := lts.peek()
		if n.Type == token.MINUS {
			fmt.Println("seeing a minus sign")
		}
		return fact, err
	}

	// factor         → unary ( ( "/" | "*" ) unary )* ;
	factor = func() (ASTnode, error) {
		una, err := unary()
		return una, err
	}

	// unary          → ( "!" | "-" ) unary
	//                | primary ;
	unary = func() (ASTnode, error) {
		t := lts.peek()
		if t.Type == token.BANG {
			b := lts.consume()
			fmt.Println(b.String())
		}

		prim, err := primary()
		return prim, err
	}

	group = func() (ASTnode, error) {
		g := ASTgroup{}

		c := lts.peek()
		switch c.Type {
		case token.EOF:
			return g, errors.New("parse_error: EOF detected in group")
		case token.RIGHT_PAREN:
			return g, errors.New("parse_error: ')' detected in group")
		default:
			node, err := expression()
			if err != nil {
				return g, err
			}
			g.Contents = node
		}

		close := lts.consume()
		if close.Type != token.RIGHT_PAREN {
			return g, errors.New("parse_error: expected ')' in group")
		}
		return g, nil
	}

	// primary        → NUMBER | STRING | "true" | "false" | "nil"
	//                | "(" expression ")" ;
	primary = func() (ASTnode, error) {
		t := lts.consume()
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
			return ASTliteral{}, fmt.Errorf("primary: unexpected character '%v' in input", t.Type)
		}
	}

	lastStr := ""
	initial := true
	for lts.peek().Type != token.EOF {
		node, err := expression()
		if err != nil {
			switch err.Error() {
			case "EOF":
				fmt.Print(lastStr)
				return nil
			default:
				return err
			}
		}
		if initial {
			initial = false
		} else {
			fmt.Print(" ")
		}

		fmt.Print(node)
	}
	return nil
}
