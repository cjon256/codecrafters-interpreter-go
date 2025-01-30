package main

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
	{
		name:   "#RY8 test-1",
		lines:  "",
		errors: "",
		output: "EOF  null\n",
		retval: 0,
	},
	{
		name:   "TZ7 test-1",
		lines:  "() 	@",
		errors: "[line 1] Error: Unexpected character: @\n",
		output: "LEFT_PAREN ( null\nRIGHT_PAREN ) null\nEOF  null\n",
		retval: 65,
	},
	{
		name:   "TZ7 test-2",
		lines:  " \t\n@",
		errors: "[line 2] Error: Unexpected character: @\n",
		output: "EOF  null\n",
		retval: 65,
	},
	{
		name: "TZ7 test-3",
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
		name:   "#TZ7 test-4",
		lines:  "({*	%})",
		errors: "[line 1] Error: Unexpected character: %\n",
		output: `LEFT_PAREN ( null
LEFT_BRACE { null
STAR * null
RIGHT_BRACE } null
RIGHT_PAREN ) null
EOF  null
`,
		retval: 65,
	},
}

func doTest(lines []byte) (string, string, int) {
	output := ""
	retval := 0

	fmt.Printf("lines is %v\n", lines)

	errCh := make(chan error)
	tokCh := make(chan tokenStruct)

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
			t.Errorf("output does not match:\n\texpected '%#v'\n\tgot   : '%#v'\n", test.output, output)
		}
		if test.errors != errs {
			t.Errorf("errors does not match:\n\texpected '%#v'\n\tgot   : '%#v'\n", test.errors, errs)
		}
		if test.retval != retval {
			t.Errorf("retval does not match:\n\texpected %d\n\tgot   : %d\n", test.retval, retval)
		}
	}
}

