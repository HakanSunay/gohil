package parser

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/HakanSunay/gohil/lexer"
	"github.com/HakanSunay/gohil/syntaxtree"
	"github.com/HakanSunay/gohil/token"
)

// TODO: try to apply DRY to test cases

func TestLetStatements(t *testing.T) {
	// TODO: add tests for bool and string as well, not only integer literals
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
		expectedValue      int
	}{
		{"x", 6},
		{"y", 77},
		{"zzz", 888},
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
		if val, _ := letStmt.Value.(*syntaxtree.IntegerLiteral); val.Value != tt.expectedValue {
			t.Errorf("expected statement value to be %d, but got %d", tt.expectedValue, val.Value)
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

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int
	}{
		{"!6;", "!", 6},
		{"-66;", "-", 66},
	}
	for _, tt := range prefixTests {
		l := lexer.NewLexer(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		if len(program.Statements) != 1 {
			t.Fatalf("expected len 1, but got: %d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*syntaxtree.ExpressionStmt)
		if !ok {
			t.Fatalf("type asserting to ExpressionStmt failed, got :%v",
				reflect.TypeOf(program.Statements[0]))
		}

		exp, ok := stmt.Expression.(*syntaxtree.PrefixExpr)
		if !ok {
			t.Fatalf("type asserting to PrefixExpr failed, got :%v",
				reflect.TypeOf(stmt.Expression))
		}

		if exp.Operator != tt.operator {
			t.Fatalf("expected operator: %s, but got %s",
				tt.operator, exp.Operator)
		}

		if !testIntegerLiteral(t, exp.Right, tt.integerValue) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  int
		operator   string
		rightValue int
	}{
		{"6 + 6;", 6, "+", 6},
		{"6 - 6;", 6, "-", 6},
		{"6 * 6;", 6, "*", 6},
		{"6 / 6;", 6, "/", 6},
		{"6 > 6;", 6, ">", 6},
		{"6 < 6;", 6, "<", 6},
		{"6 == 6;", 6, "==", 6},
		{"6 != 6;", 6, "!=", 6},
	}
	for _, tt := range infixTests {
		l := lexer.NewLexer(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		if len(program.Statements) != 1 {
			t.Fatalf("expected len 1, but got: %d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*syntaxtree.ExpressionStmt)
		if !ok {
			t.Fatalf("type asserting to ExpressionStmt failed, got :%v",
				reflect.TypeOf(program.Statements[0]))
		}

		exp, ok := stmt.Expression.(*syntaxtree.InfixExpr)
		if !ok {
			t.Fatalf("type asserting to InfixExpr failed, got :%v",
				reflect.TypeOf(stmt.Expression))
		}

		if !testIntegerLiteral(t, exp.Left, tt.leftValue) {
			return
		}

		if exp.Operator != tt.operator {
			t.Fatalf("expected %s, but got: %s",
				tt.operator, exp.Operator)
		}

		if !testIntegerLiteral(t, exp.Right, tt.rightValue) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, il syntaxtree.Expr, value int) bool {
	integer, ok := il.(*syntaxtree.IntegerLiteral)
	if !ok {
		t.Errorf("unable to type assert to IntegerLiteral, got %v",
			reflect.TypeOf(il))
		return false
	}

	if integer.Value != value {
		t.Errorf("expected %d, but got %d", value, integer.Value)
		return false
	}

	if integer.GetTokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("expected %v, but got %v", value,
			integer.GetTokenLiteral())
		return false
	}

	return true
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-x * y",
			"((-x) * y)",
		},
		{
			"!-x",
			"(!(-x))",
		},
		{
			"x + y + z",
			"((x + y) + z)",
		},
		{
			"x + y - z",
			"((x + y) - z)",
		},
		{
			"x * y * z",
			"((x * y) * z)",
		},
		{
			"x * y / z",
			"((x * y) / z)",
		},
		{
			"x + y / z",
			"(x + (y / z))",
		},
		{
			"x + y * z + d / e - f",
			"(((x + (y * z)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
	}

	for _, tt := range tests {
		l := lexer.NewLexer(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		// String will wrap the literals in parentheses
		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected %s, but got %s", tt.expected, actual)
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()

	if len(program.Statements) != 1 {
		t.Fatalf("expected len (%d), but got %d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*syntaxtree.ExpressionStmt)
	if !ok {
		t.Fatalf("type asserting to ExpressionsStmt failed, got %T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*syntaxtree.IfExpr)
	if !ok {
		t.Fatalf("type asserting to IfExpr failed, got %T", stmt.Expression)
	}

	conditionExpr, ok := exp.Condition.(*syntaxtree.InfixExpr)
	if !ok {
		t.Errorf("type asserting to ExpressionsStmt failed, got %T", exp.Condition)
	}
	if conditionExpr.Left.String() != "x" {
		t.Errorf("expected left side of condition expr to be x, but got %v", conditionExpr.Left.String())
	}
	if conditionExpr.Operator != "<" {
		t.Errorf("expected < operator, but got %v", conditionExpr.Operator)
	}
	if conditionExpr.Right.String() != "y" {
		t.Errorf("expected right side of condition expr to be y, but got %v", conditionExpr.Left.String())
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("expected len 1, but got %d", len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*syntaxtree.ExpressionStmt)
	if !ok {
		t.Fatalf("type asserting to ExpressionsStmt failed, got %T", exp.Consequence.Statements[0])
	}

	if consequence.Expression.String() != "x" {
		t.Errorf("expected consequence expression to be x, but got %v", consequence.Expression.String())
	}

	if exp.Alternative != nil {
		t.Errorf("expected else alternative to be nil, but got %v", exp.Alternative.String())
	}
}

func TestFunctionLiteral(t *testing.T) {
	input := `fn(x, y, z) { x + y + z; }`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()

	if len(program.Statements) != 1 {
		t.Fatalf("expected len (%d), but got %d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*syntaxtree.ExpressionStmt)
	if !ok {
		t.Fatalf("type asserting to ExpressionsStmt failed, got %T", program.Statements[0])
	}

	fn, ok := stmt.Expression.(*syntaxtree.FunctionLiteral)
	if !ok {
		t.Fatalf("type asserting to FunctionLiteral failed, got %T", stmt.Expression)
	}

	if fn.Token.Type != token.Function {
		t.Errorf("expected %v, but got %v", token.Function, fn.Token.Type)
	}

	if fn.Parameters[0].GetTokenLiteral() != "x" {
		t.Errorf("expected x as first parameter literal, but got %v", fn.Parameters[0].GetTokenLiteral())
	}

	if fn.Parameters[1].GetTokenLiteral() != "y" {
		t.Errorf("expected y as second parameter literal, but got %v", fn.Parameters[1].GetTokenLiteral())
	}

	if fn.Parameters[2].GetTokenLiteral() != "z" {
		t.Errorf("expected z as first parameter literal, but got %v", fn.Parameters[2].GetTokenLiteral())
	}

	if len(fn.Body.Statements) != 1 {
		t.Errorf("expected 1 statement, got %d", len(fn.Body.Statements))
	}

	// we know that infix parsing operations are enclosed from right to left,
	// when they have the same precedence
	if fn.Body.String() != "((x + y) + z)" {
		t.Errorf("expected ((x + y) + z), but got %v", fn.Body.String())
	}
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {};", expectedParams: []string{}},
		{input: "fn(x) {};", expectedParams: []string{"x"}},
		{input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}
	for _, tt := range tests {
		l := lexer.NewLexer(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		stmt := program.Statements[0].(*syntaxtree.ExpressionStmt)
		function := stmt.Expression.(*syntaxtree.FunctionLiteral)
		if len(function.Parameters) != len(tt.expectedParams) {
			t.Errorf("expected len %d, but got %d", len(tt.expectedParams), len(function.Parameters))
		}
		for i, ident := range tt.expectedParams {
			id := function.Parameters[i]
			if id.GetTokenLiteral() != ident {
				t.Errorf("expected %v, but got %v", ident, id.GetTokenLiteral())
			}
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"
	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()

	if len(program.Statements) != 1 {
		t.Fatalf("expected statement len %d, but got %d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*syntaxtree.ExpressionStmt)
	if !ok {
		t.Fatalf("stmt is not syntaxtree.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*syntaxtree.CallExpr)
	if !ok {
		t.Fatalf("stmt.Expression is not syntaxtree.CallExpression. got=%T",
			stmt.Expression)
	}

	if exp.Function.GetTokenLiteral() != "add" {
		t.Errorf("expected token literal %v, but got %v", "add", exp.Function.GetTokenLiteral())
	}

	if len(exp.Arguments) != 3 {
		t.Fatalf("expected length of arguments %d, but got=%d", 3, len(exp.Arguments))
	}

	if exp.Arguments[0].String() != "1" {
		t.Errorf("expected 1, but got %v", exp.Arguments[0].String())
	}

	firstExpr, ok := exp.Arguments[1].(*syntaxtree.InfixExpr)
	if !ok {
		t.Fatalf("unable to type assert to InfixExpr, got %T", exp.Arguments[1])
	}

	if firstExpr.Left.String() != "2" {
		t.Errorf("expected 2, but got %v", firstExpr.Left.String())
	}

	if firstExpr.Operator != "*" {
		t.Errorf("expected *, but got %v", firstExpr.Operator)
	}

	if firstExpr.Right.String() != "3" {
		t.Errorf("expected 3, but got %v", firstExpr.Right.String())
	}

	secondExpr, ok := exp.Arguments[2].(*syntaxtree.InfixExpr)
	if !ok {
		t.Fatalf("unable to type assert to InfixExpr, got %T", exp.Arguments[2])
	}

	if secondExpr.Left.String() != "4" {
		t.Errorf("expected 4, but got %v", secondExpr.Left.String())
	}

	if secondExpr.Operator != "+" {
		t.Errorf("expected +, but got %v", secondExpr.Operator)
	}

	if secondExpr.Right.String() != "5" {
		t.Errorf("expected 5, but got %v", secondExpr.Right.String())
	}
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world";`
	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()

	stmt := program.Statements[0].(*syntaxtree.ExpressionStmt)
	literal, ok := stmt.Expression.(*syntaxtree.StringLiteral)
	if !ok {
		t.Fatalf("expected type StringLiteral, but got %T", stmt.Expression)
	}
	if literal.Value != "hello world" {
		t.Errorf("expected literal value %v, but got %v", "hello world", literal.Value)
	}
}
