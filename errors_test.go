package simplexer

import "testing"

func TestSyntaxError(t *testing.T) {
	err := SyntaxError("test")
	except := "SyntaxError: \"test\""

	if err.Error() != except {
		t.Errorf("excepted %#v but got %s", except, err.Error())
	}
}
