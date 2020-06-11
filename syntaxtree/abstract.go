package syntaxtree

type Program struct {
	Statements []Stmt
}

func (p *Program) GetTokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].GetTokenLiteral()
	}

	return ""
}


