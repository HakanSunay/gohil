package syntaxtree

import "github.com/HakanSunay/gohil/token"

// Identifier represents identifiers that are used in statements and expressions.
// x is an identifier and to represent it we need a token - token.Identifier
// and a value in our case, the value is "x"
type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) GetTokenLiteral() string { return i.Token.Literal }

// identifier doesn't produce values, why is it implementing the Expr iface?
// we will have functions that produce values, which will be assigned to identifiers.
func (i *Identifier) exprNode() {}
