import os

GoHeader = """
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
"""

GoTail = """
}


func doTest(lines []byte) (string, string, int) {
	output := ""
	retval := 0

	fmt.Printf("lines is %v\\n", lines)

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
            t.Errorf("%s: output does not match:\\n\\texpected '%#v'\\n\\tgot    : '%#v'\\n", test.name, test.output, output)
		}
		if test.errors != errs {
            t.Errorf("%s: errors does not match:\\n\\texpected '%#v'\\n\\tgot    : '%#v'\\n", test.name, test.errors, errs)
		}
		if test.retval != retval {
            t.Errorf("%s: retval does not match:\\n\\texpected %d\\n\\tgot       : %d\\n", test.name, test.retval, retval)
		}
	}
}

"""

TestCases = [
    {
        "Group": "TZ7",
        "Test": "test-1",
        "Input": """
() 	@""",
        "Output": """\
LEFT_PAREN ( null
RIGHT_PAREN ) null
EOF  null
""",
        "Error": """\
[line 2] Error: Unexpected character: @
""",
        "ReturnCode": 65,
    },
    {
        "Group": "TZ7",
        "Test": "test-2",
        "Input": """\
 	
@
""",
        "Output": """\
EOF  null
""",
        "Error": """\
[line 2] Error: Unexpected character: @
""",
        "ReturnCode": 65,
    },
    {
        "Group": "TZ7",
        "Test": "test-3",
        "Input": """\
()  #	{}
@
$
+++
// Let's Go!
+++
#
""",
        "Output": """\
LEFT_PAREN ( null
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
""",
        "Error": """\
[line 1] Error: Unexpected character: #
[line 2] Error: Unexpected character: @
[line 3] Error: Unexpected character: $
[line 7] Error: Unexpected character: #
""",
        "ReturnCode": 65,
    },
    {
        "Group": "TZ7",
        "Test": "test-4",
        "Input": """\
({*	%})
""",
        "Output": """\
LEFT_PAREN ( null
LEFT_BRACE { null
STAR * null
RIGHT_BRACE } null
RIGHT_PAREN ) null
EOF  null
""",
        "Error": """\
[line 1] Error: Unexpected character: %
""",
        "ReturnCode": 65,
    },
    {
        "Group": "R2",
        "Test": "test-1",
        "Input": """\
 
""",
        "Output": """\
EOF  null
""",
        "Error": """\
""",
        "ReturnCode": 0,
    },
    {
        "Group": "ER2",
        "Test": "test-2",
        "Input": """\
 	
 
""",
        "Output": """\
EOF  null
""",
        "Error": "",
        "ReturnCode": 0,
    },
    {
        "Group": "ER2",
        "Test": "test-3",
        "Input": """\
{
	}
((	-.*,))
""",
        "Output": """\
LEFT_BRACE { null
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
""",
        "Error": "",
        "ReturnCode": 0,
    },
    {
        "Group": "ER2",
        "Test": "test-4",
        "Input": """\
{ 		

}
((- >	))
""",
        "Output": """\
LEFT_BRACE { null
RIGHT_BRACE } null
LEFT_PAREN ( null
LEFT_PAREN ( null
MINUS - null
GREATER > null
RIGHT_PAREN ) null
RIGHT_PAREN ) null
EOF  null
""",
        "Error": "",
        "ReturnCode": 0,
    },
    {
        "Group": "ML2",
        "Test": "test-1",
        "Input": """\
//Comment
""",
        "Output": """\
EOF  null
""",
        "Error": "",
        "ReturnCode": 0,
    },
    {
        "Group": "ML2",
        "Test": "test-2",
        "Input": """\
(///Unicode:£§᯽☺♣)
""",
        "Output": """\
LEFT_PAREN ( null
EOF  null
""",
        "Error": "",
        "ReturnCode": 0,
    },
    {
        "Group": "ML2",
        "Test": "test-3",
        "Input": """\
/
""",
        "Output": """\
SLASH / null
EOF  null
""",
        "Error": "",
        "ReturnCode": 0,
    },
    {
        "Group": "ML2",
        "Test": "test-4",
        "Input": """\
({(===!)})//Comment
""",
        "Output": """\
LEFT_PAREN ( null
LEFT_BRACE { null
LEFT_PAREN ( null
EQUAL_EQUAL == null
EQUAL = null
BANG ! null
RIGHT_PAREN ) null
RIGHT_BRACE } null
RIGHT_PAREN ) null
EOF  null
""",
        "Error": "",
        "ReturnCode": 0,
    },
    {
        "Group": "ET2",
        "Test": "test-1",
        "Input": """\
>=
""",
        "Output": """\
GREATER_EQUAL >= null
EOF  null
""",
        "Error": "",
        "ReturnCode": 0,
    },
    {
        "Group": "ET2",
        "Test": "test-2",
        "Input": """\
<<<=>>>=
""",
        "Output": """\
LESS < null
LESS < null
LESS_EQUAL <= null
GREATER > null
GREATER > null
GREATER_EQUAL >= null
EOF  null
""",
        "Error": "",
        "ReturnCode": 0,
    },
    {
        "Group": "ET2",
        "Test": "test-3",
        "Input": """\
>=<>=>>
""",
        "Output": """\
GREATER_EQUAL >= null
LESS < null
GREATER_EQUAL >= null
GREATER > null
GREATER > null
EOF  null
""",
        "Error": "",
        "ReturnCode": 0,
    },
    {
        "Group": "ET2",
        "Test": "test-4",
        "Input": """\
(){>!>=}
""",
        "Output": """\
LEFT_PAREN ( null
RIGHT_PAREN ) null
LEFT_BRACE { null
GREATER > null
BANG ! null
GREATER_EQUAL >= null
RIGHT_BRACE } null
EOF  null
""",
        "Error": "",
        "ReturnCode": 0,
    },
    {
        "Group": "BU3",
        "Test": "test-1",
        "Input": """\
!=
""",
        "Output": """\
BANG_EQUAL != null
EOF  null
""",
        "Error": "",
        "ReturnCode": 0,
    },
    {
        "Group": "BU3",
        "Test": "test-2",
        "Input": """\
!!===
""",
        "Output": """\
BANG ! null
BANG_EQUAL != null
EQUAL_EQUAL == null
EOF  null
""",
        "Error": "",
        "ReturnCode": 0,
    },
    {
        "Group": "BU3",
        "Test": "test-3",
        "Input": """\
!{!}(!===)=
""",
        "Output": """\
BANG ! null
LEFT_BRACE { null
BANG ! null
RIGHT_BRACE } null
LEFT_PAREN ( null
BANG_EQUAL != null
EQUAL_EQUAL == null
RIGHT_PAREN ) null
EQUAL = null
EOF  null
""",
        "Error": "",
        "ReturnCode": 0,
    },
    {
        "Group": "BU3",
        "Test": "test-4",
        "Input": """\
{(%=#!==)}
""",
        "Output": """\
LEFT_BRACE { null
LEFT_PAREN ( null
EQUAL = null
BANG_EQUAL != null
EQUAL = null
RIGHT_PAREN ) null
RIGHT_BRACE } null
EOF  null
""",
        "Error": """\
[line 1] Error: Unexpected character: %
[line 1] Error: Unexpected character: #
""",
        "ReturnCode": 65,
    },
    {
        "Group": "MP7",
        "Test": "test-1",
        "Input": """\
=
""",
        "Output": """\
EQUAL = null
EOF  null
""",
        "Error": "",
        "ReturnCode": 0,
    },
    {
        "Group": "MP7",
        "Test": "test-2",
        "Input": """\
==
""",
        "Output": """\
EQUAL_EQUAL == null
EOF  null
""",
        "Error": "",
        "ReturnCode": 0,
    },
    {
        "Group": "MP7",
        "Test": "test-3",
        "Input": """\
({=}){==}
""",
        "Output": """\
LEFT_PAREN ( null
LEFT_BRACE { null
EQUAL = null
RIGHT_BRACE } null
RIGHT_PAREN ) null
LEFT_BRACE { null
EQUAL_EQUAL == null
RIGHT_BRACE } null
EOF  null
""",
        "Error": "",
        "ReturnCode": 0,
    },
    {
        "Group": "MP7",
        "Test": "test-4",
        "Input": """\
(($#===%))
""",
        "Output": """\
LEFT_PAREN ( null
LEFT_PAREN ( null
EQUAL_EQUAL == null
EQUAL = null
RIGHT_PAREN ) null
RIGHT_PAREN ) null
EOF  null
""",
        "Error": """\
[line 1] Error: Unexpected character: $
[line 1] Error: Unexpected character: #
[line 1] Error: Unexpected character: %
""",
        "ReturnCode": 65,
    },
    {
        "Group": "EA6",
        "Test": "test-1",
        "Input": """\
@
""",
        "Output": """\
EOF  null
""",
        "Error": """\
[line 1] Error: Unexpected character: @
""",
        "ReturnCode": 65,
    },
    {
        "Group": "EA6",
        "Test": "test-2",
        "Input": """\
,.$(#
""",
        "Output": """\
COMMA , null
DOT . null
LEFT_PAREN ( null
EOF  null
""",
        "Error": """\
[line 1] Error: Unexpected character: $
[line 1] Error: Unexpected character: #
""",
        "ReturnCode": 65,
    },
    {
        "Group": "EA6",
        "Test": "test-3",
        "Input": """\
%$#@%
""",
        "Output": """\
EOF  null
""",
        "Error": """\
[line 1] Error: Unexpected character: %
[line 1] Error: Unexpected character: $
[line 1] Error: Unexpected character: #
[line 1] Error: Unexpected character: @
[line 1] Error: Unexpected character: %
""",
        "ReturnCode": 65,
    },
    {
        "Group": "EA6",
        "Test": "test-4",
        "Input": """\
{(*-#;.$%)}
""",
        "Output": """\
LEFT_BRACE { null
LEFT_PAREN ( null
STAR * null
MINUS - null
SEMICOLON ; null
DOT . null
RIGHT_PAREN ) null
RIGHT_BRACE } null
EOF  null
""",
        "Error": """\
[line 1] Error: Unexpected character: #
[line 1] Error: Unexpected character: $
[line 1] Error: Unexpected character: %
""",
        "ReturnCode": 65,
    },
    {
        "Group": "XC5",
        "Test": "test-1",
        "Input": """\
+-
""",
        "Output": """\
PLUS + null
MINUS - null
EOF  null
""",
        "Error": "",
        "ReturnCode": 0,
    },
    {
        "Group": "XC5",
        "Test": "test-2",
        "Input": """\
++--**..,,;;
""",
        "Output": """\
PLUS + null
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
""",
        "Error": "",
        "ReturnCode": 0,
    },
    {
        "Group": "XC5",
        "Test": "test-3",
        "Input": """\
+;*.-,*
""",
        "Output": """\
PLUS + null
SEMICOLON ; null
STAR * null
DOT . null
MINUS - null
COMMA , null
STAR * null
EOF  null
""",
        "Error": "",
        "ReturnCode": 0,
    },
    {
        "Group": "XC5",
        "Test": "test-4",
        "Input": """\
({.,;*+})
""",
        "Output": """\
LEFT_PAREN ( null
LEFT_BRACE { null
DOT . null
COMMA , null
SEMICOLON ; null
STAR * null
PLUS + null
RIGHT_BRACE } null
RIGHT_PAREN ) null
EOF  null
""",
        "Error": "",
        "ReturnCode": 0,
    },
    {
        "Group": "OE8",
        "Test": "test-1",
        "Input": """\
}
""",
        "Output": """\
RIGHT_BRACE } null
EOF  null
""",
        "Error": "",
        "ReturnCode": 0,
    },
    {
        "Group": "OE8",
        "Test": "test-2",
        "Input": """\
{{}}
""",
        "Output": """\
LEFT_BRACE { null
LEFT_BRACE { null
RIGHT_BRACE } null
RIGHT_BRACE } null
EOF  null
""",
        "Error": "",
        "ReturnCode": 0,
    },
    {
        "Group": "OE8",
        "Test": "test-3",
        "Input": """\
{}{{}
""",
        "Output": """\
LEFT_BRACE { null
RIGHT_BRACE } null
LEFT_BRACE { null
LEFT_BRACE { null
RIGHT_BRACE } null
EOF  null
""",
        "Error": "",
        "ReturnCode": 0,
    },
    {
        "Group": "OE8",
        "Test": "test-4",
        "Input": """\
}){{)}(
""",
        "Output": """\
RIGHT_BRACE } null
RIGHT_PAREN ) null
LEFT_BRACE { null
LEFT_BRACE { null
RIGHT_PAREN ) null
RIGHT_BRACE } null
LEFT_PAREN ( null
EOF  null
""",
        "Error": "",
        "ReturnCode": 0,
    },
    {
        "Group": "OL4",
        "Test": "test-1",
        "Input": """\
(
""",
        "Output": """\
LEFT_PAREN ( null
EOF  null
""",
        "Error": "",
        "ReturnCode": 0,
    },
    {
        "Group": "OL4",
        "Test": "test-2",
        "Input": """\
))
""",
        "Output": """\
RIGHT_PAREN ) null
RIGHT_PAREN ) null
EOF  null
""",
        "Error": "",
        "ReturnCode": 0,
    },
    {
        "Group": "OL4",
        "Test": "test-3",
        "Input": """\
())()
""",
        "Output": """\
LEFT_PAREN ( null
RIGHT_PAREN ) null
RIGHT_PAREN ) null
LEFT_PAREN ( null
RIGHT_PAREN ) null
EOF  null
""",
        "Error": "",
        "ReturnCode": 0,
    },
    {
        "Group": "OL4",
        "Test": "test-4",
        "Input": """\
(())))(
""",
        "Output": """\
LEFT_PAREN ( null
LEFT_PAREN ( null
RIGHT_PAREN ) null
RIGHT_PAREN ) null
RIGHT_PAREN ) null
RIGHT_PAREN ) null
LEFT_PAREN ( null
EOF  null
""",
        "Error": "",
        "ReturnCode": 0,
    },
    {
        "Group": "RY8",
        "Test": "test-1",
        "Input": "",
        "Output": """\
EOF  null
""",
        "Error": "",
        "ReturnCode": 0,
    },
]

print(GoHeader)

for test in TestCases:
    body = f"""
    {{
		name:   "[{test["Group"]}] [{test["Test"]}]",
		lines:  `{test["Input"]}`,
		errors: `{test["Error"]}`,
		output: `{test["Output"]}`,
		retval: {test["ReturnCode"]},
    }},
    """
    print(body)
    # testDir = f"./Tests/{test['Group']}/{test['Test']}"
    # os.makedirs(testDir, exist_ok=True)
    # print(f"Creating {testDir}...")
    # with open(f"{testDir}/test.lox", "w") as err:
    #     err.write(test["Input"])
    # with open(f"{testDir}/output.text", "w") as err:
    #     err.write(test["Output"])
    # with open(f"{testDir}/errors.text", "w") as err:
    #     err.write(test["Error"])


print(GoTail)
