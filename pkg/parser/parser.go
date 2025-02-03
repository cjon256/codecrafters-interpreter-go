package parser

import (
	"fmt"

	"example.com/cjon/tokenizer"
)

func Parse(tokens chan tokenizer.TokenStruct) error {
	for t := range tokens {
		if t.Literal == "null" {
			fmt.Println(t.Lexeme)
		} else {
			fmt.Println(t.Literal)
		}
	}
	return nil
}
