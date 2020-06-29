package eval

import (
	"github.com/HakanSunay/gohil/object"
	"github.com/HakanSunay/gohil/syntaxtree"
)

var (
	Null  = &object.Null{}
	True  = &object.Boolean{Value: true}
	False = &object.Boolean{Value: false}
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
		return parseToBooleanInstance(node.Value)
	// hil supports 2 prefix operators: ! (excl. Mark / Bang) and - (minus)
	case *syntaxtree.PrefixExpr:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)

	case *syntaxtree.InfixExpr:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(node.Operator, left, right)
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

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalNegativeValueExpression(right)
	default:
		return nil
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case True:
		return False
	case False:
		return True
	case Null:
		return True
	default:
		return False
	}
}

func evalNegativeValueExpression(right object.Object) object.Object {
	if right.Type() != object.IntegerObject {
		// TODO: log here
		return nil
	}

	return &object.Integer{Value: -(right.(*object.Integer).Value)}
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == object.IntegerObject && right.Type() == object.IntegerObject:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.BooleanObject && right.Type() == object.BooleanObject:
		return evalBooleanInfixExpression(operator, left, right)
	default:
		return nil
	}
}

func evalIntegerInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	// type assertions should have already been done using the Type method
	// still might need to guard for panics
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "==":
		return parseToBooleanInstance(leftVal == rightVal)
	case "!=":
		return parseToBooleanInstance(leftVal != rightVal)
	case ">":
		return parseToBooleanInstance(leftVal > rightVal)
	case "<":
		return parseToBooleanInstance(leftVal < rightVal)
	default:
		return nil
	}
}

func evalBooleanInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch operator {
	case "==":
		return parseToBooleanInstance(left == right)
	case "!=":
		return parseToBooleanInstance(left != right)
	default:
		return Null
	}
}

func parseToBooleanInstance(p bool) *object.Boolean {
	if p {
		return True
	}
	return False
}
