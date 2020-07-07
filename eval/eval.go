package eval

import (
	"fmt"

	"github.com/HakanSunay/gohil/object"
	"github.com/HakanSunay/gohil/syntaxtree"
)

var (
	Null  = &object.Null{}
	True  = &object.Boolean{Value: true}
	False = &object.Boolean{Value: false}
)

func Eval(node syntaxtree.Node, environment *object.Environment) object.Object {
	switch node := node.(type) {
	// Statements:
	case *syntaxtree.Program:
		return evalProgram(node.Statements, environment) // start traversing the program tree
	case *syntaxtree.ExpressionStmt:
		return Eval(node.Expression, environment)
	case *syntaxtree.BlockStmt:
		return evalBlockStatement(node, environment)
	case *syntaxtree.ReturnStmt:
		val := Eval(node.ReturnValue, environment)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *syntaxtree.LetStmt:
		val := Eval(node.Value, environment)
		if isError(val) {
			return val
		}
		environment.Set(node.Name.Value, val)

	// Expressions:
	case *syntaxtree.Identifier:
		return evalIdentifier(node, environment)
	case *syntaxtree.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *syntaxtree.StringLiteral:
		return &object.String{Value: node.Value}
	case *syntaxtree.BooleanLiteral:
		return parseToBooleanInstance(node.Value)
	// hil supports 2 prefix operators: ! (excl. Mark / Bang) and - (minus)
	case *syntaxtree.PrefixExpr:
		right := Eval(node.Right, environment)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *syntaxtree.InfixExpr:
		left := Eval(node.Left, environment)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, environment)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *syntaxtree.IfExpr:
		return evalIfExpression(node, environment)
	case *syntaxtree.CallExpr:
		function := Eval(node.Function, environment)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, environment)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)
	case *syntaxtree.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{
			Parameters: params,
			Body:       body,
			Env:        environment,
		}
	}

	return nil
}

func evalExpressions(exprs []syntaxtree.Expr, env *object.Environment) []object.Object {
	var result []object.Object

	// also evaluation from LEFT to RIGHT
	for _, e := range exprs {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			// this ensures the error check for len 1
			return []object.Object{evaluated}
		}

		result = append(result, evaluated)
	}

	return result
}

func evalProgram(statements []syntaxtree.Stmt, environment *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range statements {
		// last evaluated statement will end up as the result
		result = Eval(stmt, environment)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(block *syntaxtree.BlockStmt, environment *object.Environment) object.Object {
	var result object.Object
	for _, statement := range block.Statements {
		result = Eval(statement, environment)
		if result != nil {
			rt := result.Type()
			if rt == object.ReturnValueObject || rt == object.ErrorObject {
				return result
			}
		}
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
		return newError("unknown operator: %s%s", operator, right.Type())
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
		return newError("unknown operator: -%s", right.Type())
	}

	return &object.Integer{Value: -(right.(*object.Integer).Value)}
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == object.IntegerObject && right.Type() == object.IntegerObject:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.BooleanObject && right.Type() == object.BooleanObject:
		return evalBooleanInfixExpression(operator, left, right)
	case left.Type() == object.StringObject && right.Type() == object.StringObject:
		return evalStringInfixExpression(operator, left, right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
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
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	if operator != "+" {
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}

	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	return &object.String{Value: leftVal + rightVal}
}

func evalBooleanInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch operator {
	case "==":
		return parseToBooleanInstance(left == right)
	case "!=":
		return parseToBooleanInstance(left != right)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func parseToBooleanInstance(p bool) *object.Boolean {
	if p {
		return True
	}
	return False
}

func evalIfExpression(node *syntaxtree.IfExpr, environment *object.Environment) object.Object {
	condition := Eval(node.Condition, environment)
	if isError(condition) {
		return condition
	}

	// This is referred to as being "truthy"
	// this means that we can evaluate expr like if 5 { ... }
	if truthy := condition != Null && condition != False; truthy {
		return Eval(node.Consequence, environment)
	} else if node.Alternative != nil {
		return Eval(node.Alternative, environment)
	} else {
		return Null
	}
}

func evalIdentifier(node *syntaxtree.Identifier, environment *object.Environment) object.Object {
	val, ok := environment.Get(node.Value)
	if !ok {
		return newError("identifier not found: %s", node.Value)
	}

	return val
}

func newError(format string, args ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, args...)}
}

// isError is used when checking for errors whenever we call Eval inside of Eval,
// in order to stop errors from being passed around and then bubbling up far away from their origin
func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ErrorObject
	}
	return false
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	function, ok := fn.(*object.Function)
	if !ok {
		return newError("not a function: %s", fn.Type())
	}

	extendedEnv := extendFunctionEnv(function, args)
	evaluated := Eval(function.Body, extendedEnv)
	return unwrapReturnValue(evaluated)
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)
	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}
	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}
