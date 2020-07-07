package syntaxtree

import (
	"strings"

	"github.com/HakanSunay/gohil/token"
)

// Identifier represents identifiers that are used in statements and expressions.
// x is an identifier and to represent it we need a token - token.Identifier
// and a value in our case, the value is "x"
type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) GetTokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) String() string {
	return i.Value
}

// identifier doesn't produce values, why is it implementing the Expr iface?
// we will have functions that produce values, which will be assigned to identifiers.
func (i *Identifier) exprNode() {}

type IntegerLiteral struct {
	Token token.Token
	Value int
}

func (il *IntegerLiteral) GetTokenLiteral() string {
	return il.Token.Literal
}

func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}

func (il *IntegerLiteral) exprNode() {}

type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (b *BooleanLiteral) GetTokenLiteral() string {
	return b.Token.Literal
}

func (b *BooleanLiteral) String() string {
	return b.Token.Literal
}

func (b *BooleanLiteral) exprNode() {}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) GetTokenLiteral() string {
	return sl.Token.Literal
}

func (sl *StringLiteral) String() string {
	return sl.Token.Literal
}

func (sl *StringLiteral) exprNode() {}

// PrefixExpr describes prefix expressions in gohil.
// There are 2 types of prefix expressions in the language: ! and -
// E.g: -66
// Token: -
// Operation: -
// Right: 66
type PrefixExpr struct {
	Token    token.Token // The prefix token, e.g. ! or -
	Operator string
	Right    Expr
}

func (p *PrefixExpr) GetTokenLiteral() string {
	return p.Token.Literal
}

func (p *PrefixExpr) String() string {
	var builder strings.Builder

	builder.WriteString("(")

	builder.WriteString(p.Operator)
	builder.WriteString(p.Right.String())

	builder.WriteString(")")

	return builder.String()
}

func (p *PrefixExpr) exprNode() {}

// InfixExpr describes infix expressions.
// There are many infix expressions supported by gohil.
// All of the arithmetic operations are considered infix expressions.
// E.g: 6 + 8
// Should result in:
// Token: +
// Left: IntegerLiteral(6)
// Operator: +
// Right: IntegerLiteral(8)
type InfixExpr struct {
	Token token.Token

	Left     Expr
	Operator string
	Right    Expr
}

func (i *InfixExpr) GetTokenLiteral() string {
	return i.Token.Literal
}

func (i *InfixExpr) String() string {
	var builder strings.Builder

	builder.WriteString("(")
	builder.WriteString(i.Left.String())

	builder.WriteString(" ")
	builder.WriteString(i.Operator)
	builder.WriteString(" ")

	builder.WriteString(i.Right.String())
	builder.WriteString(")")

	return builder.String()
}

func (i *InfixExpr) exprNode() {}

type IfExpr struct {
	Token       token.Token // if
	Condition   Expr
	Consequence *BlockStmt
	Alternative *BlockStmt
}

func (ie *IfExpr) GetTokenLiteral() string {
	return ie.Token.Literal
}

func (ie *IfExpr) String() string {
	var builder strings.Builder

	builder.WriteString("if")
	builder.WriteString(ie.Condition.String())
	builder.WriteString(" ")
	builder.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		builder.WriteString("else ")
		builder.WriteString(ie.Alternative.String())
	}

	return builder.String()
}

func (ie *IfExpr) exprNode() {}

type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStmt // reminder: 1 block statement has many statements
}

func (f *FunctionLiteral) String() string {
	var builder strings.Builder

	var params []string
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	// fn
	builder.WriteString(f.GetTokenLiteral())

	// params (x, y, z)
	builder.WriteString("(")
	builder.WriteString(strings.Join(params, ", "))
	builder.WriteString(") ")

	// function body
	builder.WriteString(f.Body.String())

	return builder.String()
}

func (f *FunctionLiteral) GetTokenLiteral() string {
	return f.Token.Literal
}

func (f *FunctionLiteral) exprNode() {}

// CallExpr identifies a callable expression.
// call expressions are of this structure:
// <expression>(<comma separated expressions>)
// sum(1, 2)
// sum(1 + 2, 3 + 4)
// fn(x, y) { x + y; }(1, 2)
type CallExpr struct {
	Token     token.Token // '(' left parenthesis
	Function  Expr        // either an identifier or a function literal
	Arguments []Expr
}

func (c *CallExpr) String() string {
	var builder strings.Builder

	var args []string
	for _, a := range c.Arguments {
		args = append(args, a.String())
	}

	builder.WriteString(c.Function.String())
	builder.WriteString("(")
	builder.WriteString(strings.Join(args, ", "))
	builder.WriteString(")")

	return builder.String()
}

func (c *CallExpr) GetTokenLiteral() string {
	return c.Token.Literal
}

func (c *CallExpr) exprNode() {}
