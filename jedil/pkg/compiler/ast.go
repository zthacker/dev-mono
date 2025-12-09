package compiler

// AST Node types for the JEDIL language

// Node is the base interface for all AST nodes
type Node interface {
	node()
}

// Expression nodes

type Expr interface {
	Node
	expr()
}

// NumberLiteral represents a numeric constant
type NumberLiteral struct {
	Value float64
}

// Identifier represents a variable reference
type Identifier struct {
	Name string
}

// BinaryOp represents a binary operation (a + b, a * b, etc.)
type BinaryOp struct {
	Left  Expr
	Op    TokenType // TOKEN_PLUS, TOKEN_MINUS, etc.
	Right Expr
}

// UnaryOp represents a unary operation (-x)
type UnaryOp struct {
	Op   TokenType // TOKEN_MINUS
	Expr Expr
}

// CallExpr represents a function call
type CallExpr struct {
	Callee string // Function name
	Args   []Expr
}

// Vec3Literal represents vec3(x, y, z)
type Vec3Literal struct {
	X Expr
	Y Expr
	Z Expr
}

// Statement nodes

type Stmt interface {
	Node
	stmt()
}

// VarDecl represents: let x = expr
type VarDecl struct {
	Name  string
	Value Expr
}

// ReturnStmt represents: return expr
type ReturnStmt struct {
	Value Expr
}

// FnDecl represents a function declaration
type FnDecl struct {
	Name   string
	Params []string
	Body   []Stmt
}

// Program represents the entire program
type Program struct {
	Statements []Stmt
}

// Implement the marker interfaces
func (NumberLiteral) node() {}
func (NumberLiteral) expr()  {}
func (Identifier) node()     {}
func (Identifier) expr()      {}
func (BinaryOp) node()       {}
func (BinaryOp) expr()        {}
func (UnaryOp) node()        {}
func (UnaryOp) expr()         {}
func (CallExpr) node()       {}
func (CallExpr) expr()        {}
func (Vec3Literal) node()    {}
func (Vec3Literal) expr()     {}

func (VarDecl) node()    {}
func (VarDecl) stmt()     {}
func (ReturnStmt) node() {}
func (ReturnStmt) stmt()  {}
func (FnDecl) node()     {}
func (FnDecl) stmt()      {}
func (Program) node()    {}
