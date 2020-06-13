package syntaxtree

import "strings"

type Program struct {
	Statements []Stmt
}

func (p *Program) GetTokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].GetTokenLiteral()
	}

	return ""
}

func (p *Program) String() string {
	var builder strings.Builder

	for _, stm := range p.Statements {
		builder.WriteString(stm.String())
	}

	return builder.String()
}
