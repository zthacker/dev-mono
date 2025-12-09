package compiler

// TokenType represents the type of a lexical token
type TokenType int

const (
	// Special tokens
	TOKEN_EOF TokenType = iota
	TOKEN_ERROR

	// Literals
	TOKEN_NUMBER     // 123.45
	TOKEN_IDENTIFIER // variable_name

	// Keywords
	TOKEN_LET    // let
	TOKEN_FN     // fn
	TOKEN_RETURN // return
	TOKEN_VEC3   // vec3

	// Operators
	TOKEN_PLUS  // +
	TOKEN_MINUS // -
	TOKEN_STAR  // *
	TOKEN_SLASH // /
	TOKEN_EQUAL // =

	// Delimiters
	TOKEN_LPAREN // (
	TOKEN_RPAREN // )
	TOKEN_LBRACE // {
	TOKEN_RBRACE // }
	TOKEN_COMMA  // ,
)

// Token represents a lexical token
type Token struct {
	Type    TokenType
	Lexeme  string  // Original text
	Literal float64 // For numbers
	Line    int     // Line number (for error reporting)
}

func (t TokenType) String() string {
	switch t {
	case TOKEN_EOF:
		return "EOF"
	case TOKEN_ERROR:
		return "ERROR"
	case TOKEN_NUMBER:
		return "NUMBER"
	case TOKEN_IDENTIFIER:
		return "IDENTIFIER"
	case TOKEN_LET:
		return "LET"
	case TOKEN_FN:
		return "FN"
	case TOKEN_RETURN:
		return "RETURN"
	case TOKEN_VEC3:
		return "VEC3"
	case TOKEN_PLUS:
		return "PLUS"
	case TOKEN_MINUS:
		return "MINUS"
	case TOKEN_STAR:
		return "STAR"
	case TOKEN_SLASH:
		return "SLASH"
	case TOKEN_EQUAL:
		return "EQUAL"
	case TOKEN_LPAREN:
		return "LPAREN"
	case TOKEN_RPAREN:
		return "RPAREN"
	case TOKEN_LBRACE:
		return "LBRACE"
	case TOKEN_RBRACE:
		return "RBRACE"
	case TOKEN_COMMA:
		return "COMMA"
	default:
		return "UNKNOWN"
	}
}
