package compiler

import (
	"fmt"
	"jedil/pkg/bytecode"
)

// Compiler converts AST to bytecode instructions
type Compiler struct {
	instructions []bytecode.Instruction
	variables    map[string]int // variable name -> stack offset
	stackDepth   int
}

// NewCompiler creates a new compiler
func NewCompiler() *Compiler {
	return &Compiler{
		instructions: make([]bytecode.Instruction, 0),
		variables:    make(map[string]int),
		stackDepth:   0,
	}
}

// Compile compiles a program to bytecode
func (c *Compiler) Compile(program *Program) ([]bytecode.Instruction, error) {
	for _, stmt := range program.Statements {
		if err := c.compileStmt(stmt); err != nil {
			return nil, err
		}
	}

	// Add HALT at the end
	c.emit(bytecode.OP_HALT, 0)

	return c.instructions, nil
}

// Statement compilation

func (c *Compiler) compileStmt(stmt Stmt) error {
	switch s := stmt.(type) {
	case *VarDecl:
		return c.compileVarDecl(s)
	case *ReturnStmt:
		return c.compileReturnStmt(s)
	case *FnDecl:
		return fmt.Errorf("functions not yet implemented")
	default:
		return fmt.Errorf("unknown statement type")
	}
}

func (c *Compiler) compileVarDecl(stmt *VarDecl) error {
	// Compile the value expression
	if err := c.compileExpr(stmt.Value); err != nil {
		return err
	}

	// Store the variable at current stack depth
	c.variables[stmt.Name] = c.stackDepth
	c.stackDepth++

	return nil
}

func (c *Compiler) compileReturnStmt(stmt *ReturnStmt) error {
	// Compile the return value expression
	return c.compileExpr(stmt.Value)
}

// Expression compilation

func (c *Compiler) compileExpr(expr Expr) error {
	switch e := expr.(type) {
	case *NumberLiteral:
		return c.compileNumber(e)
	case *Identifier:
		return c.compileIdentifier(e)
	case *BinaryOp:
		return c.compileBinaryOp(e)
	case *UnaryOp:
		return c.compileUnaryOp(e)
	case *Vec3Literal:
		return c.compileVec3Literal(e)
	case *CallExpr:
		return c.compileCallExpr(e)
	default:
		return fmt.Errorf("unknown expression type")
	}
}

func (c *Compiler) compileNumber(expr *NumberLiteral) error {
	c.emit(bytecode.OP_PUSH, expr.Value)
	return nil
}

func (c *Compiler) compileIdentifier(expr *Identifier) error {
	// Look up the variable
	offset, ok := c.variables[expr.Name]
	if !ok {
		return fmt.Errorf("undefined variable: %s", expr.Name)
	}

	// Emit LOAD instruction with stack offset
	c.emit(bytecode.OP_LOAD, float64(offset))
	return nil
}

func (c *Compiler) compileBinaryOp(expr *BinaryOp) error {
	// Compile left operand
	if err := c.compileExpr(expr.Left); err != nil {
		return err
	}

	// Compile right operand
	if err := c.compileExpr(expr.Right); err != nil {
		return err
	}

	// Emit the operation
	switch expr.Op {
	case TOKEN_PLUS:
		// Need to figure out if this is scalar or vector addition
		// For now, assume scalar (we'll detect type later)
		c.emit(bytecode.OP_ADD, 0)
	case TOKEN_MINUS:
		c.emit(bytecode.OP_SUB, 0)
	case TOKEN_STAR:
		c.emit(bytecode.OP_MUL, 0)
	case TOKEN_SLASH:
		c.emit(bytecode.OP_DIV, 0)
	default:
		return fmt.Errorf("unknown binary operator: %v", expr.Op)
	}

	return nil
}

func (c *Compiler) compileUnaryOp(expr *UnaryOp) error {
	// Compile the operand
	if err := c.compileExpr(expr.Expr); err != nil {
		return err
	}

	// Emit the operation
	switch expr.Op {
	case TOKEN_MINUS:
		// Negate: multiply by -1
		c.emit(bytecode.OP_PUSH, -1.0)
		c.emit(bytecode.OP_MUL, 0)
	default:
		return fmt.Errorf("unknown unary operator: %v", expr.Op)
	}

	return nil
}

func (c *Compiler) compileVec3Literal(expr *Vec3Literal) error {
	// Compile x, y, z components
	if err := c.compileExpr(expr.X); err != nil {
		return err
	}
	if err := c.compileExpr(expr.Y); err != nil {
		return err
	}
	if err := c.compileExpr(expr.Z); err != nil {
		return err
	}

	// Emit VEC3 instruction to combine them
	c.emit(bytecode.OP_VEC3, 0)

	return nil
}

func (c *Compiler) compileCallExpr(expr *CallExpr) error {
	// Handle built-in functions
	switch expr.Callee {
	case "cross":
		if len(expr.Args) != 2 {
			return fmt.Errorf("cross() expects 2 arguments, got %d", len(expr.Args))
		}
		if err := c.compileExpr(expr.Args[0]); err != nil {
			return err
		}
		if err := c.compileExpr(expr.Args[1]); err != nil {
			return err
		}
		c.emit(bytecode.OP_VCROSS, 0)
		return nil

	case "dot":
		if len(expr.Args) != 2 {
			return fmt.Errorf("dot() expects 2 arguments, got %d", len(expr.Args))
		}
		if err := c.compileExpr(expr.Args[0]); err != nil {
			return err
		}
		if err := c.compileExpr(expr.Args[1]); err != nil {
			return err
		}
		c.emit(bytecode.OP_VMUL, 0) // OP_VMUL is the dot product
		return nil

	case "mag":
		if len(expr.Args) != 1 {
			return fmt.Errorf("mag() expects 1 argument, got %d", len(expr.Args))
		}
		if err := c.compileExpr(expr.Args[0]); err != nil {
			return err
		}
		c.emit(bytecode.OP_VMAG, 0)
		return nil

	default:
		return fmt.Errorf("unknown function: %s", expr.Callee)
	}
}

// Helper to emit an instruction
func (c *Compiler) emit(op bytecode.OpCode, args float64) {
	c.instructions = append(c.instructions, bytecode.Instruction{
		Op:   op,
		Args: args,
	})
}

// CompileSource is a convenience function that lexes, parses, and compiles in one go
func CompileSource(source string) ([]bytecode.Instruction, error) {
	// Parse
	parser := NewParser(source)
	program, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	// Compile
	compiler := NewCompiler()
	return compiler.Compile(program)
}
