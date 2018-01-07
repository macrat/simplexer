package simplexer_test

import (
	"testing"

	"github.com/macrat/simplexer"
)

func TestSyntaxError(t *testing.T) {
	err := simplexer.SyntaxError("test")
	except := "SyntaxError: \"test\""

	if err.Error() != except {
		t.Errorf("excepted %#v but got %s", except, err.Error())
	}
}
