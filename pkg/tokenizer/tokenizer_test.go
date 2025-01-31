package tokenizer

import (
	"fmt"
	"io"
	"os"
	"testing"
)

type testStruct struct {
	name   string
	lines  string
	errors string
	output string
	retval int
}

var tests []testStruct = []testStruct{
	// 	{
	// 		name: "",
	// 		lines: `
	// `,
	// 	errors: `
	// `,
	// 		output: `
	// `,
	// 		retval: 0,
	// 	},
	{
		name: "[UE7] [test-4]",
		lines: `("world"+"hello") != "other_string"
`,
		errors: ``,
		output: `LEFT_PAREN ( null
STRING "world" world
PLUS + null
STRING "hello" hello
RIGHT_PAREN ) null
BANG_EQUAL != null
STRING "other_string" other_string
EOF  null
`,
		retval: 0,
	},
	{
		name:  "String unterminated",
		lines: `"baz" "unterminated`,
		errors: `[line 1] Error: Unterminated string.
`,
		output: `STRING "baz" baz
EOF  null
`,
		retval: 65,
	},
	{
		name: "String works",
		lines: `
"foo baz"
`,
		errors: ``,
		output: `STRING "foo baz" foo baz
EOF  null
`,
		retval: 0,
	},
	{
		name: "[TZ7] [test-1]",
		lines: `
() 	@`,
		errors: `[line 2] Error: Unexpected character: @
`,
		output: `LEFT_PAREN ( null
RIGHT_PAREN ) null
EOF  null
`,
		retval: 65,
	},

	{
		name: "[TZ7] [test-2]",
		lines: ` 	
@
`,
		errors: `[line 2] Error: Unexpected character: @
`,
		output: `EOF  null
`,
		retval: 65,
	},

	{
		name: "[TZ7] [test-3]",
		lines: `()  #	{}
@
$
+++
// Let's Go!
+++
#
`,
		errors: `[line 1] Error: Unexpected character: #
[line 2] Error: Unexpected character: @
[line 3] Error: Unexpected character: $
[line 7] Error: Unexpected character: #
`,
		output: `LEFT_PAREN ( null
RIGHT_PAREN ) null
LEFT_BRACE { null
RIGHT_BRACE } null
PLUS + null
PLUS + null
PLUS + null
PLUS + null
PLUS + null
PLUS + null
EOF  null
`,
		retval: 65,
	},

	{
		name: "[TZ7] [test-4]",
		lines: `({*	%})
`,
		errors: `[line 1] Error: Unexpected character: %
`,
		output: `LEFT_PAREN ( null
LEFT_BRACE { null
STAR * null
RIGHT_BRACE } null
RIGHT_PAREN ) null
EOF  null
`,
		retval: 65,
	},

	{
		name: "[R2] [test-1]",
		lines: ` 
`,
		errors: ``,
		output: `EOF  null
`,
		retval: 0,
	},

	{
		name: "[ER2] [test-2]",
		lines: ` 	
 
`,
		errors: ``,
		output: `EOF  null
`,
		retval: 0,
	},

	{
		name: "[ER2] [test-3]",
		lines: `{
	}
((	-.*,))
`,
		errors: ``,
		output: `LEFT_BRACE { null
RIGHT_BRACE } null
LEFT_PAREN ( null
LEFT_PAREN ( null
MINUS - null
DOT . null
STAR * null
COMMA , null
RIGHT_PAREN ) null
RIGHT_PAREN ) null
EOF  null
`,
		retval: 0,
	},

	{
		name: "[ER2] [test-4]",
		lines: `{ 		

}
((- >	))
`,
		errors: ``,
		output: `LEFT_BRACE { null
RIGHT_BRACE } null
LEFT_PAREN ( null
LEFT_PAREN ( null
MINUS - null
GREATER > null
RIGHT_PAREN ) null
RIGHT_PAREN ) null
EOF  null
`,
		retval: 0,
	},

	{
		name: "[ML2] [test-1]",
		lines: `//Comment
`,
		errors: ``,
		output: `EOF  null
`,
		retval: 0,
	},

	{
		name: "[ML2] [test-2]",
		lines: `(///Unicode:£§᯽☺♣)
`,
		errors: ``,
		output: `LEFT_PAREN ( null
EOF  null
`,
		retval: 0,
	},

	{
		name: "[ML2] [test-3]",
		lines: `/
`,
		errors: ``,
		output: `SLASH / null
EOF  null
`,
		retval: 0,
	},

	{
		name: "[ML2] [test-4]",
		lines: `({(===!)})//Comment
`,
		errors: ``,
		output: `LEFT_PAREN ( null
LEFT_BRACE { null
LEFT_PAREN ( null
EQUAL_EQUAL == null
EQUAL = null
BANG ! null
RIGHT_PAREN ) null
RIGHT_BRACE } null
RIGHT_PAREN ) null
EOF  null
`,
		retval: 0,
	},

	{
		name: "[ET2] [test-1]",
		lines: `>=
`,
		errors: ``,
		output: `GREATER_EQUAL >= null
EOF  null
`,
		retval: 0,
	},

	{
		name: "[ET2] [test-2]",
		lines: `<<<=>>>=
`,
		errors: ``,
		output: `LESS < null
LESS < null
LESS_EQUAL <= null
GREATER > null
GREATER > null
GREATER_EQUAL >= null
EOF  null
`,
		retval: 0,
	},

	{
		name: "[ET2] [test-3]",
		lines: `>=<>=>>
`,
		errors: ``,
		output: `GREATER_EQUAL >= null
LESS < null
GREATER_EQUAL >= null
GREATER > null
GREATER > null
EOF  null
`,
		retval: 0,
	},

	{
		name: "[ET2] [test-4]",
		lines: `(){>!>=}
`,
		errors: ``,
		output: `LEFT_PAREN ( null
RIGHT_PAREN ) null
LEFT_BRACE { null
GREATER > null
BANG ! null
GREATER_EQUAL >= null
RIGHT_BRACE } null
EOF  null
`,
		retval: 0,
	},

	{
		name: "[BU3] [test-1]",
		lines: `!=
`,
		errors: ``,
		output: `BANG_EQUAL != null
EOF  null
`,
		retval: 0,
	},

	{
		name: "[BU3] [test-2]",
		lines: `!!===
`,
		errors: ``,
		output: `BANG ! null
BANG_EQUAL != null
EQUAL_EQUAL == null
EOF  null
`,
		retval: 0,
	},

	{
		name: "[BU3] [test-3]",
		lines: `!{!}(!===)=
`,
		errors: ``,
		output: `BANG ! null
LEFT_BRACE { null
BANG ! null
RIGHT_BRACE } null
LEFT_PAREN ( null
BANG_EQUAL != null
EQUAL_EQUAL == null
RIGHT_PAREN ) null
EQUAL = null
EOF  null
`,
		retval: 0,
	},

	{
		name: "[BU3] [test-4]",
		lines: `{(%=#!==)}
`,
		errors: `[line 1] Error: Unexpected character: %
[line 1] Error: Unexpected character: #
`,
		output: `LEFT_BRACE { null
LEFT_PAREN ( null
EQUAL = null
BANG_EQUAL != null
EQUAL = null
RIGHT_PAREN ) null
RIGHT_BRACE } null
EOF  null
`,
		retval: 65,
	},

	{
		name: "[MP7] [test-1]",
		lines: `=
`,
		errors: ``,
		output: `EQUAL = null
EOF  null
`,
		retval: 0,
	},

	{
		name: "[MP7] [test-2]",
		lines: `==
`,
		errors: ``,
		output: `EQUAL_EQUAL == null
EOF  null
`,
		retval: 0,
	},

	{
		name: "[MP7] [test-3]",
		lines: `({=}){==}
`,
		errors: ``,
		output: `LEFT_PAREN ( null
LEFT_BRACE { null
EQUAL = null
RIGHT_BRACE } null
RIGHT_PAREN ) null
LEFT_BRACE { null
EQUAL_EQUAL == null
RIGHT_BRACE } null
EOF  null
`,
		retval: 0,
	},

	{
		name: "[MP7] [test-4]",
		lines: `(($#===%))
`,
		errors: `[line 1] Error: Unexpected character: $
[line 1] Error: Unexpected character: #
[line 1] Error: Unexpected character: %
`,
		output: `LEFT_PAREN ( null
LEFT_PAREN ( null
EQUAL_EQUAL == null
EQUAL = null
RIGHT_PAREN ) null
RIGHT_PAREN ) null
EOF  null
`,
		retval: 65,
	},

	{
		name: "[EA6] [test-1]",
		lines: `@
`,
		errors: `[line 1] Error: Unexpected character: @
`,
		output: `EOF  null
`,
		retval: 65,
	},

	{
		name: "[EA6] [test-2]",
		lines: `,.$(#
`,
		errors: `[line 1] Error: Unexpected character: $
[line 1] Error: Unexpected character: #
`,
		output: `COMMA , null
DOT . null
LEFT_PAREN ( null
EOF  null
`,
		retval: 65,
	},

	{
		name: "[EA6] [test-3]",
		lines: `%$#@%
`,
		errors: `[line 1] Error: Unexpected character: %
[line 1] Error: Unexpected character: $
[line 1] Error: Unexpected character: #
[line 1] Error: Unexpected character: @
[line 1] Error: Unexpected character: %
`,
		output: `EOF  null
`,
		retval: 65,
	},

	{
		name: "[EA6] [test-4]",
		lines: `{(*-#;.$%)}
`,
		errors: `[line 1] Error: Unexpected character: #
[line 1] Error: Unexpected character: $
[line 1] Error: Unexpected character: %
`,
		output: `LEFT_BRACE { null
LEFT_PAREN ( null
STAR * null
MINUS - null
SEMICOLON ; null
DOT . null
RIGHT_PAREN ) null
RIGHT_BRACE } null
EOF  null
`,
		retval: 65,
	},

	{
		name: "[XC5] [test-1]",
		lines: `+-
`,
		errors: ``,
		output: `PLUS + null
MINUS - null
EOF  null
`,
		retval: 0,
	},

	{
		name: "[XC5] [test-2]",
		lines: `++--**..,,;;
`,
		errors: ``,
		output: `PLUS + null
PLUS + null
MINUS - null
MINUS - null
STAR * null
STAR * null
DOT . null
DOT . null
COMMA , null
COMMA , null
SEMICOLON ; null
SEMICOLON ; null
EOF  null
`,
		retval: 0,
	},

	{
		name: "[XC5] [test-3]",
		lines: `+;*.-,*
`,
		errors: ``,
		output: `PLUS + null
SEMICOLON ; null
STAR * null
DOT . null
MINUS - null
COMMA , null
STAR * null
EOF  null
`,
		retval: 0,
	},

	{
		name: "[XC5] [test-4]",
		lines: `({.,;*+})
`,
		errors: ``,
		output: `LEFT_PAREN ( null
LEFT_BRACE { null
DOT . null
COMMA , null
SEMICOLON ; null
STAR * null
PLUS + null
RIGHT_BRACE } null
RIGHT_PAREN ) null
EOF  null
`,
		retval: 0,
	},

	{
		name: "[OE8] [test-1]",
		lines: `}
`,
		errors: ``,
		output: `RIGHT_BRACE } null
EOF  null
`,
		retval: 0,
	},

	{
		name: "[OE8] [test-2]",
		lines: `{{}}
`,
		errors: ``,
		output: `LEFT_BRACE { null
LEFT_BRACE { null
RIGHT_BRACE } null
RIGHT_BRACE } null
EOF  null
`,
		retval: 0,
	},

	{
		name: "[OE8] [test-3]",
		lines: `{}{{}
`,
		errors: ``,
		output: `LEFT_BRACE { null
RIGHT_BRACE } null
LEFT_BRACE { null
LEFT_BRACE { null
RIGHT_BRACE } null
EOF  null
`,
		retval: 0,
	},

	{
		name: "[OE8] [test-4]",
		lines: `}){{)}(
`,
		errors: ``,
		output: `RIGHT_BRACE } null
RIGHT_PAREN ) null
LEFT_BRACE { null
LEFT_BRACE { null
RIGHT_PAREN ) null
RIGHT_BRACE } null
LEFT_PAREN ( null
EOF  null
`,
		retval: 0,
	},

	{
		name: "[OL4] [test-1]",
		lines: `(
`,
		errors: ``,
		output: `LEFT_PAREN ( null
EOF  null
`,
		retval: 0,
	},

	{
		name: "[OL4] [test-2]",
		lines: `))
`,
		errors: ``,
		output: `RIGHT_PAREN ) null
RIGHT_PAREN ) null
EOF  null
`,
		retval: 0,
	},

	{
		name: "[OL4] [test-3]",
		lines: `())()
`,
		errors: ``,
		output: `LEFT_PAREN ( null
RIGHT_PAREN ) null
RIGHT_PAREN ) null
LEFT_PAREN ( null
RIGHT_PAREN ) null
EOF  null
`,
		retval: 0,
	},

	{
		name: "[OL4] [test-4]",
		lines: `(())))(
`,
		errors: ``,
		output: `LEFT_PAREN ( null
LEFT_PAREN ( null
RIGHT_PAREN ) null
RIGHT_PAREN ) null
RIGHT_PAREN ) null
RIGHT_PAREN ) null
LEFT_PAREN ( null
EOF  null
`,
		retval: 0,
	},

	{
		name:   "[RY8] [test-1]",
		lines:  ``,
		errors: ``,
		output: `EOF  null
`,
		retval: 0,
	},
}

func doTest(lines []byte) (string, string, int) {
	output := ""
	retval := 0

	fmt.Printf("lines is %v\n", lines)

	errCh := make(chan error)
	tokCh := make(chan TokenStruct)

	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w
	go Tokenize(tokCh, errCh, lines)
	for t := range tokCh {
		output = output + fmt.Sprintln(t)
	}
	os.Stderr = oldStderr
	w.Close()

	syntaxErrs, _ := io.ReadAll(r)

	err := <-errCh
	if err != nil {
		retval = 65
	}
	return output, string(syntaxErrs), retval
}

func TestTokenize(t *testing.T) {
	for _, test := range tests {
		output, errs, retval := doTest([]byte(test.lines))
		if test.output != output {
			t.Errorf("%s: output does not match:\n\texpected '%#v'\n\tgot    : '%#v'\n", test.name, test.output, output)
		}
		if test.errors != errs {
			t.Errorf("%s: errors does not match:\n\texpected '%#v'\n\tgot    : '%#v'\n", test.name, test.errors, errs)
		}
		if test.retval != retval {
			t.Errorf("%s: retval does not match:\n\texpected %d\n\tgot       : %d\n", test.name, test.retval, retval)
		}
	}
}
