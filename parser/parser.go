package parser

import (
	"fmt"
	"strconv"

	"github.com/HakanSunay/gohil/lexer"
	"github.com/HakanSunay/gohil/syntaxtree"
	"github.com/HakanSunay/gohil/token"
)

const (
	Lowest = iota + 1
	Equals
	LessGreater
	Sum
	Product
	Prefix
	Call
)

var precedences = map[token.Type]int{
	token.Equal:    Equals,
	token.NotEqual: Equals,

	token.LessThan:    LessGreater,
	token.GreaterThan: LessGreater,

	token.Plus:  Sum,
	token.Minus: Sum,

	token.Slash:    Product,
	token.Asterisk: Product,

	token.Function:        Call,
	token.LeftParenthesis: Call,
}

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

	// prefix funcs
	parser.addPrefixFunc(token.Identifier, parser.parseIdentifier)
	parser.addPrefixFunc(token.Int, parser.parseIntegerLiteral)
	parser.addPrefixFunc(token.True, parser.parseBooleanLiteral)
	parser.addPrefixFunc(token.False, parser.parseBooleanLiteral)
	parser.addPrefixFunc(token.ExclamationMark, parser.parsePrefixExpression)
	parser.addPrefixFunc(token.Minus, parser.parsePrefixExpression)
	parser.addPrefixFunc(token.LeftParenthesis, parser.parseGroupedExpression)
	parser.addPrefixFunc(token.If, parser.parseIfExpression)
	parser.addPrefixFunc(token.Function, parser.parseFunctionLiteral)
	parser.addPrefixFunc(token.String, parser.parseStringLiteral)
	parser.addPrefixFunc(token.LeftBracket, parser.parseArrayLiteral)

	// infix funcs
	parser.addInfixFunc(token.Plus, parser.parseInfixExpression)
	parser.addInfixFunc(token.Minus, parser.parseInfixExpression)
	parser.addInfixFunc(token.Slash, parser.parseInfixExpression)
	parser.addInfixFunc(token.Asterisk, parser.parseInfixExpression)
	parser.addInfixFunc(token.Equal, parser.parseInfixExpression)
	parser.addInfixFunc(token.NotEqual, parser.parseInfixExpression)
	parser.addInfixFunc(token.LessThan, parser.parseInfixExpression)
	parser.addInfixFunc(token.GreaterThan, parser.parseInfixExpression)
	parser.addInfixFunc(token.LeftParenthesis, parser.parseCallExpression)

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
	// jump to to the assign token
	p.jump()

	// jump to the value
	p.jump()
	stmt.Value = p.parseExpression(Lowest)

	if p.nextToken.Type == token.SemiColon {
		p.jump()
	}

	return stmt
}

