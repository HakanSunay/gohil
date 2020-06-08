package token

// Token represents a token abstraction that is used in lexer logic
type Token struct {
	Type    Type
	Literal string
}

// Set sets the fields of the token type
func (t *Token) Set(typ Type, literal byte) {
	t.Type = typ
	t.Literal = string(literal)
}
