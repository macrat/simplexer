package simplexer

import "fmt"

// The error that returns when found an unknown token.
type SyntaxError string

func (se SyntaxError) Error() string {
	return fmt.Sprintf("SyntaxError: %#v", se)
}
