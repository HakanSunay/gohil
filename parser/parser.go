package parser

import (
	"fmt"

	"github.com/HakanSunay/gohil/lexer"
	"github.com/HakanSunay/gohil/syntaxtree"
	"github.com/HakanSunay/gohil/token"
)

// Parser repeatedly calls lexer's NextToken to apply logic onto it.
// We need the current and the next token for every evaluation,
// because future knowledge is crucial during evaluation.
// E.g:
// currentToken is 6
// nextToken could be a semi-colon or the beginning of an arithmetic operation
type Parser struct {
	lxr *lexer.Lexer

	currentToken token.Token
	nextToken    token.Token

	errors []string
}

// NewParser is the constructor for the Parser type
func NewParser(lxr *lexer.Lexer) *Parser {
	parser := &Parser{lxr: lxr, errors: []string{}}

	parser.currentToken = parser.nextToken
	parser.nextToken = lxr.NextToken()

	return parser
}

// jump moves the current and next token to the corresponding next token in the lexer
func (p *Parser) jump() {
	p.currentToken = p.nextToken
	p.nextToken = p.lxr.NextToken()
}

// ParseProgram performs recursive descent parsing (aka Pratt parsing)
func (p *Parser) ParseProgram() *syntaxtree.Program {
	program := &syntaxtree.Program{Statements: []syntaxtree.Stmt{}}

	// iterate over all the tokens in the lexer until EOF is hit
	for p.currentToken.Type != token.EOF {
		// for each token double (cur, nxt) call parseStatement
		statement := p.parseStatement()

		// if the resulting statement is not nil, add it to the program
		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}

		// move to the next token in the lexer
		p.jump()
	}

	return program
}

// parseStatement handles parsing statements
func (p *Parser) parseStatement() syntaxtree.Stmt {
	switch p.currentToken.Type {
	case token.Let:
		return p.parseLetStatement()
	default:
		return nil
	}
}

// parseLetStatement takes care of parsing let statements
func (p *Parser) parseLetStatement() *syntaxtree.LetStmt {
	stmt := &syntaxtree.LetStmt{Token: p.currentToken}

	// if the next token is not an identifier, this is an invalid let statement
	if p.nextToken.Type != token.Identifier {
		msg := generateErrorMsg(p.currentToken.Type, token.Identifier, p.nextToken.Type)
		p.errors = append(p.errors, msg)
		return nil
	}

	// lets move the identifier as current token
	p.jump()

	stmt.Name = &syntaxtree.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}

	// currently, we have let identifier
	// if the next token is not an equal assign, this is an invalid let statement
	if p.nextToken.Type != token.Assign {

		return nil
	}

	// TODO: handle parsing expression to statement value
	// This is not that easy currently, since we are planning to have functions,
	// strings, integers, maps, arrays,
	// skipping the expression until we get to the semicolon
	for p.currentToken.Type != token.SemiColon {
		p.jump()
	}

	return stmt
}

// TODO: logging
// generateErrorMsg generated error message for unexpected token retrieval
func generateErrorMsg(cur token.Type, exp token.Type, actual token.Type) string {
	return fmt.Sprintf("current token (%s) expected next token to be (%s), but got (%s)",
		cur, exp, actual)
}

// GetErrors returns the encountered errors of the parser
func (p *Parser) GetErrors() []string {
	return p.errors
}
