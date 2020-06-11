package parser

import (
	"reflect"
	"testing"

	"github.com/HakanSunay/gohil/lexer"
	"github.com/HakanSunay/gohil/syntaxtree"
)

func TestLetStatements(t *testing.T) {
	input := `
let x = 6;
let y = 77;
let foobar = 888;
`
	l := lexer.NewLexer(input)
	p := NewParser(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("unexpected nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("expected 3 statements, but got: %d",
			len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}
	for i, tt := range tests {
		stmt := program.Statements[i]
		if stmt.GetTokenLiteral() != "let" {
			t.Errorf("expected let as token literal, but got: %v",
				stmt.GetTokenLiteral())
		}
		letStmt, ok := stmt.(*syntaxtree.LetStmt)
		if !ok {
			t.Errorf("type asserting to LetStmt failed, got: %v",
				reflect.TypeOf(stmt))
		}
		if letStmt.Name.Value != tt.expectedIdentifier {
			t.Errorf("expected identifier: %v, but got %v",
				tt.expectedIdentifier, letStmt.Name.Value)
		}
		if letStmt.Name.GetTokenLiteral() != tt.expectedIdentifier {
			t.Errorf("expected token literal: %v, but got %v",
				tt.expectedIdentifier, letStmt.Name.GetTokenLiteral())
		}
	}
}
