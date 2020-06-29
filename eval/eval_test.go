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

func evaluate(input string) object.Object {
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()
	obj := Eval(program)
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
