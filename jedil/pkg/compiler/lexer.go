package compiler

import (
	"strconv"
	"unicode"
)

// Lexer tokenizes JEDIL source code
type Lexer struct {
	source  string
	start   int // Start of current token
	current int // Current position
	line    int // Current line number
}

// NewLexer creates a new lexer for the given source code
func NewLexer(source string) *Lexer {
	return &Lexer{
		source:  source,
		start:   0,
		current: 0,
		line:    1,
	}
}

// ScanToken scans and returns the next token
func (l *Lexer) ScanToken() Token {
	l.skipWhitespace()

	l.start = l.current

	if l.isAtEnd() {
		return l.makeToken(TOKEN_EOF)
	}

	c := l.advance()

	// Numbers
	if unicode.IsDigit(rune(c)) {
		return l.number()
	}

	// Identifiers and keywords
	if unicode.IsLetter(rune(c)) || c == '_' {
		return l.identifier()
	}

	// Single-character tokens
	switch c {
	case '(':
		return l.makeToken(TOKEN_LPAREN)
	case ')':
		return l.makeToken(TOKEN_RPAREN)
	case '{':
		return l.makeToken(TOKEN_LBRACE)
	case '}':
		return l.makeToken(TOKEN_RBRACE)
	case ',':
		return l.makeToken(TOKEN_COMMA)
	case '+':
		return l.makeToken(TOKEN_PLUS)
	case '-':
		return l.makeToken(TOKEN_MINUS)
	case '*':
		return l.makeToken(TOKEN_STAR)
	case '/':
		// Could be division or comment
		if l.match('/') {
			// Comment - skip to end of line
			for l.peek() != '\n' && !l.isAtEnd() {
				l.advance()
			}
			return l.ScanToken() // Recurse to get next real token
		}
		return l.makeToken(TOKEN_SLASH)
	case '=':
		return l.makeToken(TOKEN_EQUAL)
	}

	return l.errorToken("Unexpected character")
}

// Helper functions

func (l *Lexer) isAtEnd() bool {
	return l.current >= len(l.source)
}

func (l *Lexer) advance() byte {
	l.current++
	return l.source[l.current-1]
}

func (l *Lexer) peek() byte {
	if l.isAtEnd() {
		return 0
	}
	return l.source[l.current]
}

func (l *Lexer) peekNext() byte {
	if l.current+1 >= len(l.source) {
		return 0
	}
	return l.source[l.current+1]
}

func (l *Lexer) match(expected byte) bool {
	if l.isAtEnd() {
		return false
	}
	if l.source[l.current] != expected {
		return false
	}
	l.current++
	return true
}

func (l *Lexer) skipWhitespace() {
	for {
		if l.isAtEnd() {
			return
		}
		c := l.peek()
		switch c {
		case ' ', '\r', '\t':
			l.advance()
		case '\n':
			l.line++
			l.advance()
		default:
			return
		}
	}
}

func (l *Lexer) number() Token {
	// Scan integer part
	for unicode.IsDigit(rune(l.peek())) {
		l.advance()
	}

	// Decimal part
	if l.peek() == '.' && unicode.IsDigit(rune(l.peekNext())) {
		l.advance() // Consume '.'
		for unicode.IsDigit(rune(l.peek())) {
			l.advance()
		}
	}

	// Parse the number
	text := l.source[l.start:l.current]
	value, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return l.errorToken("Invalid number")
	}

	token := l.makeToken(TOKEN_NUMBER)
	token.Literal = value
	return token
}

func (l *Lexer) identifier() Token {
	for unicode.IsLetter(rune(l.peek())) || unicode.IsDigit(rune(l.peek())) || l.peek() == '_' {
		l.advance()
	}

	text := l.source[l.start:l.current]
	tokenType := l.identifierType(text)
	return l.makeToken(tokenType)
}

func (l *Lexer) identifierType(text string) TokenType {
	// Check if it's a keyword
	switch text {
	case "let":
		return TOKEN_LET
	case "fn":
		return TOKEN_FN
	case "return":
		return TOKEN_RETURN
	case "vec3":
		return TOKEN_VEC3
	default:
		return TOKEN_IDENTIFIER
	}
}

func (l *Lexer) makeToken(tokenType TokenType) Token {
	return Token{
		Type:   tokenType,
		Lexeme: l.source[l.start:l.current],
		Line:   l.line,
	}
}

func (l *Lexer) errorToken(message string) Token {
	return Token{
		Type:   TOKEN_ERROR,
		Lexeme: message,
		Line:   l.line,
	}
}
