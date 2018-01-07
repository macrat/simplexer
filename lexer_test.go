package simplexer_test

import (
	"strings"
	"testing"

	"github.com/macrat/simplexer"
)

type want struct {
	TypeID   simplexer.TokenID
	Literal  string
	Pos      simplexer.Position
	LastLine string
}

func execute(t *testing.T, input string, wants []want) {
	lexer := simplexer.NewLexer(strings.NewReader(input))

	for _, except := range wants {
		token, err := lexer.Scan()
		if err != nil {
			t.Fatalf(err.Error())
		}
		if token == nil {
			t.Fatalf("excepted token type=%s literal=%#v but got nil", except.TypeID, except.Literal)
		}

		if token.Type.ID != except.TypeID {
			t.Errorf("excepted type %s but got %s", except.TypeID, token.Type.ID)
		}
		if token.Literal != except.Literal {
			t.Errorf("excepted literal %#v but got %#v", except.Literal, token.Literal)
		}

		if lexer.Position != except.Pos {
			t.Errorf("excepted position %#v but got %#v", except.Pos, lexer.Position)
		}

		if lexer.GetLastLine() != except.LastLine {
			t.Errorf("excepted last line %#v but got %#v", except.LastLine, lexer.GetLastLine())
		}
	}

	token, err := lexer.Scan()
	if token != nil {
		t.Errorf("excepted end but got %#v", token)
	}
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestLexer(t *testing.T) {
	execute(t, "\t10; literal\nhoge = \"abc\"", []want{
		{
			TypeID:   simplexer.NUMBER,
			Literal:  "10",
			Pos:      simplexer.Position{Line: 0, Column: 1},
			LastLine: "\t10; literal",
		},
		{
			TypeID:   simplexer.OTHER,
			Literal:  ";",
			Pos:      simplexer.Position{Line: 0, Column: 3},
			LastLine: "\t10; literal",
		},
		{
			TypeID:   simplexer.IDENT,
			Literal:  "literal",
			Pos:      simplexer.Position{Line: 0, Column: 5},
			LastLine: "\t10; literal",
		},
		{
			TypeID:   simplexer.IDENT,
			Literal:  "hoge",
			Pos:      simplexer.Position{Line: 1, Column: 0},
			LastLine: "hoge = \"abc\"",
		},
		{
			TypeID:   simplexer.OTHER,
			Literal:  "=",
			Pos:      simplexer.Position{Line: 1, Column: 5},
			LastLine: "hoge = \"abc\"",
		},
		{
			TypeID:   simplexer.STRING,
			Literal:  "\"abc\"",
			Pos:      simplexer.Position{Line: 1, Column: 7},
			LastLine: "hoge = \"abc\"",
		},
	})
}

func TestLexer_oneLine(t *testing.T) {
	execute(t, "this is \"one line\"", []want{
		{
			TypeID:   simplexer.IDENT,
			Literal:  "this",
			Pos:      simplexer.Position{Line: 0, Column: 0},
			LastLine: "this is \"one line\"",
		},
		{
			TypeID:   simplexer.IDENT,
			Literal:  "is",
			Pos:      simplexer.Position{Line: 0, Column: 5},
			LastLine: "this is \"one line\"",
		},
		{
			TypeID:   simplexer.STRING,
			Literal:  "\"one line\"",
			Pos:      simplexer.Position{Line: 0, Column: 8},
			LastLine: "this is \"one line\"",
		},
	})
}

func TestLexer_reportingError(t *testing.T) {
	lexer := simplexer.NewLexer(strings.NewReader("1 2 error 3 4"))
	lexer.TokenTypes = []simplexer.TokenType{
		simplexer.NewTokenType(0, `^[0-9]+`),
	}

	if token, err := lexer.Scan(); err != nil {
		t.Fatalf("%s", err.Error())
	} else if token.Literal != "1" {
		t.Fatalf("except 1 but got %s", token.Literal)
	}

	if token, err := lexer.Scan(); err != nil {
		t.Fatalf("%s", err.Error())
	} else if token.Literal != "2" {
		t.Fatalf("except 2 but got %s", token.Literal)
	}

	token, e := lexer.Scan()
	if e == nil {
		t.Fatalf("except error but got nil")
	}
	if token != nil {
		t.Errorf("token when error except nil but got %s", token)
	}

	err, ok := e.(simplexer.UnknownTokenError)
	if !ok {
		t.Fatalf("except UnknownTokenError but got other error")
	}

	exceptPos := simplexer.Position{Line: 0, Column: 4}
	if err.Position != exceptPos {
		t.Errorf("position of error excepts %v but got %v", exceptPos, err.Position)
	}

	exceptLiteral := "error"
	if err.Literal != exceptLiteral {
		t.Errorf("literal of error excepts %s but got %s", exceptLiteral, err.Literal)
	}
}
