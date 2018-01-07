package simplexer

import (
	"io"
	"regexp"
	"strings"
)

// Defined default values for properties of Lexer as a package value.
var (
	DefaultWhitespace = regexp.MustCompile(`^(\s|\r|\n)+`)

	DefaultTokenTypes = []TokenType{
		NewTokenType(IDENT, `^[a-zA-Z_][a-zA-Z0-9_]*`),
		NewTokenType(NUMBER, `^[0-9]+(\.[0-9]+)?`),
		NewTokenType(STRING, `^\"([^"]*)\"`),
		NewTokenType(OTHER, `^.`),
	}
)

/*
The lexical analyzer.

Whitespace is a regular expression for skipping characters like whitespaces.
The default value is simplexer.DefaultWhitespace.

TokenTypes is an array of TokenType.
Lexer will sequential check TokenTypes, and return first matched token.
Default is simplexer.DefaultTokenTypes.

Please be careful, Lexer will never use it even if append TokenType after OTHER.
Because OTHER will accept any single character.
*/
type Lexer struct {
	reader     io.Reader
	buf        string
	loadedLine string
	nextPos    Position
	Whitespace *regexp.Regexp
	TokenTypes []TokenType
}

// Make a new Lexer.
func NewLexer(reader io.Reader) *Lexer {
	l := new(Lexer)
	l.reader = reader

	l.Whitespace = DefaultWhitespace
	l.TokenTypes = DefaultTokenTypes

	return l
}

func (l *Lexer) readBufIfNeed() {
	if len(l.buf) < 1024 {
		buf := make([]byte, 2048)
		l.reader.Read(buf)
		l.buf += strings.TrimRight(string(buf), "\x00")
	}
}

/*
Mathing buffer with a regular expression.

Returns submatches.
*/
func (l *Lexer) Match(re *regexp.Regexp) []string {
	l.readBufIfNeed()

	if m := l.Whitespace.FindString(l.buf); m != "" {
		l.consumeBuffer(m)
	}

	l.readBufIfNeed()

	return re.FindStringSubmatch(l.buf)
}

func (l *Lexer) consumeBuffer(s string) {
	l.buf = l.buf[len(s):]

	l.nextPos = shiftPos(l.nextPos, s)

	if idx := strings.LastIndex(s, "\n"); idx >= 0 {
		l.loadedLine = s[idx+1:]
	} else {
		l.loadedLine += s
	}
}

func (l *Lexer) makeError() error {
	for i, _ := range l.buf {
		if l.Whitespace.MatchString(l.buf[i:]) {
			return UnknownTokenError{
				Literal:  l.buf[:i],
				Position: l.nextPos,
			}
		}

		for _, tokenType := range l.TokenTypes {
			if tokenType.Re.MatchString(l.buf[i:]) {
				return UnknownTokenError{
					Literal:  l.buf[:i],
					Position: l.nextPos,
				}
			}
		}
	}

	return UnknownTokenError{
		Literal:  l.buf,
		Position: l.nextPos,
	}
}

/*
Peek the first token in the buffer.

Returns nil as *Token if the buffer is empty.
*/
func (l *Lexer) Peek() (*Token, error) {
	for _, tokenType := range l.TokenTypes {
		if m := l.Match(tokenType.Re); len(m) > 0 {
			return &Token{
				Type:       &tokenType,
				Literal:    m[0],
				Submatches: m[1:],
				Position:   l.nextPos,
			}, nil
		}
	}

	if len(l.buf) > 0 {
		return nil, l.makeError()
	}

	return nil, nil
}

/*
Scan will get the first token in the buffer and remove it from the buffer.

This function using Lexer.Peek. Please read document of Peek.
*/
func (l *Lexer) Scan() (*Token, error) {
	t, e := l.Peek()

	if t != nil {
		l.consumeBuffer(t.Literal)
	}

	return t, e
}

/*
GetCurrentLine returns line of last scanned token.
*/
func (l *Lexer) GetLastLine() string {
	l.readBufIfNeed()

	if idx := strings.Index(l.buf, "\n"); idx >= 0 {
		return l.loadedLine + l.buf[:strings.Index(l.buf, "\n")]
	} else {
		return l.loadedLine + l.buf
	}
}