// TODO: logging
// generateErrorMsg generated error message for unexpected token retrieval
func generateErrorMsg(cur token.Type, exp token.Type, actual token.Type) string {
	return fmt.Sprintf("Current token of type (%s) expected next token of type (%s), but got (%s)",
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

	// parse the return value
	stmt.ReturnValue = p.parseExpression(Lowest)

	if p.nextToken.Type == token.SemiColon {
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

	for p.nextToken.Type != token.SemiColon && precedence < p.getNextPrecedence() {
		infix, ok := p.infixMap[p.nextToken.Type]
		if !ok {
			msg := fmt.Sprintf("no infix parse function for (%s) found",
				p.nextToken.Type)
			p.errors = append(p.errors, msg)
			return leftExpr
		}

		p.jump()

		// the infix takes the left expr as a parameter and updates it
		leftExpr = infix(leftExpr)
	}

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

func (p *Parser) parseBooleanLiteral() syntaxtree.Expr {
	return &syntaxtree.BooleanLiteral{
		Token: p.currentToken,
		// if the token type is True -> assigning to true; else False
		Value: p.currentToken.Type == token.True,
	}
}

func (p *Parser) parsePrefixExpression() syntaxtree.Expr {
	// imagine getting !66 as parameter
	// ! becomes the current expression and its token is !
	expression := &syntaxtree.PrefixExpr{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
	}

	// moving to the next token
	p.jump()

	// now assigning 66 integerLiteral to the right side
	// of the current expression recursively
	// second highest operator precedence, after func calls
	expression.Right = p.parseExpression(Prefix)

	return expression
}

func (p *Parser) getCurrentPrecedence() int {
	val, ok := precedences[p.currentToken.Type]
	if !ok {
		return Lowest
	}

	return val
}

func (p *Parser) getNextPrecedence() int {
	val, ok := precedences[p.nextToken.Type]
	if !ok {
		return Lowest
	}

	return val
}

func (p *Parser) parseInfixExpression(leftExpr syntaxtree.Expr) syntaxtree.Expr {
	expr := &syntaxtree.InfixExpr{
		Token:    p.currentToken,
		Left:     leftExpr,
		Operator: p.currentToken.Literal,
	}

	precedence := p.getCurrentPrecedence()
	p.jump()
	expr.Right = p.parseExpression(precedence)

	return expr
}

func (p *Parser) parseGroupedExpression() syntaxtree.Expr {
	p.jump()

	expr := p.parseExpression(Lowest)

	if p.nextToken.Type != token.RightParenthesis {
		msg := generateErrorMsg(p.currentToken.Type, token.RightParenthesis, p.nextToken.Type)
		p.errors = append(p.errors, msg)
		return nil
	}

	p.jump()
	return expr
}

func (p *Parser) parseIfExpression() syntaxtree.Expr {
	ifExpr := &syntaxtree.IfExpr{Token: p.currentToken}

	if p.nextToken.Type != token.LeftParenthesis {
		msg := generateErrorMsg(p.currentToken.Type, token.LeftParenthesis, p.nextToken.Type)
		p.errors = append(p.errors, msg)
		return nil
	}
	// jump to the left parenthesis
	p.jump()

	// jump to the condition
	p.jump()
	// parse the condition
	ifExpr.Condition = p.parseExpression(Lowest)

	if p.nextToken.Type != token.RightParenthesis {
		msg := generateErrorMsg(p.currentToken.Type, token.RightParenthesis, p.nextToken.Type)
		p.errors = append(p.errors, msg)
		return nil
	}
	p.jump()

	// parse the block statement leading (
	if p.nextToken.Type != token.LeftBrace {
		msg := generateErrorMsg(p.currentToken.Type, token.LeftBrace, p.nextToken.Type)
		p.errors = append(p.errors, msg)
		return nil
	}
	p.jump()

	// parse the block statement
	ifExpr.Consequence = p.parseBlockStatement()

	// if the ELSE is present, parse it and its block statement
	if p.nextToken.Type == token.Else {
		p.jump()

		if p.nextToken.Type != token.LeftBrace {
			msg := generateErrorMsg(p.currentToken.Type, token.LeftBrace, p.nextToken.Type)
			p.errors = append(p.errors, msg)
			return nil
		}
		p.jump()

		ifExpr.Alternative = p.parseBlockStatement()
	}

	return ifExpr
}

func (p *Parser) parseBlockStatement() *syntaxtree.BlockStmt {
	blockStmt := &syntaxtree.BlockStmt{
		Token:      p.currentToken,
		Statements: []syntaxtree.Stmt{},
	}

	p.jump()
	for p.currentToken.Type != token.RightBrace {
		stmt := p.parseStatement()
		if stmt != nil {
			blockStmt.Statements = append(blockStmt.Statements, stmt)
		}
		p.jump()
	}

	return blockStmt
}

func (p *Parser) parseFunctionLiteral() syntaxtree.Expr {
	fnLiteral := &syntaxtree.FunctionLiteral{Token: p.currentToken}

	if p.nextToken.Type != token.LeftParenthesis {
		msg := generateErrorMsg(p.currentToken.Type, token.LeftParenthesis, p.nextToken.Type)
		p.errors = append(p.errors, msg)
		return nil
	}
	// jump to the left parenthesis
	p.jump()

	fnLiteral.Parameters = p.parseFunctionParameters()

	// after parsing the parameters, the next token must be the left brace.
	// opening the function body
	if p.nextToken.Type != token.LeftBrace {
		msg := generateErrorMsg(p.currentToken.Type, token.LeftBrace, p.nextToken.Type)
		p.errors = append(p.errors, msg)
		return nil
	}
	p.jump()

	// same as if expr consequence / alternative parsing
	fnLiteral.Body = p.parseBlockStatement()

	return fnLiteral
}

func (p *Parser) parseFunctionParameters() []*syntaxtree.Identifier {
	var ids []*syntaxtree.Identifier

	// no parameters
	if p.nextToken.Type == token.RightParenthesis {
		p.jump()
		return ids
	}

	p.jump()

	// since there is right parenthesis, there is at least 1 parameters
	// therefore we parse it manually
	id := &syntaxtree.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
	ids = append(ids, id)

	// while there are parameters left, add them to the ids
	for p.nextToken.Type == token.Comma {
		// jump to the comma
		p.jump()
		// jump to the parameter
		p.jump()

		id := &syntaxtree.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
		ids = append(ids, id)
	}

	// no more comma, the next token must be a right parenthesis
	if p.nextToken.Type != token.RightParenthesis {
		msg := generateErrorMsg(p.currentToken.Type, token.RightParenthesis, p.nextToken.Type)
		p.errors = append(p.errors, msg)
		return nil
	}
	// jump to the right parenthesis
	p.jump()

	return ids
}

func (p *Parser) parseCallExpression(fn syntaxtree.Expr) syntaxtree.Expr {
	exp := &syntaxtree.CallExpr{Token: p.currentToken, Function: fn}
	exp.Arguments = p.parseCallArguments()

	return exp
}

func (p *Parser) parseCallArguments() []syntaxtree.Expr {
	var args []syntaxtree.Expr

	// no args
	if p.nextToken.Type == token.RightParenthesis {
		p.jump()
		return args
	}

	// jump to first arg
	p.jump()

	// since every argument is an expression, lets parse the first one
	args = append(args, p.parseExpression(Lowest))

	// similar to parsing function parameters
	// while there is an argument left, keep on parsing them
	for p.nextToken.Type == token.Comma {
		// jump to the comma
		p.jump()
		// jump to the arg
		p.jump()
		args = append(args, p.parseExpression(Lowest))
	}

	// no more comma, the next token must be a right parenthesis
	if p.nextToken.Type != token.RightParenthesis {
		msg := generateErrorMsg(p.currentToken.Type, token.RightParenthesis, p.nextToken.Type)
		p.errors = append(p.errors, msg)
		return nil
	}
	// jump to the right parenthesis
	p.jump()

	return args
}

func (p *Parser) parseStringLiteral() syntaxtree.Expr {
	return &syntaxtree.StringLiteral{Token: p.currentToken, Value: p.currentToken.Literal}
}

func (p *Parser) parseArrayLiteral() syntaxtree.Expr {
	array := &syntaxtree.ArrayLiteral{Token: p.currentToken}
	array.Elements = p.parseExpressionList(token.RightBracket)
	return array
}

// TODO: merge this with parseCallArguments
func (p *Parser) parseExpressionList(bracket token.Type) []syntaxtree.Expr {
	var list []syntaxtree.Expr

	// no elements
	if p.nextToken.Type == bracket {
		p.jump()
		return list
	}

	// jump to first element
	p.jump()

	// since every element is an expression, lets parse the first one
	list = append(list, p.parseExpression(Lowest))

	// similar to parsing function parameters
	// while there is an element left, keep on parsing them
	for p.nextToken.Type == token.Comma {
		// jump to the comma
		p.jump()
		// jump to the element
		p.jump()
		list = append(list, p.parseExpression(Lowest))
	}

	// no more comma, the next token must be a right parenthesis
	if p.nextToken.Type != bracket {
		msg := generateErrorMsg(p.currentToken.Type, bracket, p.nextToken.Type)
		p.errors = append(p.errors, msg)
		return nil
	}

	// jump to the right parenthesis
	p.jump()

	return list
}
