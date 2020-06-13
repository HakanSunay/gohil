package parser

import (
	"fmt"
	"strconv"

	"github.com/HakanSunay/gohil/lexer"
	"github.com/HakanSunay/gohil/syntaxtree"
	"github.com/HakanSunay/gohil/token"
)

const (
	Lowest      = iota + 1
	Equals      // ==
	LessGreater // > or <
	Sum         // +
	Product     // *
	Prefix      // -X or !X
	Call        // myFunction(X)
)

// parserFunc types are used for Pratt parsing
// only the infixParseFN takes an argument,
// because that is the left side of the operator that is being parsed
// E.g: 6 + _ (+ is the operator and 6 is the argument)
// prefixParseFN doesn't have a left side for its operator
type (
	prefixParseFN func() syntaxtree.Expr
	infixParseFN  func(syntaxtree.Expr) syntaxtree.Expr
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

	prefixMap map[token.Type]prefixParseFN
	infixMap  map[token.Type]infixParseFN

	errors []string
}

// NewParser is the constructor for the Parser type
func NewParser(lxr *lexer.Lexer) *Parser {
	parser := &Parser{
		lxr: lxr,

		prefixMap: make(map[token.Type]prefixParseFN),
		infixMap:  make(map[token.Type]infixParseFN),

		errors: []string{},
	}

	parser.jump()
	parser.jump()

	// register the available parsing functions
	parser.addPrefixFunc(token.Identifier, parser.parseIdentifier)
	parser.addPrefixFunc(token.Int, parser.parseIntegerLiteral)

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
	case token.Return:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
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
		msg := generateErrorMsg(p.currentToken.Type, token.Assign, p.nextToken.Type)
		p.errors = append(p.errors, msg)
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

func (p *Parser) parseReturnStatement() *syntaxtree.ReturnStmt {
	stmt := &syntaxtree.ReturnStmt{Token: p.currentToken}

	// let's move to the next token - the expression
	p.jump()

	// TODO: handle expression parsing to ReturnValue

	for p.currentToken.Type != token.SemiColon {
		p.jump()
	}

	return stmt
}

func (p *Parser) addPrefixFunc(tokenType token.Type, fn prefixParseFN) {
	p.prefixMap[tokenType] = fn
}

func (p *Parser) addInfixFunc(tokenType token.Type, fn infixParseFN) {
	p.infixMap[tokenType] = fn
}

func (p *Parser) parseExpressionStatement() *syntaxtree.ExpressionStmt {
	stmt := &syntaxtree.ExpressionStmt{Token: p.currentToken}

	stmt.Expression = p.parseExpression(Lowest)

	// 6;
	// 6
	// are both valid in gohil
	if p.nextToken.Type == token.SemiColon {
		p.jump()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) syntaxtree.Expr {
	// is there a parsing function that can handle the current token type
	prefix, ok := p.prefixMap[p.currentToken.Type]
	if !ok {
		msg := fmt.Sprintf("no prefix parse function for (%s) found",
			p.currentToken.Type)
		p.errors = append(p.errors, msg)
		return nil
	}
	leftExpr := prefix()
	return leftExpr
}

func (p *Parser) parseIdentifier() syntaxtree.Expr {
	return &syntaxtree.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
}

func (p *Parser) parseIntegerLiteral() syntaxtree.Expr {
	integerLiteral := &syntaxtree.IntegerLiteral{Token: p.currentToken}
	value, err := strconv.Atoi(p.currentToken.Literal)
	if err != nil {
		msg := fmt.Sprintf("could not parse (%s) to integer", p.currentToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	integerLiteral.Value = value
	return integerLiteral
}
