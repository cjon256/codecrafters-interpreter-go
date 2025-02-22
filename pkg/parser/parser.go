package parser

import (
	"errors"
	"fmt"

	"example.com/cjon/interpreter-starter-go/pkg/token"
)

type ASTnode interface {
	fmt.Stringer
	Evaluate() string
}

type ASTerror struct {
	err error
}

func (e ASTerror) String() string {
	return e.err.Error()
}

func (e ASTerror) Evaluate() string {
	return e.err.Error()
}

type ASTgroup struct {
	Contents ASTnode
}

func (g ASTgroup) String() string {
	return fmt.Sprintf("(group %s)", g.Contents)
}

func (g ASTgroup) Evaluate() string {
	return ""
}

type ASTliteral struct {
	Contents string
}

func (l ASTliteral) String() string {
	return l.Contents
}

func (l ASTliteral) Evaluate() string {
	return l.Contents
}

type ASTunary struct {
	Operator token.Type
	Contents ASTnode
}

func (l ASTunary) String() string {
	switch l.Operator {
	case token.BANG:
		str := fmt.Sprintf("(! %s)", l.Contents)
		return str
	case token.MINUS:
		str := fmt.Sprintf("(- %s)", l.Contents)
		return str
	default:
		return l.Contents.String()
	}
}

func (l ASTunary) Evaluate() string {
	return ""
}

type ASTbinary struct {
	Operator token.Type
	Left     ASTnode
	Right    ASTnode
}

func (b ASTbinary) String() string {
	var str string
	switch b.Operator {
	case token.SLASH:
		str = fmt.Sprintf("(/ %s %s)", b.Left, b.Right)
	case token.STAR:
		str = fmt.Sprintf("(* %s %s)", b.Left, b.Right)
	case token.PLUS:
		str = fmt.Sprintf("(+ %s %s)", b.Left, b.Right)
	case token.MINUS:
		str = fmt.Sprintf("(- %s %s)", b.Left, b.Right)
	case token.BANG_EQUAL:
		str = fmt.Sprintf("(!= %s %s)", b.Left, b.Right)
	case token.EQUAL_EQUAL:
		str = fmt.Sprintf("(== %s %s)", b.Left, b.Right)
	case token.GREATER:
		str = fmt.Sprintf("(> %s %s)", b.Left, b.Right)
	case token.GREATER_EQUAL:
		str = fmt.Sprintf("(>= %s %s)", b.Left, b.Right)
	case token.LESS:
		str = fmt.Sprintf("(< %s %s)", b.Left, b.Right)
	case token.LESS_EQUAL:
		str = fmt.Sprintf("(<= %s %s)", b.Left, b.Right)
	default:
		str = fmt.Sprintf("(?? %s %s)", b.Left, b.Right)
	}
	return str
}

func (b ASTbinary) Evaluate() string {
	return ""
}

type lookaheadTokenStream struct {
	ch   <-chan token.Struct
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

func Parse(tokens <-chan token.Struct, astNodes chan<- ASTnode) {
	var group func() (ASTnode, error)
	var primary func() (ASTnode, error)
	var expression func() (ASTnode, error)
	var equality func() (ASTnode, error)
	var comparison func() (ASTnode, error)
	var term func() (ASTnode, error)
	var factor func() (ASTnode, error)
	var unary func() (ASTnode, error)
	lts := lookaheadTokenStream{ch: tokens}

	// expression     → equality ;
	expression = func() (ASTnode, error) {
		expr, err := equality()
		return expr, err
	}

	// equality       → comparison ( ( "!=" | "==" ) comparison )* ;
	equality = func() (ASTnode, error) {
		left, err := comparison()
		if err != nil {
			return left, err
		}

	done:
		for {
			o := lts.peek()
			switch o.Type {
			case token.BANG_EQUAL:
				fallthrough
			case token.EQUAL_EQUAL:
				o = lts.consume()
				var right ASTnode
				right, err = comparison()
				tmp := ASTbinary{o.Type, left, right}
				if err != nil {
					return tmp, err
				}
				left = tmp
			default:
				break done
			}
		}

		return left, err
	}

	// comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
	comparison = func() (ASTnode, error) {
		left, err := term()
		if err != nil {
			return left, err
		}

	done:
		for {
			o := lts.peek()
			switch o.Type {
			case token.GREATER:
				fallthrough
			case token.GREATER_EQUAL:
				fallthrough
			case token.LESS:
				fallthrough
			case token.LESS_EQUAL:
				o = lts.consume()
				var right ASTnode
				right, err = term()
				tmp := ASTbinary{o.Type, left, right}
				if err != nil {
					return tmp, err
				}
				left = tmp
			default:
				break done
			}
		}

		return left, err
	}

	// term           → factor ( ( "-" | "+" ) factor )* ;
	term = func() (ASTnode, error) {
		left, err := factor()
		if err != nil {
			return left, err
		}

	done:
		for {
			o := lts.peek()
			switch o.Type {
			case token.PLUS:
				fallthrough
			case token.MINUS:
				o = lts.consume()
				var right ASTnode
				right, err = factor()
				tmp := ASTbinary{o.Type, left, right}
				if err != nil {
					return tmp, err
				}
				left = tmp
			default:
				break done
			}
		}
		return left, err
	}

	// factor         → unary ( ( "/" | "*" ) unary )* ;
	factor = func() (ASTnode, error) {
		left, err := unary()
		if err != nil {
			return left, err
		}

	done:
		for {
			o := lts.peek()
			switch o.Type {
			case token.STAR:
				fallthrough
			case token.SLASH:
				o = lts.consume()
				var right ASTnode
				right, err = unary()
				tmp := ASTbinary{o.Type, left, right}
				if err != nil {
					return tmp, err
				}
				left = tmp
			default:
				break done
			}
		}
		return left, err
	}

	// unary          → ( "!" | "-" ) unary | primary ;
	unary = func() (ASTnode, error) {
		t := lts.peek()
		if t.Type == token.BANG || t.Type == token.MINUS {
			t = lts.consume()
			prim, err := unary()
			wrapper := ASTunary{Operator: t.Type, Contents: prim}
			return wrapper, err
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
			return ASTliteral{}, errors.New("parse_error: EOF detected, expected literal")
		case token.ERROR:
			return ASTliteral{}, errors.New(t.Lexeme)
		case token.LEFT_PAREN:
			node, err := group()
			if err != nil {
				return ASTliteral{}, err
			}
			return node, nil
		case token.RIGHT_PAREN:
			return ASTliteral{}, fmt.Errorf("[line %d] Error at ')': Expect expression.", t.Line)
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

	for lts.peek().Type != token.EOF {
		node, err := expression()
		if err != nil {
			node = ASTerror{err}
		}
		astNodes <- node
	}
	close(astNodes)
}
