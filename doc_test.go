package simplexer_test

import (
	"fmt"
	"strings"

	"github.com/macrat/simplexer"
)

func Example() {
	input := "hello_world = \"hello world\"\nnumber = 1"
	lexer := simplexer.NewLexer(strings.NewReader(input))

	fmt.Println(input)
	fmt.Println("==========")

	for {
		token, err := lexer.Scan()
		if err != nil {
			panic(err.Error())
		}
		if token == nil {
			fmt.Println("==========")
			return
		}

		fmt.Printf("line %2d, column %2d: %s: %s\n",
			lexer.Position.Line,
			lexer.Position.Column,
			token.Type,
			token.Literal)
	}

	// Output:
	// hello_world = "hello world"
	// number = 1
	// ==========
	// line  0, column  0: IDENT: hello_world
	// line  0, column 12: OTHER: =
	// line  0, column 14: STRING: "hello world"
	// line  1, column  0: IDENT: number
	// line  1, column  7: OTHER: =
	// line  1, column  9: NUMBER: 1
	// ==========
}
