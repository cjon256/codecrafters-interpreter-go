package parser

import (
	"fmt"

	"example.com/cjon/tokenizer"
)

func Parse(tokens chan tokenizer.TokenStruct) error {
	for t := range tokens {
		fmt.Println(t.Lexeme)
	}
	return nil
}
