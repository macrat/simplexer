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

func ExampleGetLastLine() {
	input := "this is a\ntest string\n"
	lexer := simplexer.NewLexer(strings.NewReader(input))

	for {
		token, err := lexer.Scan()
		if err != nil {
			panic(err.Error())
		}
		if token == nil {
			break
		}

		fmt.Println(lexer.GetLastLine())
		fmt.Printf("%s%s\n\n",
			strings.Repeat(" ", lexer.Position.Column),
			strings.Repeat("=", len(token.Literal)))
	}

	// Output:
	// this is a
	// ====
	//
	// this is a
	//      ==
	//
	// this is a
	//         =
	//
	// test string
	// ====
	//
	// test string
	//      ======
}

func ExampleNewTokenType() {
	number := simplexer.TokenID(0)
	others := simplexer.TokenID(1)

	lexer := simplexer.NewLexer(strings.NewReader("123this is test456"))

	lexer.TokenTypes = []simplexer.TokenType{
		simplexer.NewTokenType(number, `^[0-9]+`),
		simplexer.NewTokenType(others, `^[^0-9]+`),
	}

	for {
		token, _ := lexer.Scan()
		if token == nil {
			break
		}

		if token.Type.ID == number {
			fmt.Printf("%s is number\n", token.Literal)
		}

		if token.Type.ID == others {
			fmt.Printf("%s is not number\n", token.Literal)
		}
	}

	// Output:
	// 123 is number
	// this is test is not number
	// 456 is number
}
