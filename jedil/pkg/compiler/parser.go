package compiler

import (
	"fmt"
)

// Parser converts tokens into an AST
type Parser struct {
	lexer   *Lexer
	prev    Token
	current Token
	peek    Token
}

// NewParser creates a new parser
func NewParser(source string) *Parser {
	lexer := NewLexer(source)
	p := &Parser{
		lexer: lexer,
	}
	// Prime the pump with two tokens
	p.advance()
	p.advance()
	return p
}

// Parse parses the entire program
func (p *Parser) Parse() (*Program, error) {
	var statements []Stmt

	for !p.check(TOKEN_EOF) {
		stmt, err := p.statement()
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)
	}

	return &Program{Statements: statements}, nil
}

// Statement parsing

func (p *Parser) statement() (Stmt, error) {
	if p.match(TOKEN_LET) {
		return p.varDecl()
	}
	if p.match(TOKEN_FN) {
		return p.fnDecl()
	}
	if p.match(TOKEN_RETURN) {
		return p.returnStmt()
	}

	return nil, p.error("Expected statement")
}

func (p *Parser) varDecl() (Stmt, error) {
	if !p.check(TOKEN_IDENTIFIER) {
		return nil, p.error("Expected variable name")
	}
	name := p.current.Lexeme
	p.advance()

	if !p.match(TOKEN_EQUAL) {
		return nil, p.error("Expected '=' after variable name")
	}

	value, err := p.expression()
	if err != nil {
		return nil, err
	}

	return &VarDecl{Name: name, Value: value}, nil
}

func (p *Parser) fnDecl() (Stmt, error) {
	if !p.check(TOKEN_IDENTIFIER) {
		return nil, p.error("Expected function name")
	}
	name := p.current.Lexeme
	p.advance()

	if !p.match(TOKEN_LPAREN) {
		return nil, p.error("Expected '(' after function name")
	}

	// Parse parameters
	var params []string
	if !p.check(TOKEN_RPAREN) {
		for {
			if !p.check(TOKEN_IDENTIFIER) {
				return nil, p.error("Expected parameter name")
			}
			params = append(params, p.current.Lexeme)
			p.advance()

			if !p.match(TOKEN_COMMA) {
				break
			}
		}
	}

	if !p.match(TOKEN_RPAREN) {
		return nil, p.error("Expected ')' after parameters")
	}

	if !p.match(TOKEN_LBRACE) {
		return nil, p.error("Expected '{' before function body")
	}

	// Parse body
	var body []Stmt
	for !p.check(TOKEN_RBRACE) && !p.check(TOKEN_EOF) {
		stmt, err := p.statement()
		if err != nil {
			return nil, err
		}
		body = append(body, stmt)
	}

	if !p.match(TOKEN_RBRACE) {
		return nil, p.error("Expected '}' after function body")
	}

	return &FnDecl{Name: name, Params: params, Body: body}, nil
}

func (p *Parser) returnStmt() (Stmt, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}
	return &ReturnStmt{Value: value}, nil
}

// Expression parsing (precedence climbing)

func (p *Parser) expression() (Expr, error) {
	return p.addition()
}

func (p *Parser) addition() (Expr, error) {
	expr, err := p.multiplication()
	if err != nil {
		return nil, err
	}

	for p.match(TOKEN_PLUS, TOKEN_MINUS) {
		op := p.previous().Type
		right, err := p.multiplication()
		if err != nil {
			return nil, err
		}
		expr = &BinaryOp{Left: expr, Op: op, Right: right}
	}

	return expr, nil
}

func (p *Parser) multiplication() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(TOKEN_STAR, TOKEN_SLASH) {
		op := p.previous().Type
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = &BinaryOp{Left: expr, Op: op, Right: right}
	}

	return expr, nil
}

func (p *Parser) unary() (Expr, error) {
	if p.match(TOKEN_MINUS) {
		expr, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &UnaryOp{Op: TOKEN_MINUS, Expr: expr}, nil
	}

	return p.call()
}

func (p *Parser) call() (Expr, error) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}

	// Check for function call
	if p.match(TOKEN_LPAREN) {
		// Must be a function call - expr should be an identifier
		if ident, ok := expr.(*Identifier); ok {
			args, err := p.arguments()
			if err != nil {
				return nil, err
			}
			if !p.match(TOKEN_RPAREN) {
				return nil, p.error("Expected ')' after arguments")
			}
			return &CallExpr{Callee: ident.Name, Args: args}, nil
		}
		return nil, p.error("Only identifiers can be called")
	}

	return expr, nil
}

func (p *Parser) arguments() ([]Expr, error) {
	var args []Expr

	if p.check(TOKEN_RPAREN) {
		return args, nil
	}

	for {
		arg, err := p.expression()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)

		if !p.match(TOKEN_COMMA) {
			break
		}
	}

	return args, nil
}

func (p *Parser) primary() (Expr, error) {
	if p.match(TOKEN_NUMBER) {
		return &NumberLiteral{Value: p.previous().Literal}, nil
	}

	if p.match(TOKEN_IDENTIFIER) {
		return &Identifier{Name: p.previous().Lexeme}, nil
	}

	if p.match(TOKEN_VEC3) {
		if !p.match(TOKEN_LPAREN) {
			return nil, p.error("Expected '(' after 'vec3'")
		}

		x, err := p.expression()
		if err != nil {
			return nil, err
		}

		if !p.match(TOKEN_COMMA) {
			return nil, p.error("Expected ',' after x component")
		}

		y, err := p.expression()
		if err != nil {
			return nil, err
		}

		if !p.match(TOKEN_COMMA) {
			return nil, p.error("Expected ',' after y component")
		}

		z, err := p.expression()
		if err != nil {
			return nil, err
		}

		if !p.match(TOKEN_RPAREN) {
			return nil, p.error("Expected ')' after z component")
		}

		return &Vec3Literal{X: x, Y: y, Z: z}, nil
	}

	if p.match(TOKEN_LPAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		if !p.match(TOKEN_RPAREN) {
			return nil, p.error("Expected ')' after expression")
		}
		return expr, nil
	}

	return nil, p.error(fmt.Sprintf("Unexpected token: %s", p.current.Type))
}

// Helper functions

func (p *Parser) match(types ...TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(t TokenType) bool {
	return p.current.Type == t
}

func (p *Parser) advance() {
	p.prev = p.current
	p.current = p.peek
	p.peek = p.lexer.ScanToken()
}

func (p *Parser) previous() Token {
	return p.prev
}

func (p *Parser) error(message string) error {
	return fmt.Errorf("Parse error at line %d: %s (token: %s '%s')",
		p.current.Line, message, p.current.Type, p.current.Lexeme)
}
