package simplexer

import (
	"regexp"
	"strconv"
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

/*
Convert to readable string.

Be careful, user added token ID's will convert to UNKNOWN.
*/
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

// Compare TokenID as int.
func (id TokenID) Compare(another TokenID) int {
	return int(id - another)
}

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

// Get readable string of TokenID.
func (tt TokenType) String() string {
	return tt.ID.String()
}

// Compare TokenType of ID.
func (tt TokenType) Compare(another TokenType) int {
	return tt.ID.Compare(another.ID)
}

// A data of found Token.
type Token struct {
	Type       *TokenType
	Literal    string   // The string of matched.
	Submatches []string // Submatches of regular expression.
}
