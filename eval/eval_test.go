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

func evaluate(input string) object.Object {
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()
	obj := Eval(program)
	return obj
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
