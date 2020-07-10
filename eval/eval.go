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
	case *syntaxtree.ArrayLiteral:
		elements := evalExpressions(node.Elements, environment)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}
	case *syntaxtree.HashLiteral:
		return evalHashLiteral(node, environment)
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
	case *syntaxtree.IndexExpression:
		left := Eval(node.Left, environment)
		if isError(left) {
			return left
		}

		index := Eval(node.Index, environment)
		if isError(index) {
			return index
		}

		return evalIndexExpression(left, index)
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
	// check for existence in env
	if val, ok := environment.Get(node.Value); ok {
		return val
	}

	// check if it is a builtin function
	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	// neither env var nor builtin
	return newError("identifier not found: " + node.Value)
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
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function: %s", fn.Type())
	}
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

// evalIndexExpressions handles evaluation logic for all kinds of index operations
// the main purpose is to switch on parameter types to decide what logic to apply
func evalIndexExpression(left object.Object, index object.Object) object.Object {
	switch {
	// arr[INTEGER]
	case left.Type() == object.ArrayObject && index.Type() == object.IntegerObject:
		return evalArrayIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalArrayIndexExpression(arr object.Object, index object.Object) object.Object {
	// type assertion wont fail, guaranteed before
	arrayObject := arr.(*object.Array)

	// TODO: can do some python magic like returning last element if -1 is the index
	// or some special case if the index value is 999 999 always return the fist element
	i := index.(*object.Integer).Value
	max := len(arrayObject.Elements) - 1
	if i < 0 || i > max {
		return Null
	}

	return arrayObject.Elements[i]
}

func evalHashLiteral(node *syntaxtree.HashLiteral, environment *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)
	for keyNode, valueNode := range node.Pairs {
		// evaluate the key
		key := Eval(keyNode, environment)
		if isError(key) {
			return key
		}

		// verify it is hashable
		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", key.Type())
		}

		// evaluate the value
		value := Eval(valueNode, environment)
		if isError(value) {
			return value
		}

		// get the hash of the key
		hashed := hashKey.HashKey()

		// add to pairs map as HashPair object
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return &object.Hash{Pairs: pairs}
}
