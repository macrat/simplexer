package simplexer_test

import (
	"fmt"
	"strings"

	"github.com/macrat/simplexer"
)

func ExampleNewTokenType() {
	const (
		NUMBER simplexer.TokenID = iota
		OTHERS
	)

	lexer := simplexer.NewLexer(strings.NewReader("123this is test456"))

	lexer.TokenTypes = []simplexer.TokenType{
		simplexer.NewTokenType(NUMBER, `^[0-9]+`),
		simplexer.NewTokenType(OTHERS, `^[^0-9]+`),
	}

	for {
		token, _ := lexer.Scan()
		if token == nil {
			break
		}

		if token.Type.ID == NUMBER {
			fmt.Printf("%s is number\n", token.Literal)
		}

		if token.Type.ID == OTHERS {
			fmt.Printf("%s is not number\n", token.Literal)
		}
	}

	// Output:
	// 123 is number
	// this is test is not number
	// 456 is number
}
