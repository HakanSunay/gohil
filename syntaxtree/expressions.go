package syntaxtree

import (
	"github.com/HakanSunay/gohil/token"
	"strings"
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
