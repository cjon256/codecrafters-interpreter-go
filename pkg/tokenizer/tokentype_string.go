// Code generated by "stringer -type=TokenType"; DO NOT EDIT.

package tokenizer

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[EOF-0]
	_ = x[LEFT_PAREN-1]
	_ = x[RIGHT_PAREN-2]
	_ = x[LEFT_BRACE-3]
	_ = x[RIGHT_BRACE-4]
	_ = x[SEMICOLON-5]
	_ = x[COMMA-6]
	_ = x[PLUS-7]
	_ = x[MINUS-8]
	_ = x[EQUAL-9]
	_ = x[STAR-10]
	_ = x[BANG_EQUAL-11]
	_ = x[EQUAL_EQUAL-12]
	_ = x[LESS_EQUAL-13]
	_ = x[GREATER_EQUAL-14]
	_ = x[LESS-15]
	_ = x[GREATER-16]
	_ = x[SLASH-17]
	_ = x[DOT-18]
	_ = x[BANG-19]
	_ = x[STRING-20]
	_ = x[NUMBER-21]
}

const _TokenType_name = "EOFLEFT_PARENRIGHT_PARENLEFT_BRACERIGHT_BRACESEMICOLONCOMMAPLUSMINUSEQUALSTARBANG_EQUALEQUAL_EQUALLESS_EQUALGREATER_EQUALLESSGREATERSLASHDOTBANGSTRINGNUMBER"

var _TokenType_index = [...]uint8{0, 3, 13, 24, 34, 45, 54, 59, 63, 68, 73, 77, 87, 98, 108, 121, 125, 132, 137, 140, 144, 150, 156}

func (i TokenType) String() string {
	if i < 0 || i >= TokenType(len(_TokenType_index)-1) {
		return "TokenType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TokenType_name[_TokenType_index[i]:_TokenType_index[i+1]]
}
