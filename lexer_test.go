package simplexer

import (
	"strings"
	"testing"
)

func TestLexer(t *testing.T) {
	lexer := NewLexer(strings.NewReader("\t10; literal\n\"abc\""))

	wants := []struct {
		TypeID  TokenID
		Literal string
	}{
		{
			TypeID:  NUMBER,
			Literal: "10",
		},
		{
			TypeID:  OTHER,
			Literal: ";",
		},
		{
			TypeID:  IDENT,
			Literal: "literal",
		},
		{
			TypeID:  STRING,
			Literal: "\"abc\"",
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
	}

	token, err := lexer.Scan()
	if token != nil {
		t.Errorf("excepted end but got %#v", token)
	}
	if err != nil {
		t.Errorf(err.Error())
	}
}