/*
[tester::#ER2] Running tests for Stage #ER2 (Scanning: Whitespace)
[tester::#ER2] [test-1] Running test case: 1
[tester::#ER2] [test-1] Writing contents to ./test.lox:
[tester::#ER2] [test-1] [test.lox] <|SPACE|>
[tester::#ER2] [test-1] $ ./your_program.sh tokenize test.lox
[your_program] EOF  null
[tester::#ER2] [test-1] ✓ 1 line(s) match on stdout
[tester::#ER2] [test-1] ✓ Received exit code 0.
[tester::#ER2] [test-2] Running test case: 2
[tester::#ER2] [test-2] Writing contents to ./test.lox:
[tester::#ER2] [test-2] [test.lox]  <|TAB|>
[tester::#ER2] [test-2] [test.lox] <|SPACE|>
[tester::#ER2] [test-2] $ ./your_program.sh tokenize test.lox
[your_program] EOF  null
[tester::#ER2] [test-2] ✓ 1 line(s) match on stdout
[tester::#ER2] [test-2] ✓ Received exit code 0.
[tester::#ER2] [test-3] Running test case: 3
[tester::#ER2] [test-3] Writing contents to ./test.lox:
[tester::#ER2] [test-3] [test.lox] {
[tester::#ER2] [test-3] [test.lox] <|TAB|>}
[tester::#ER2] [test-3] [test.lox] ((<|TAB|>-.*,))
[tester::#ER2] [test-3] $ ./your_program.sh tokenize test.lox
[your_program] LEFT_BRACE { null
[your_program] RIGHT_BRACE } null
[your_program] LEFT_PAREN ( null
[your_program] LEFT_PAREN ( null
[your_program] MINUS - null
[your_program] DOT . null
[your_program] STAR * null
[your_program] COMMA , null
[your_program] RIGHT_PAREN ) null
[your_program] RIGHT_PAREN ) null
[your_program] EOF  null
[tester::#ER2] [test-3] ✓ 11 line(s) match on stdout
[tester::#ER2] [test-3] ✓ Received exit code 0.
[tester::#ER2] [test-4] Running test case: 4
[tester::#ER2] [test-4] Writing contents to ./test.lox:
[tester::#ER2] [test-4] [test.lox] { <|TAB|><|TAB|>
[tester::#ER2] [test-4] [test.lox]
[tester::#ER2] [test-4] [test.lox] }
[tester::#ER2] [test-4] [test.lox] ((-<|SPACE|>><|TAB|>))
[tester::#ER2] [test-4] $ ./your_program.sh tokenize test.lox
[your_program] LEFT_BRACE { null
[your_program] RIGHT_BRACE } null
[your_program] LEFT_PAREN ( null
[your_program] LEFT_PAREN ( null
[your_program] MINUS - null
[your_program] GREATER > null
[your_program] RIGHT_PAREN ) null
[your_program] RIGHT_PAREN ) null
[your_program] EOF  null
[tester::#ER2] [test-4] ✓ 9 line(s) match on stdout
[tester::#ER2] [test-4] ✓ Received exit code 0.
[tester::#ER2] Test passed.

[tester::#ML2] Running tests for Stage #ML2 (Scanning: Division operator & comments)
[tester::#ML2] [test-1] Running test case: 1
[tester::#ML2] [test-1] Writing contents to ./test.lox:
[tester::#ML2] [test-1] [test.lox] //Comment
[tester::#ML2] [test-1] $ ./your_program.sh tokenize test.lox
[your_program] EOF  null
[tester::#ML2] [test-1] ✓ 1 line(s) match on stdout
[tester::#ML2] [test-1] ✓ Received exit code 0.
[tester::#ML2] [test-2] Running test case: 2
[tester::#ML2] [test-2] Writing contents to ./test.lox:
[tester::#ML2] [test-2] [test.lox] (///Unicode:£§᯽☺♣)
[tester::#ML2] [test-2] $ ./your_program.sh tokenize test.lox
[your_program] LEFT_PAREN ( null
[your_program] EOF  null
[tester::#ML2] [test-2] ✓ 2 line(s) match on stdout
[tester::#ML2] [test-2] ✓ Received exit code 0.
[tester::#ML2] [test-3] Running test case: 3
[tester::#ML2] [test-3] Writing contents to ./test.lox:
[tester::#ML2] [test-3] [test.lox] /
[tester::#ML2] [test-3] $ ./your_program.sh tokenize test.lox
[your_program] SLASH / null
[your_program] EOF  null
[tester::#ML2] [test-3] ✓ 2 line(s) match on stdout
[tester::#ML2] [test-3] ✓ Received exit code 0.
[tester::#ML2] [test-4] Running test case: 4
[tester::#ML2] [test-4] Writing contents to ./test.lox:
[tester::#ML2] [test-4] [test.lox] ({(===!)})//Comment
[tester::#ML2] [test-4] $ ./your_program.sh tokenize test.lox
[your_program] LEFT_PAREN ( null
[your_program] LEFT_BRACE { null
[your_program] LEFT_PAREN ( null
[your_program] EQUAL_EQUAL == null
[your_program] EQUAL = null
[your_program] BANG ! null
[your_program] RIGHT_PAREN ) null
[your_program] RIGHT_BRACE } null
[your_program] RIGHT_PAREN ) null
[your_program] EOF  null
[tester::#ML2] [test-4] ✓ 10 line(s) match on stdout
[tester::#ML2] [test-4] ✓ Received exit code 0.
[tester::#ML2] Test passed.

[tester::#ET2] Running tests for Stage #ET2 (Scanning: Relational operators)
[tester::#ET2] [test-1] Running test case: 1
[tester::#ET2] [test-1] Writing contents to ./test.lox:
[tester::#ET2] [test-1] [test.lox] >=
[tester::#ET2] [test-1] $ ./your_program.sh tokenize test.lox
[your_program] GREATER_EQUAL >= null
[your_program] EOF  null
[tester::#ET2] [test-1] ✓ 2 line(s) match on stdout
[tester::#ET2] [test-1] ✓ Received exit code 0.
[tester::#ET2] [test-2] Running test case: 2
[tester::#ET2] [test-2] Writing contents to ./test.lox:
[tester::#ET2] [test-2] [test.lox] <<<=>>>=
[tester::#ET2] [test-2] $ ./your_program.sh tokenize test.lox
[your_program] LESS < null
[your_program] LESS < null
[your_program] LESS_EQUAL <= null
[your_program] GREATER > null
[your_program] GREATER > null
[your_program] GREATER_EQUAL >= null
[your_program] EOF  null
[tester::#ET2] [test-2] ✓ 7 line(s) match on stdout
[tester::#ET2] [test-2] ✓ Received exit code 0.
[tester::#ET2] [test-3] Running test case: 3
[tester::#ET2] [test-3] Writing contents to ./test.lox:
[tester::#ET2] [test-3] [test.lox] >=<>=>>
[tester::#ET2] [test-3] $ ./your_program.sh tokenize test.lox
[your_program] GREATER_EQUAL >= null
[your_program] LESS < null
[your_program] GREATER_EQUAL >= null
[your_program] GREATER > null
[your_program] GREATER > null
[your_program] EOF  null
[tester::#ET2] [test-3] ✓ 6 line(s) match on stdout
[tester::#ET2] [test-3] ✓ Received exit code 0.
[tester::#ET2] [test-4] Running test case: 4
[tester::#ET2] [test-4] Writing contents to ./test.lox:
[tester::#ET2] [test-4] [test.lox] (){>!>=}
[tester::#ET2] [test-4] $ ./your_program.sh tokenize test.lox
[your_program] LEFT_PAREN ( null
[your_program] RIGHT_PAREN ) null
[your_program] LEFT_BRACE { null
[your_program] GREATER > null
[your_program] BANG ! null
[your_program] GREATER_EQUAL >= null
[your_program] RIGHT_BRACE } null
[your_program] EOF  null
[tester::#ET2] [test-4] ✓ 8 line(s) match on stdout
[tester::#ET2] [test-4] ✓ Received exit code 0.
[tester::#ET2] Test passed.

[tester::#BU3] Running tests for Stage #BU3 (Scanning: Negation & inequality operators)
[tester::#BU3] [test-1] Running test case: 1
[tester::#BU3] [test-1] Writing contents to ./test.lox:
[tester::#BU3] [test-1] [test.lox] !=
[tester::#BU3] [test-1] $ ./your_program.sh tokenize test.lox
[your_program] BANG_EQUAL != null
[your_program] EOF  null
[tester::#BU3] [test-1] ✓ 2 line(s) match on stdout
[tester::#BU3] [test-1] ✓ Received exit code 0.
[tester::#BU3] [test-2] Running test case: 2
[tester::#BU3] [test-2] Writing contents to ./test.lox:
[tester::#BU3] [test-2] [test.lox] !!===
[tester::#BU3] [test-2] $ ./your_program.sh tokenize test.lox
[your_program] BANG ! null
[your_program] BANG_EQUAL != null
[your_program] EQUAL_EQUAL == null
[your_program] EOF  null
[tester::#BU3] [test-2] ✓ 4 line(s) match on stdout
[tester::#BU3] [test-2] ✓ Received exit code 0.
[tester::#BU3] [test-3] Running test case: 3
[tester::#BU3] [test-3] Writing contents to ./test.lox:
[tester::#BU3] [test-3] [test.lox] !{!}(!===)=
[tester::#BU3] [test-3] $ ./your_program.sh tokenize test.lox
[your_program] BANG ! null
[your_program] LEFT_BRACE { null
[your_program] BANG ! null
[your_program] RIGHT_BRACE } null
[your_program] LEFT_PAREN ( null
[your_program] BANG_EQUAL != null
[your_program] EQUAL_EQUAL == null
[your_program] RIGHT_PAREN ) null
[your_program] EQUAL = null
[your_program] EOF  null
[tester::#BU3] [test-3] ✓ 10 line(s) match on stdout
[tester::#BU3] [test-3] ✓ Received exit code 0.
[tester::#BU3] [test-4] Running test case: 4
[tester::#BU3] [test-4] Writing contents to ./test.lox:
[tester::#BU3] [test-4] [test.lox] {(%=#!==)}
[tester::#BU3] [test-4] $ ./your_program.sh tokenize test.lox
[your_program] LEFT_BRACE { null
[your_program] LEFT_PAREN ( null
[your_program] [line 1] Error: Unexpected character: %
[your_program] EQUAL = null
[your_program] BANG_EQUAL != null
[your_program] EQUAL = null
[your_program] [line 1] Error: Unexpected character: #
[your_program] RIGHT_PAREN ) null
[your_program] RIGHT_BRACE } null
[your_program] EOF  null
[tester::#BU3] [test-4] ✓ 2 line(s) match on stderr
[tester::#BU3] [test-4] ✓ 8 line(s) match on stdout
[tester::#BU3] [test-4] ✓ Received exit code 65.
[tester::#BU3] Test passed.

[tester::#MP7] Running tests for Stage #MP7 (Scanning: Assignment & equality Operators)
[tester::#MP7] [test-1] Running test case: 1
[tester::#MP7] [test-1] Writing contents to ./test.lox:
[tester::#MP7] [test-1] [test.lox] =
[tester::#MP7] [test-1] $ ./your_program.sh tokenize test.lox
[your_program] EQUAL = null
[your_program] EOF  null
[tester::#MP7] [test-1] ✓ 2 line(s) match on stdout
[tester::#MP7] [test-1] ✓ Received exit code 0.
[tester::#MP7] [test-2] Running test case: 2
[tester::#MP7] [test-2] Writing contents to ./test.lox:
[tester::#MP7] [test-2] [test.lox] ==
[tester::#MP7] [test-2] $ ./your_program.sh tokenize test.lox
[your_program] EQUAL_EQUAL == null
[your_program] EOF  null
[tester::#MP7] [test-2] ✓ 2 line(s) match on stdout
[tester::#MP7] [test-2] ✓ Received exit code 0.
[tester::#MP7] [test-3] Running test case: 3
[tester::#MP7] [test-3] Writing contents to ./test.lox:
[tester::#MP7] [test-3] [test.lox] ({=}){==}
[tester::#MP7] [test-3] $ ./your_program.sh tokenize test.lox
[your_program] LEFT_PAREN ( null
[your_program] LEFT_BRACE { null
[your_program] EQUAL = null
[your_program] RIGHT_BRACE } null
[your_program] RIGHT_PAREN ) null
[your_program] LEFT_BRACE { null
[your_program] EQUAL_EQUAL == null
[your_program] RIGHT_BRACE } null
[your_program] EOF  null
[tester::#MP7] [test-3] ✓ 9 line(s) match on stdout
[tester::#MP7] [test-3] ✓ Received exit code 0.
[tester::#MP7] [test-4] Running test case: 4
[tester::#MP7] [test-4] Writing contents to ./test.lox:
[tester::#MP7] [test-4] [test.lox] (($#===%))
[tester::#MP7] [test-4] $ ./your_program.sh tokenize test.lox
[your_program] LEFT_PAREN ( null
[your_program] LEFT_PAREN ( null
[your_program] EQUAL_EQUAL == null
[your_program] [line 1] Error: Unexpected character: $
[your_program] EQUAL = null
[your_program] RIGHT_PAREN ) null
[your_program] [line 1] Error: Unexpected character: #
[your_program] RIGHT_PAREN ) null
[your_program] EOF  null
[your_program] [line 1] Error: Unexpected character: %
[tester::#MP7] [test-4] ✓ 3 line(s) match on stderr
[tester::#MP7] [test-4] ✓ 7 line(s) match on stdout
[tester::#MP7] [test-4] ✓ Received exit code 65.
[tester::#MP7] Test passed.

[tester::#EA6] Running tests for Stage #EA6 (Scanning: Lexical errors)
[tester::#EA6] [test-1] Running test case: 1
[tester::#EA6] [test-1] Writing contents to ./test.lox:
[tester::#EA6] [test-1] [test.lox] @
[tester::#EA6] [test-1] $ ./your_program.sh tokenize test.lox
[your_program] EOF  null
[your_program] [line 1] Error: Unexpected character: @
[tester::#EA6] [test-1] ✓ 1 line(s) match on stderr
[tester::#EA6] [test-1] ✓ 1 line(s) match on stdout
[tester::#EA6] [test-1] ✓ Received exit code 65.
[tester::#EA6] [test-2] Running test case: 2
[tester::#EA6] [test-2] Writing contents to ./test.lox:
[tester::#EA6] [test-2] [test.lox] ,.$(#
[tester::#EA6] [test-2] $ ./your_program.sh tokenize test.lox
[your_program] [line 1] Error: Unexpected character: $
[your_program] [line 1] Error: Unexpected character: #
[your_program] COMMA , null
[your_program] DOT . null
[your_program] LEFT_PAREN ( null
[your_program] EOF  null
[tester::#EA6] [test-2] ✓ 2 line(s) match on stderr
[tester::#EA6] [test-2] ✓ 4 line(s) match on stdout
[tester::#EA6] [test-2] ✓ Received exit code 65.
[tester::#EA6] [test-3] Running test case: 3
[tester::#EA6] [test-3] Writing contents to ./test.lox:
[tester::#EA6] [test-3] [test.lox] %$#@%
[tester::#EA6] [test-3] $ ./your_program.sh tokenize test.lox
[your_program] [line 1] Error: Unexpected character: %
[your_program] [line 1] Error: Unexpected character: $
[your_program] EOF  null
[your_program] [line 1] Error: Unexpected character: #
[your_program] [line 1] Error: Unexpected character: @
[your_program] [line 1] Error: Unexpected character: %
[tester::#EA6] [test-3] ✓ 5 line(s) match on stderr
[tester::#EA6] [test-3] ✓ 1 line(s) match on stdout
[tester::#EA6] [test-3] ✓ Received exit code 65.
[tester::#EA6] [test-4] Running test case: 4
[tester::#EA6] [test-4] Writing contents to ./test.lox:
[tester::#EA6] [test-4] [test.lox] {(*-#;.$%)}
[tester::#EA6] [test-4] $ ./your_program.sh tokenize test.lox
[your_program] LEFT_BRACE { null
[your_program] LEFT_PAREN ( null
[your_program] STAR * null
[your_program] [line 1] Error: Unexpected character: #
[your_program] [line 1] Error: Unexpected character: $
[your_program] [line 1] Error: Unexpected character: %
[your_program] MINUS - null
[your_program] SEMICOLON ; null
[your_program] DOT . null
[your_program] RIGHT_PAREN ) null
[your_program] RIGHT_BRACE } null
[your_program] EOF  null
[tester::#EA6] [test-4] ✓ 3 line(s) match on stderr
[tester::#EA6] [test-4] ✓ 9 line(s) match on stdout
[tester::#EA6] [test-4] ✓ Received exit code 65.
[tester::#EA6] Test passed.

[tester::#XC5] Running tests for Stage #XC5 (Scanning: Other single-character tokens)
[tester::#XC5] [test-1] Running test case: 1
[tester::#XC5] [test-1] Writing contents to ./test.lox:
[tester::#XC5] [test-1] [test.lox] +-
[tester::#XC5] [test-1] $ ./your_program.sh tokenize test.lox
[your_program] PLUS + null
[your_program] MINUS - null
[your_program] EOF  null
[tester::#XC5] [test-1] ✓ 3 line(s) match on stdout
[tester::#XC5] [test-1] ✓ Received exit code 0.
[tester::#XC5] [test-2] Running test case: 2
[tester::#XC5] [test-2] Writing contents to ./test.lox:
[tester::#XC5] [test-2] [test.lox] ++--**..,,;;
[tester::#XC5] [test-2] $ ./your_program.sh tokenize test.lox
[your_program] PLUS + null
[your_program] PLUS + null
[your_program] MINUS - null
[your_program] MINUS - null
[your_program] STAR * null
[your_program] STAR * null
[your_program] DOT . null
[your_program] DOT . null
[your_program] COMMA , null
[your_program] COMMA , null
[your_program] SEMICOLON ; null
[your_program] SEMICOLON ; null
[your_program] EOF  null
[tester::#XC5] [test-2] ✓ 13 line(s) match on stdout
[tester::#XC5] [test-2] ✓ Received exit code 0.
[tester::#XC5] [test-3] Running test case: 3
[tester::#XC5] [test-3] Writing contents to ./test.lox:
[tester::#XC5] [test-3] [test.lox] +;*.-,*
[tester::#XC5] [test-3] $ ./your_program.sh tokenize test.lox
[your_program] PLUS + null
[your_program] SEMICOLON ; null
[your_program] STAR * null
[your_program] DOT . null
[your_program] MINUS - null
[your_program] COMMA , null
[your_program] STAR * null
[your_program] EOF  null
[tester::#XC5] [test-3] ✓ 8 line(s) match on stdout
[tester::#XC5] [test-3] ✓ Received exit code 0.
[tester::#XC5] [test-4] Running test case: 4
[tester::#XC5] [test-4] Writing contents to ./test.lox:
[tester::#XC5] [test-4] [test.lox] ({.,;*+})
[tester::#XC5] [test-4] $ ./your_program.sh tokenize test.lox
[your_program] LEFT_PAREN ( null
[your_program] LEFT_BRACE { null
[your_program] DOT . null
[your_program] COMMA , null
[your_program] SEMICOLON ; null
[your_program] STAR * null
[your_program] PLUS + null
[your_program] RIGHT_BRACE } null
[your_program] RIGHT_PAREN ) null
[your_program] EOF  null
[tester::#XC5] [test-4] ✓ 10 line(s) match on stdout
[tester::#XC5] [test-4] ✓ Received exit code 0.
[tester::#XC5] Test passed.

[tester::#OE8] Running tests for Stage #OE8 (Scanning: Braces)
[tester::#OE8] [test-1] Running test case: 1
[tester::#OE8] [test-1] Writing contents to ./test.lox:
[tester::#OE8] [test-1] [test.lox] }
[tester::#OE8] [test-1] $ ./your_program.sh tokenize test.lox
[your_program] RIGHT_BRACE } null
[your_program] EOF  null
[tester::#OE8] [test-1] ✓ 2 line(s) match on stdout
[tester::#OE8] [test-1] ✓ Received exit code 0.
[tester::#OE8] [test-2] Running test case: 2
[tester::#OE8] [test-2] Writing contents to ./test.lox:
[tester::#OE8] [test-2] [test.lox] {{}}
[tester::#OE8] [test-2] $ ./your_program.sh tokenize test.lox
[your_program] LEFT_BRACE { null
[your_program] LEFT_BRACE { null
[your_program] RIGHT_BRACE } null
[your_program] RIGHT_BRACE } null
[your_program] EOF  null
[tester::#OE8] [test-2] ✓ 5 line(s) match on stdout
[tester::#OE8] [test-2] ✓ Received exit code 0.
[tester::#OE8] [test-3] Running test case: 3
[tester::#OE8] [test-3] Writing contents to ./test.lox:
[tester::#OE8] [test-3] [test.lox] {}{{}
[tester::#OE8] [test-3] $ ./your_program.sh tokenize test.lox
[your_program] LEFT_BRACE { null
[your_program] RIGHT_BRACE } null
[your_program] LEFT_BRACE { null
[your_program] LEFT_BRACE { null
[your_program] RIGHT_BRACE } null
[your_program] EOF  null
[tester::#OE8] [test-3] ✓ 6 line(s) match on stdout
[tester::#OE8] [test-3] ✓ Received exit code 0.
[tester::#OE8] [test-4] Running test case: 4
[tester::#OE8] [test-4] Writing contents to ./test.lox:
[tester::#OE8] [test-4] [test.lox] }){{)}(
[tester::#OE8] [test-4] $ ./your_program.sh tokenize test.lox
[your_program] RIGHT_BRACE } null
[your_program] RIGHT_PAREN ) null
[your_program] LEFT_BRACE { null
[your_program] LEFT_BRACE { null
[your_program] RIGHT_PAREN ) null
[your_program] RIGHT_BRACE } null
[your_program] LEFT_PAREN ( null
[your_program] EOF  null
[tester::#OE8] [test-4] ✓ 8 line(s) match on stdout
[tester::#OE8] [test-4] ✓ Received exit code 0.
[tester::#OE8] Test passed.

[tester::#OL4] Running tests for Stage #OL4 (Scanning: Parentheses)
[tester::#OL4] [test-1] Running test case: 1
[tester::#OL4] [test-1] Writing contents to ./test.lox:
[tester::#OL4] [test-1] [test.lox] (
[tester::#OL4] [test-1] $ ./your_program.sh tokenize test.lox
[your_program] LEFT_PAREN ( null
[your_program] EOF  null
[tester::#OL4] [test-1] ✓ 2 line(s) match on stdout
[tester::#OL4] [test-1] ✓ Received exit code 0.
[tester::#OL4] [test-2] Running test case: 2
[tester::#OL4] [test-2] Writing contents to ./test.lox:
[tester::#OL4] [test-2] [test.lox] ))
[tester::#OL4] [test-2] $ ./your_program.sh tokenize test.lox
[your_program] RIGHT_PAREN ) null
[your_program] RIGHT_PAREN ) null
[your_program] EOF  null
[tester::#OL4] [test-2] ✓ 3 line(s) match on stdout
[tester::#OL4] [test-2] ✓ Received exit code 0.
[tester::#OL4] [test-3] Running test case: 3
[tester::#OL4] [test-3] Writing contents to ./test.lox:
[tester::#OL4] [test-3] [test.lox] ())()
[tester::#OL4] [test-3] $ ./your_program.sh tokenize test.lox
[your_program] LEFT_PAREN ( null
[your_program] RIGHT_PAREN ) null
[your_program] RIGHT_PAREN ) null
[your_program] LEFT_PAREN ( null
[your_program] RIGHT_PAREN ) null
[your_program] EOF  null
[tester::#OL4] [test-3] ✓ 6 line(s) match on stdout
[tester::#OL4] [test-3] ✓ Received exit code 0.
[tester::#OL4] [test-4] Running test case: 4
[tester::#OL4] [test-4] Writing contents to ./test.lox:
[tester::#OL4] [test-4] [test.lox] (())))(
[tester::#OL4] [test-4] $ ./your_program.sh tokenize test.lox
[your_program] LEFT_PAREN ( null
[your_program] LEFT_PAREN ( null
[your_program] RIGHT_PAREN ) null
[your_program] RIGHT_PAREN ) null
[your_program] RIGHT_PAREN ) null
[your_program] RIGHT_PAREN ) null
[your_program] LEFT_PAREN ( null
[your_program] EOF  null
[tester::#OL4] [test-4] ✓ 8 line(s) match on stdout
[tester::#OL4] [test-4] ✓ Received exit code 0.
[tester::#OL4] Test passed.

*/
