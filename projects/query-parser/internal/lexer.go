package internal

import (
	"strings"
	"unicode"
)

// TokenType represents the type of a token
type TokenType string

const (
	TokenSelect     TokenType = "SELECT"
	TokenFrom       TokenType = "FROM"
	TokenWhere      TokenType = "WHERE"
	TokenAnd        TokenType = "AND"
	TokenOr         TokenType = "OR"
	TokenNot        TokenType = "NOT"
	TokenInsert     TokenType = "INSERT"
	TokenInto       TokenType = "INTO"
	TokenValues     TokenType = "VALUES"
	TokenUpdate     TokenType = "UPDATE"
	TokenSet        TokenType = "SET"
	TokenDelete     TokenType = "DELETE"
	TokenCreate     TokenType = "CREATE"
	TokenTable      TokenType = "TABLE"
	TokenJoin       TokenType = "JOIN"
	TokenOn         TokenType = "ON"
	TokenAs         TokenType = "AS"
	TokenOrderBy    TokenType = "ORDER BY"
	TokenGroupBy    TokenType = "GROUP BY"
	TokenLimit      TokenType = "LIMIT"
	TokenOffset     TokenType = "OFFSET"
	TokenIdentifier TokenType = "IDENTIFIER"
	TokenNumber     TokenType = "NUMBER"
	TokenString     TokenType = "STRING"
	TokenOperator   TokenType = "OPERATOR"
	TokenComma      TokenType = "COMMA"
	TokenStar       TokenType = "STAR"
	TokenLParen     TokenType = "LPAREN"
	TokenRParen     TokenType = "RPAREN"
	TokenSemicolon  TokenType = "SEMICOLON"
	TokenEOF        TokenType = "EOF"
	TokenError      TokenType = "ERROR"
)

// Token represents a lexical token
type Token struct {
	Type     TokenType `json:"type"`
	Value    string    `json:"value"`
	Position Position  `json:"position"`
}

// Position represents the position of a token in the source
type Position struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

// Lexer tokenizes SQL input
type Lexer struct {
	input   string
	pos     int
	tokens  []Token
	current int
}

// Keywords maps SQL keywords to token types
var Keywords = map[string]TokenType{
	"SELECT":   TokenSelect,
	"FROM":     TokenFrom,
	"WHERE":    TokenWhere,
	"AND":      TokenAnd,
	"OR":       TokenOr,
	"NOT":      TokenNot,
	"INSERT":   TokenInsert,
	"INTO":     TokenInto,
	"VALUES":   TokenValues,
	"UPDATE":   TokenUpdate,
	"SET":      TokenSet,
	"DELETE":   TokenDelete,
	"CREATE":   TokenCreate,
	"TABLE":    TokenTable,
	"JOIN":     TokenJoin,
	"ON":       TokenOn,
	"AS":       TokenAs,
	"ORDER":    TokenOrderBy,
	"GROUP":    TokenGroupBy,
	"BY":       TokenIdentifier, // Will be combined with ORDER/GROUP
	"LIMIT":    TokenLimit,
	"OFFSET":   TokenOffset,
}

// NewLexer creates a new lexer for the given input
func NewLexer(input string) *Lexer {
	return &Lexer{
		input:   input,
		pos:     0,
		tokens:  []Token{},
		current: 0,
	}
}

// Tokenize tokenizes the entire input
func (l *Lexer) Tokenize() []Token {
	l.tokens = []Token{}
	l.pos = 0

	for l.pos < len(l.input) {
		l.skipWhitespace()
		if l.pos >= len(l.input) {
			break
		}

		token := l.nextToken()
		l.tokens = append(l.tokens, token)
	}

	// Add EOF token
	l.tokens = append(l.tokens, Token{
		Type:     TokenEOF,
		Value:    "",
		Position: Position{Start: l.pos, End: l.pos},
	})

	return l.tokens
}

// TokenizeStep returns tokens one at a time for step-by-step visualization
func (l *Lexer) TokenizeStep() *Token {
	l.skipWhitespace()
	if l.pos >= len(l.input) {
		eof := Token{
			Type:     TokenEOF,
			Value:    "",
			Position: Position{Start: l.pos, End: l.pos},
		}
		return &eof
	}

	token := l.nextToken()
	l.tokens = append(l.tokens, token)
	return &token
}

