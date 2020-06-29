package eval

import (
	"github.com/HakanSunay/gohil/object"
	"github.com/HakanSunay/gohil/syntaxtree"
)

func Eval(node syntaxtree.Node) object.Object {
	switch node := node.(type) {
	// Statements:
	case *syntaxtree.Program:
		return evalStatements(node.Statements) // start traversing the program tree
	case *syntaxtree.ExpressionStmt:
		return Eval(node.Expression)

	// Expressions
	case *syntaxtree.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *syntaxtree.BooleanLiteral:
		// TODO: Performance Boost
		// Performance boost idea: Use same instance for TRUE and FALSE,
		// instead of allocating new instance every time we encounter BooleanLiterals
		return &object.Boolean{Value: node.Value}
	}

	return nil
}

func evalStatements(statements []syntaxtree.Stmt) object.Object {
	var result object.Object

	for _, stmt := range statements {
		// last evaluated statement will end up as the result
		result = Eval(stmt)
	}

	return result
}
