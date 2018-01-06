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

// Position in the file.
type Position struct {
	Line   int
	Column int
}

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
	Whitespace *regexp.Regexp
	TokenTypes []TokenType
	Position   Position // Current position in the input.
	NextPos    Position // Position of next token.
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

func shiftPos(p Position, s string) Position {
	lines := strings.Split(s, "\n")
	lineShift := len(lines) - 1

	if lineShift == 0 {
		p.Column += len(lines[0])
	} else {
		p.Column = len(lines[len(lines)-1])
	}
	p.Line += lineShift

	return p
}

/*
Mathing buffer with a regular expression.

It won't consume buffer. Please use Lexer.Eat if want consuming.

Returns submatches.
*/
func (l *Lexer) Match(re *regexp.Regexp) []string {
	l.readBufIfNeed()

	if m := l.Whitespace.FindString(l.buf); m != "" {
		l.consumeBuffer(m, false)
	}

	l.readBufIfNeed()

	return re.FindStringSubmatch(l.buf)
}

func (l *Lexer) consumeBuffer(s string, isToken bool) {
	if len(s) > 0 {
		l.buf = l.buf[len(s):]

		if isToken {
			l.Position = l.NextPos
		}
		l.NextPos = shiftPos(l.NextPos, s)

		if idx := strings.LastIndex(s, "\n"); idx >= 0 {
			l.loadedLine = s[idx+1:]
		} else {
			l.loadedLine += s
		}
	}
}

/*
Matching buffer with a regular expression, and consume buffer if matched.

If don't want consume buffer, please use Lexer.Match.

Returns submatches.
*/
func (l *Lexer) Eat(re *regexp.Regexp) []string {
	match := l.Match(re)

	if len(match) > 0 {
		l.consumeBuffer(match[0], true)
	}

	return match
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
			}, nil
		}
	}

	if len(l.buf) > 0 {
		return nil, SyntaxError(l.buf)
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
		l.consumeBuffer(t.Literal, true)
	}

	return t, e
}

/*
GetCurrentLine returns line of last scanned token.
*/
func (l *Lexer) GetLastLine() string {
	l.readBufIfNeed()
	return l.loadedLine + l.buf[:strings.Index(l.buf, "\n")]
}
