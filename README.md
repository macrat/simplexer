simplexer
=========

[![GoDoc](https://godoc.org/github.com/macrat/simplexer?status.svg)](https://godoc.org/github.com/macrat/simplexer)

A simple lexical analyzser for Go.

## example
### simplest usage
``` go
package main

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
}
```

It is output as follow.

``` text
hello_world = "hello world"
number = 1
==========
line  0, column  0: IDENT: hello_world
line  0, column 12: OTHER: =
line  0, column 14: STRING: "hello world"
line  1, column  0: IDENT: number
line  1, column  7: OTHER: =
line  1, column  9: NUMBER: 1
==========
```

### define original token
Define new token type SUBSITUATION and NEWLINE. And append into Lexer.

``` go
package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/macrat/simplexer"
)

const (
	SUBSITUATION simplexer.TokenID = iota
	NEWLINE
)

func main() {
	input := "hello_world = \"hello world\"\nnumber = 1"
	lexer := simplexer.NewLexer(strings.NewReader(input))

	lexer.Whitespace = regexp.MustCompile(`^[\t ]`)

	lexer.TokenTypes = append([]simplexer.TokenType{
		simplexer.NewTokenType(SUBSITUATION, `^=`),
		simplexer.NewTokenType(NEWLINE, `^[\n\r]+`),
	}, lexer.TokenTypes...)

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

		fmt.Printf("%s: %#v\n", token.Type, token.Literal)
	}
}
```

Please be careful, all regular expression must starts with `^`.

It is output as follow.

``` text
hello_world = "hello world"
number = 1
==========
IDENT: "hello_world"
UNKNOWN(0): "="
STRING: "hello world"
UNKNOWN(1): "\n"
IDENT: "number"
UNKNOWN(0): "="
NUMBER: "1"
==========
```
