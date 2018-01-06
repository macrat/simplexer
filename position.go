package simplexer

import "fmt"

// Position in the file.
type Position struct {
	Line   int
	Column int
}

// Convert to string.
func (p Position) String() string {
	return fmt.Sprintf("[line:%d, column:%d]", p.Line, p.Column)
}

// Position.Before will check p is before than x.
func (p Position) Before(x Position) bool {
	return p.Line < x.Line || (p.Line == x.Line && p.Column < x.Column)
}

// Position.After will check p is after than x.
func (p Position) After(x Position) bool {
	return p.Line > x.Line || (p.Line == x.Line && p.Column > x.Column)
}