// GetTokens returns all collected tokens
func (l *Lexer) GetTokens() []Token {
	return l.tokens
}

func (l *Lexer) skipWhitespace() {
	for l.pos < len(l.input) && unicode.IsSpace(rune(l.input[l.pos])) {
		l.pos++
	}
}

func (l *Lexer) nextToken() Token {
	start := l.pos
	ch := l.input[l.pos]

	// Single character tokens
	switch ch {
	case ',':
		l.pos++
		return Token{Type: TokenComma, Value: ",", Position: Position{Start: start, End: l.pos}}
	case '*':
		l.pos++
		return Token{Type: TokenStar, Value: "*", Position: Position{Start: start, End: l.pos}}
	case '(':
		l.pos++
		return Token{Type: TokenLParen, Value: "(", Position: Position{Start: start, End: l.pos}}
	case ')':
		l.pos++
		return Token{Type: TokenRParen, Value: ")", Position: Position{Start: start, End: l.pos}}
	case ';':
		l.pos++
		return Token{Type: TokenSemicolon, Value: ";", Position: Position{Start: start, End: l.pos}}
	}

	// Operators
	if ch == '=' || ch == '<' || ch == '>' || ch == '!' {
		return l.readOperator(start)
	}

	// String literals
	if ch == '\'' || ch == '"' {
		return l.readString(start, ch)
	}

	// Numbers
	if unicode.IsDigit(rune(ch)) {
		return l.readNumber(start)
	}

	// Identifiers and keywords
	if unicode.IsLetter(rune(ch)) || ch == '_' {
		return l.readIdentifier(start)
	}

	// Unknown character
	l.pos++
	return Token{Type: TokenError, Value: string(ch), Position: Position{Start: start, End: l.pos}}
}

func (l *Lexer) readOperator(start int) Token {
	ch := l.input[l.pos]
	l.pos++

	// Check for two-character operators
	if l.pos < len(l.input) {
		next := l.input[l.pos]
		if (ch == '<' && next == '=') || (ch == '>' && next == '=') || (ch == '!' && next == '=') || (ch == '<' && next == '>') {
			l.pos++
			return Token{Type: TokenOperator, Value: string(ch) + string(next), Position: Position{Start: start, End: l.pos}}
		}
	}

	return Token{Type: TokenOperator, Value: string(ch), Position: Position{Start: start, End: l.pos}}
}

func (l *Lexer) readString(start int, quote byte) Token {
	l.pos++ // Skip opening quote
	for l.pos < len(l.input) && l.input[l.pos] != quote {
		if l.input[l.pos] == '\\' && l.pos+1 < len(l.input) {
			l.pos += 2 // Skip escaped character
		} else {
			l.pos++
		}
	}

	if l.pos >= len(l.input) {
		return Token{Type: TokenError, Value: "unterminated string", Position: Position{Start: start, End: l.pos}}
	}

	l.pos++ // Skip closing quote
	value := l.input[start+1 : l.pos-1]
	return Token{Type: TokenString, Value: value, Position: Position{Start: start, End: l.pos}}
}

func (l *Lexer) readNumber(start int) Token {
	for l.pos < len(l.input) && (unicode.IsDigit(rune(l.input[l.pos])) || l.input[l.pos] == '.') {
		l.pos++
	}
	return Token{Type: TokenNumber, Value: l.input[start:l.pos], Position: Position{Start: start, End: l.pos}}
}

func (l *Lexer) readIdentifier(start int) Token {
	for l.pos < len(l.input) && (unicode.IsLetter(rune(l.input[l.pos])) || unicode.IsDigit(rune(l.input[l.pos])) || l.input[l.pos] == '_') {
		l.pos++
	}

	value := l.input[start:l.pos]
	upper := strings.ToUpper(value)

	if tokenType, ok := Keywords[upper]; ok {
		return Token{Type: tokenType, Value: upper, Position: Position{Start: start, End: l.pos}}
	}

	return Token{Type: TokenIdentifier, Value: value, Position: Position{Start: start, End: l.pos}}
}

// HasMore returns true if there are more characters to tokenize
func (l *Lexer) HasMore() bool {
	l.skipWhitespace()
	return l.pos < len(l.input)
}

// GetPosition returns current position
func (l *Lexer) GetPosition() int {
	return l.pos
}
