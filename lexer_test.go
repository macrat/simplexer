package simplexer

import (
	"strings"
	"testing"
)

func TestLexer(t *testing.T) {
	lexer := NewLexer(strings.NewReader("\t10; literal\nhoge = \"abc\""))

	wants := []struct {
		TypeID  TokenID
		Literal string
		Pos     Position
	}{
		{
			TypeID:  NUMBER,
			Literal: "10",
			Pos:     Position{Line: 0, Column: 1},
		},
		{
			TypeID:  OTHER,
			Literal: ";",
			Pos:     Position{Line: 0, Column: 3},
		},
		{
			TypeID:  IDENT,
			Literal: "literal",
			Pos:     Position{Line: 0, Column: 5},
		},
		{
			TypeID:  IDENT,
			Literal: "hoge",
			Pos:     Position{Line: 1, Column: 0},
		},
		{
			TypeID:  OTHER,
			Literal: "=",
			Pos:     Position{Line: 1, Column: 5},
		},
		{
			TypeID:  STRING,
			Literal: "\"abc\"",
			Pos:     Position{Line: 1, Column: 7},
		},
	}

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
	}

	token, err := lexer.Scan()
	if token != nil {
		t.Errorf("excepted end but got %#v", token)
	}
	if err != nil {
		t.Errorf(err.Error())
	}
}
