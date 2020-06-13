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
let zzz = 888;
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
		{"zzz"},
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

func TestReturnStatement(t *testing.T) {
	input := `
return x;
return z;
return 7;
`
	l := lexer.NewLexer(input)
	p := NewParser(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("unexpected nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("expected len 3, but got: %d",
			len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*syntaxtree.ReturnStmt)
		if !ok {
			t.Errorf("typer asserting to ReturnStmt failed, got: %v",
				reflect.TypeOf(stmt))
			continue
		}

		if returnStmt.GetTokenLiteral() != "return" {
			t.Errorf("expected return, but got %v",
				returnStmt.GetTokenLiteral())
		}
	}
}

func TestParseIdentifierExpression(t *testing.T) {
	input := "grade;"

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	if len(program.Statements) != 1 {
		t.Errorf("expected len 1, but got: %d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*syntaxtree.ExpressionStmt)
	if !ok {
		t.Fatalf("expected ExpressionStmt type, but got: %v",
			reflect.TypeOf(program.Statements[0]))
	}

	ident, ok := stmt.Expression.(*syntaxtree.Identifier)
	if !ok {
		t.Fatalf("Unable to type assertion expr (%v) to Identifier",
			stmt.Expression)
	}

	// the semi colon is ignored
	if ident.Value != "grade" {
		t.Errorf("expected grade as identifier value, but got %v",
			ident.Value)
	}

	if ident.GetTokenLiteral() != "grade" {
		t.Errorf("expected grade as token literal, but got %v",
			ident.GetTokenLiteral())
	}
}

func TestParseIntegerLiteralExpression(t *testing.T) {
	input := "6;"

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	if len(program.Statements) != 1 {
		t.Errorf("expected len 1, but got: %d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*syntaxtree.ExpressionStmt)
	if !ok {
		t.Fatalf("expected ExpressionStmt type, but got: %v",
			reflect.TypeOf(program.Statements[0]))
	}

	literal, ok := stmt.Expression.(*syntaxtree.IntegerLiteral)
	if !ok {
		t.Fatalf("Unable to type assertion expr (%v) to IntegerLiteral",
			stmt.Expression)
	}

	if literal.Value != 6 {
		t.Errorf("expected 6 as integer value, but got %d", literal.Value)
	}

	if literal.GetTokenLiteral() != "6" {
		t.Errorf("expected token literal 6, but got %v",
			literal.GetTokenLiteral())
	}
}