package syntaxtree

import (
	"strings"

	"github.com/HakanSunay/gohil/token"
)

// LetStmt defines a let statement.
// E.g: let x = 6
// This means that we need a token that identifies this statement - token.Let.
// We need an identifier - x.
// We also need a value - 6.
type LetStmt struct {
	Token token.Token
	Name  *Identifier
	Value Expr
}

func (l *LetStmt) String() string {
	var builder strings.Builder

	builder.WriteString(l.GetTokenLiteral())
	builder.WriteString(" ")
	builder.WriteString(l.Name.String())
	builder.WriteString(" = ")

	if l.Value != nil {
		builder.WriteString(l.Value.String())
	}

	builder.WriteString(";")

	return builder.String()
}

func (l *LetStmt) GetTokenLiteral() string {
	return l.Token.Literal
}

func (l *LetStmt) stmtNode() {}

// ReturnStmt defines a return statement.
// E.g: return 6; return keyword and expression.
// This means that we need a token that identifies this statement - token.Return.
// We also need a return value - 6, which is an expression.
type ReturnStmt struct {
	Token       token.Token
	ReturnValue Expr // return expression
}

func (r *ReturnStmt) GetTokenLiteral() string {
	return r.Token.Literal
}

func (r *ReturnStmt) stmtNode() {}

func (r *ReturnStmt) String() string {
	var builder strings.Builder

	builder.WriteString(r.GetTokenLiteral())
	builder.WriteString(" ")

	if r.ReturnValue != nil {
		builder.WriteString(r.ReturnValue.String())
	}

	builder.WriteString(";")

	return builder.String()
}

// ExpressionStmt defines an expression statement.
// The previous 2 types were either only expr or stmt, but now we have both.
// Most scripting languages support this type of statements, so will gohil.
// E.g:
// let x = 6; // we said that this was a let statement
// x + 6; // this is an expression statement
// This type implements the Stmt interface, therefore we can use it in
// the Program type, which holds a slice of statements, which in turn means
// that gohil now supports expression statements
type ExpressionStmt struct {
	Token      token.Token
	Expression Expr
}

func (e *ExpressionStmt) GetTokenLiteral() string {
	return e.Token.Literal
}

func (e *ExpressionStmt) String() string {
	if e.Expression != nil {
		return e.Expression.String()
	}

	return ""
}

func (e *ExpressionStmt) stmtNode() {}

// BlockStmt defines a block statement.
// Used in conditional expressions - if, and function definitions
type BlockStmt struct {
	Token      token.Token // the { token
	Statements []Stmt
}

func (bs *BlockStmt) GetTokenLiteral() string {
	return bs.Token.Literal
}
func (bs *BlockStmt) String() string {
	var builder strings.Builder

	for _, s := range bs.Statements {
		builder.WriteString(s.String())
	}

	return builder.String()
}

func (bs *BlockStmt) stmtNode() {}
