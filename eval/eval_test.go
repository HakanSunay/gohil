package eval

import (
	"testing"

	"github.com/HakanSunay/gohil/lexer"
	"github.com/HakanSunay/gohil/object"
	"github.com/HakanSunay/gohil/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input  string
		output int
	}{
		{"0", 0},
		{"6", 6},
		{"1337", 1337},
		// negative values
		{"-6", -6},
		{"-1337", -1337},
		{"6 + 6 + 6 - 10", 8},
		{"2 * 2 * 2 * 2", 16},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}
	for _, tt := range tests {
		evaluatedObj := evaluate(tt.input)
		verifyIntegerObj(t, evaluatedObj, tt.output)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input  string
		output bool
	}{
		{"true", true},
		{"false", false},
		{"6 < 9", true},
		{"6 > 9", false},
		{"6 == 9", false},
		{"6 != 9", true},

		{"9 > 6", true},
		{"9 < 6", false},
		{"9 == 6", false},
		{"9 != 6", true},

		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(6 < 9) == true", true},
		{"(6 < 9) == false", false},
		{"(6 > 9) == true", false},
		{"(6 > 9) == false", true},
	}
	for _, tt := range tests {
		evaluatedObj := evaluate(tt.input)
		verifyBooleanObj(t, evaluatedObj, tt.output)
	}
}

func TestEvalNegateExpression(t *testing.T) {
	tests := []struct {
		input  string
		output bool
	}{
		{"!true", false},
		{"!false", true},
		{"!6", false},
		{"!!true", true},
		{"!!false", false},
		{"!!6", true},
	}
	for _, tt := range tests {
		evaluatedObj := evaluate(tt.input)
		verifyBooleanObj(t, evaluatedObj, tt.output)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (6) { 10 }", 10},
		{"if (6 < 9) { 10 }", 10},
		{"if (6 > 9) { 10 }", nil},
		{"if (6 > 9) { 10 } else { 20 }", 20},
		{"if (6 < 9) { 10 } else { 20 }", 10},
	}
	for _, tt := range tests {
		evaluated := evaluate(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			verifyIntegerObj(t, evaluated, integer)
		} else {
			verifyNullObject(t, evaluated)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"return 6;", 6},
		{"return 16; 9;", 16},
		{"return 6 * 2; 9;", 12},
		{"9; return 2 * 6; 9;", 12},
		// nested return evaluation
		{`if (10 > 1) {
			if (10 > 1) {
				return 6;
			}
		  	return 12;
          }`, 6},
	}
	for _, tt := range tests {
		evaluated := evaluate(tt.input)
		verifyIntegerObj(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input string
		expectedMessage string
	}{
		{
			"5 + true;",
			"type mismatch: Integer + Boolean",
		},
		{
			"5 + true; 5;",
			"type mismatch: Integer + Boolean",
		},
		{
			"-true",
			"unknown operator: -Boolean",
		},
		{
			"true + false;",
			"unknown operator: Boolean + Boolean",
		},
		{
			"5; true + false; 5",
			"unknown operator: Boolean + Boolean",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: Boolean + Boolean",
		},
		{
			`if (10 > 1) {
				if (10 > 1) {
					return true + false;
				}
				return 1;
			}`,
			"unknown operator: Boolean + Boolean",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
	}
	for _, tt := range tests {
		evaluated := evaluate(tt.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)",
				evaluated, evaluated)
			continue
		}
		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q",
				tt.expectedMessage, errObj.Message)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input string
		expected int
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}
	for _, tt := range tests {
		verifyIntegerObj(t, evaluate(tt.input), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; };"
	evaluated := evaluate(input)

	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("expected Function type, but got %T", evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("expected function parameter len 1, but got %v", fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("expected parameter x, but got %v", fn.Parameters[0])
	}

	expectedBody := "(x + 2)"
	if fn.Body.String() != expectedBody {
		t.Fatalf("expected function body %v, but got %v", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input string
		expected int
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
	}
	for _, tt := range tests {
		verifyIntegerObj(t, evaluate(tt.input), tt.expected)
	}
}

func evaluate(input string) object.Object {
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()
	// new env for each test case, so as not to persist the previous state
	environment := object.NewEnvironment()
	obj := Eval(program, environment)
	return obj
}

func verifyBooleanObj(t *testing.T, obj object.Object, expected bool) {
	val, ok := obj.(*object.Boolean)
	if !ok {
		t.Fatalf("Expected Boolean object type, but got %T", obj)
	}
	if val.Value != expected {
		t.Errorf("expected %v, but got %v", expected, val.Value)
	}
}

func verifyIntegerObj(t *testing.T, obj object.Object, expected int) {
	val, ok := obj.(*object.Integer)
	if !ok {
		t.Fatalf("Expected Integer object type, but got %T", obj)
	}
	if val.Value != expected {
		t.Errorf("expected %v, but got %v", expected, val.Value)
	}
}

func verifyNullObject(t *testing.T, obj object.Object) {
	if obj != Null {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
	}
}
