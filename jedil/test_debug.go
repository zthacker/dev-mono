package main

import (
	"fmt"
	"jedil/pkg/compiler"
)

func main() {
	source := "return 2.0 + 3.0"
	parser := compiler.NewParser(source)
	program, err := parser.Parse()
	if err != nil {
		fmt.Printf("Parse error: %v\n", err)
		return
	}
	
	// Check the AST
	ret := program.Statements[0].(*compiler.ReturnStmt)
	binop := ret.Value.(*compiler.BinaryOp)
	fmt.Printf("BinaryOp.Op type: %T, value: %v (%d)\n", binop.Op, binop.Op, binop.Op)
	fmt.Printf("TOKEN_PLUS value: %d\n", compiler.TOKEN_PLUS)
}
