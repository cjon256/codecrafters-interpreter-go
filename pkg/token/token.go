package token

import (
	"fmt"
)

type Struct struct {
	Type    Type
	Lexeme  string
	Literal string
	Line    int
}

func (t Struct) String() string {
	return fmt.Sprintf("%s %s %s", t.Type, t.Lexeme, t.Literal)
}

//go:generate stringer -type=Type
type Type int

const (
	EOF Type = iota
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
	AND
	CLASS
	ELSE
	FALSE
	FOR
	FUN
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE
	ERROR
)

var KEYWORDS map[string]Type = map[string]Type{"and": AND, "class": CLASS, "else": ELSE, "false": FALSE, "for": FOR, "fun": FUN, "if": IF, "nil": NIL, "or": OR, "print": PRINT, "return": RETURN, "super": SUPER, "this": THIS, "true": TRUE, "var": VAR, "while": WHILE, "null": EOF}
