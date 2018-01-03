package simplexer

import (
	"io"
	"regexp"
	"strconv"
	"strings"
)

// TokenID is Identifier for TokenType.
type TokenID int

// Default token IDs.
const (
	OTHER TokenID = -(iota + 1)
	IDENT
	NUMBER
	STRING
)

func (id TokenID) String() string {
	switch id {
	case OTHER:
		return "OTHER"
	case IDENT:
		return "IDENT"
	case NUMBER:
		return "NUMBER"
	case STRING:
		return "STRING"
	default:
		return "UNKNOWN(" + strconv.Itoa(int(id)) + ")"
	}
}

func (id TokenID) Compare(another TokenID) int {
	return int(id - another)
}

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

// TokenType is a rule for making Token.
type TokenType struct {
	ID TokenID
	Re *regexp.Regexp // Regular expression for taking token. Must be starts with ^.
}

/*
Make new TokenType.

token: A TokenID of new TokenType.

re: A regular expression of token. Must be starts with ^.
*/
func NewTokenType(token TokenID, re string) TokenType {
	return TokenType{
		ID: token,
		Re: regexp.MustCompile(re),
	}
}

func (tt TokenType) String() string {
	return tt.ID.String()
}

// Compare TokenType of ID.
func (tt TokenType) Compare(another TokenType) int {
	return tt.ID.Compare(another.ID)
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
	Whitespace *regexp.Regexp
	TokenTypes []TokenType
}

// A data of found Token.
type Token struct {
	Type       *TokenType
	Literal    string   // The string of matched.
	Submatches []string // Submatches of regular expression.
}

// Make a new Lexer.
func NewLexer(reader io.Reader) *Lexer {
	l := new(Lexer)
	l.reader = reader

	l.Whitespace = DefaultWhitespace
	l.TokenTypes = DefaultTokenTypes

	return l
}

/*
Mathing buffer with a regular expression.

It won't consume buffer. Please use Lexer.Eat if want consuming.

Returns submatches.
*/
func (l *Lexer) Match(re *regexp.Regexp) []string {
	if len(l.buf) < 1024 {
		buf := make([]byte, 2048)
		l.reader.Read(buf)
		l.buf += strings.TrimRight(string(buf), "\x00")
	}

	m := re.FindStringSubmatch(l.buf)
	return m
}

/*
Matching buffer with a regular expression, and consume buffer if matched.

If don't want consume buffer, please use Lexer.Match.

Returns submatches.
*/
func (l *Lexer) Eat(re *regexp.Regexp) []string {
	match := l.Match(re)
	if len(match) > 0 {
		l.buf = l.buf[len(match[0]):]
	}
	return match
}

/*
Peek the first token in the buffer.

Returns nil as *Token if the buffer is empty.
*/
func (l *Lexer) Peek() (*Token, error) {
	l.Eat(l.Whitespace)

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
	} else {
		return nil, nil
	}
}

/*
Scan will get the first token in the buffer and remove it from the buffer.

This function using Lexer.Peek. Please read document of Peek.
*/
func (l *Lexer) Scan() (*Token, error) {
	t, e := l.Peek()

	if t != nil {
		l.Eat(t.Type.Re)
	}

	return t, e
}
