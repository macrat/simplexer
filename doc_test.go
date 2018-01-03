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

		fmt.Printf("%s: %s\n", token.Type, token.Literal)
	}

	// Output:
	// hello_world = "hello world"
	// number = 1
	// ==========
	// IDENT: hello_world
	// OTHER: =
	// STRING: "hello world"
	// IDENT: number
	// OTHER: =
	// NUMBER: 1
	// ==========
}
