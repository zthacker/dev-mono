package compiler

import (
	"fmt"
	"jedil/pkg/bytecode"
)

type FunctionMetadata struct {
	name string
	address int // bytecode instruction offset
	paramCount int // num of params
	returnCount int // num of return values
	localCount int // total local vars (params + locals)
}

// Compiler converts AST to bytecode instructions
type Compiler struct {
	instructions []bytecode.Instruction
	variables    map[string]int // variable name -> stack offset
	stackDepth   int
	
	// function support
	functions map[string]*FunctionMetadata //function table
	currentFunction *FunctionMetadata // currently compiling
	localVars map[string]int // local scope vars
	inFunction bool // inside func?
}

// NewCompiler creates a new compiler
func NewCompiler() *Compiler {
	return &Compiler{
		instructions: make([]bytecode.Instruction, 0),
		variables:    make(map[string]int),
		functions: make(map[string]*FunctionMetadata),
		localVars: make(map[string]int),
		stackDepth:   0,
	}
}

// Compile compiles a program to bytecode
func (c *Compiler) Compile(program *Program) ([]bytecode.Instruction, error) {
	// PASS 1: Collect function declarations
	var functions []*FnDecl
	var mainCode []Stmt

	for _, stmt := range program.Statements {
		if fnDecl, ok := stmt.(*FnDecl); ok {
			functions = append(functions, fnDecl)
			c.functions[fnDecl.Name] = &FunctionMetadata{
				name:        fnDecl.Name,
				address:     -1, // set in pass 2
				paramCount:  len(fnDecl.Params),
				returnCount: 1, // default is single return -- extended later
			}
		} else {
			mainCode = append(mainCode, stmt)
		}
	}

	// Emit JMP to skip over all function bodies
	jmpIndex := len(c.instructions)
	c.emit(bytecode.OP_JMP, 0) // placeholder

	// PASS 2a: Compile all functions
	for _, fnDecl := range functions {
		if err := c.compileFnDecl(fnDecl); err != nil {
			return nil, err
		}
	}

	// Patch JMP to point here (after all functions)
	c.instructions[jmpIndex].Args = float64(len(c.instructions))

	// PASS 2b: Compile main code
	for _, stmt := range mainCode {
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
		return c.compileFnDecl(s)
	default:
		return fmt.Errorf("unknown statement type")
	}
}

func (c *Compiler) compileFnDecl(stmt *FnDecl) error {
    // Get function metadata
    fnMeta := c.functions[stmt.Name]

    // Record bytecode address where function starts
    fnMeta.address = len(c.instructions)

    // Enter function scope
    c.inFunction = true
    c.currentFunction = fnMeta
    c.localVars = make(map[string]int)
    prevStackDepth := c.stackDepth

    // Parameters are already on stack (pushed by caller)
    // They start at position 0 within the function's frame
    // The basePointer at runtime will point to where args begin
    c.stackDepth = 0
    for i, paramName := range stmt.Params {
        c.localVars[paramName] = i
    }
    c.stackDepth = len(stmt.Params)

    // Compile function body
    for _, bodyStmt := range stmt.Body {
        if err := c.compileStmt(bodyStmt); err != nil {
            return err
        }
    }

    // Add implicit return if function doesn't end with return
    lastIsReturn := false
    if len(stmt.Body) > 0 {
        _, lastIsReturn = stmt.Body[len(stmt.Body)-1].(*ReturnStmt)
    }
    if !lastIsReturn {
        // Implicit return: push nil/zero and return
        c.emit(bytecode.OP_PUSH, 0)
        c.emit(bytecode.OP_RET, float64(fnMeta.returnCount))
    }

    // Exit function scope
    c.inFunction = false
    c.currentFunction = nil
    c.stackDepth = prevStackDepth

    return nil
}

func (c *Compiler) compileReturnStmt(stmt *ReturnStmt) error {
    // Compile return value expression
    if err := c.compileExpr(stmt.Value); err != nil {
        return err
    }

    // Only emit OP_RET if we're inside a function
    if c.inFunction {
        // Emit RET with return count
        c.emit(bytecode.OP_RET, float64(c.currentFunction.returnCount))
    }
    // Top-level returns just leave the value on the stack

    return nil
}

func (c *Compiler) compileVarDecl(stmt *VarDecl) error {
    // Compile the value expression
    if err := c.compileExpr(stmt.Value); err != nil {
        return err
    }

    // Store variable in appropriate scope
    if c.inFunction {
        // Local variable
        c.localVars[stmt.Name] = c.stackDepth
    } else {
        // Global variable
        c.variables[stmt.Name] = c.stackDepth
    }
    c.stackDepth++

    return nil
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
	// Check local variables first (if in function)
    if c.inFunction {
        if offset, ok := c.localVars[expr.Name]; ok {
            c.emit(bytecode.OP_LOAD, float64(offset))
            return nil
        }
    }
	
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
	// Check if it's a user-defined function
    if fnMeta, ok := c.functions[expr.Callee]; ok {
        // Validate parameter count
        if len(expr.Args) != fnMeta.paramCount {
            return fmt.Errorf("function %s expects %d arguments, got %d",
                expr.Callee, fnMeta.paramCount, len(expr.Args))
        }

        // Compile arguments (left to right, pushed onto stack)
        for _, arg := range expr.Args {
            if err := c.compileExpr(arg); err != nil {
                return err
            }
        }

        // Push parameter count before CALL so VM knows where basePointer should be
        c.emit(bytecode.OP_PUSH, float64(len(expr.Args)))

        // Emit CALL instruction
        c.emit(bytecode.OP_CALL, float64(fnMeta.address))

        // Note: Return values are left on stack as temporary expression results
        // stackDepth is only incremented when stored to a variable in compileVarDecl

        return nil
    }
	
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
