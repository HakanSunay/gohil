package syntaxtree

import (
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

func (l *LetStmt) GetTokenLiteral() string {return l.Token.Literal }

func (l *LetStmt) stmtNode() {}
