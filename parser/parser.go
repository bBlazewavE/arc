package parser

import (
	"fmt"
	"strconv"

	"github.com/bBlazewavE/arc/lexer"
)

// Operator precedence levels
const (
	_ int = iota
	LOWEST
	OR_PREC      // ||
	AND_PREC     // &&
	EQUALS       // == !=
	LESSGREATER  // > < >= <=
	SUM          // + -
	PRODUCT      // * /
	PREFIX       // -x !x
	CALL         // fn(x)
)

var precedences = map[lexer.TokenType]int{
	lexer.OR:     OR_PREC,
	lexer.AND:    AND_PREC,
	lexer.EQ:     EQUALS,
	lexer.NEQ:    EQUALS,
	lexer.LT:     LESSGREATER,
	lexer.GT:     LESSGREATER,
	lexer.LTE:    LESSGREATER,
	lexer.GTE:    LESSGREATER,
	lexer.PLUS:   SUM,
	lexer.MINUS:  SUM,
	lexer.STAR:   PRODUCT,
	lexer.SLASH:  PRODUCT,
	lexer.LPAREN: CALL,
}

// Parser turns tokens into an AST
type Parser struct {
	tokens  []lexer.Token
	pos     int
	errors  []string
}

// New creates a new Parser
func New(tokens []lexer.Token) *Parser {
	return &Parser{tokens: tokens, pos: 0}
}

// Errors returns parsing errors
func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) current() lexer.Token {
	if p.pos >= len(p.tokens) {
		return lexer.Token{Type: lexer.EOF}
	}
	return p.tokens[p.pos]
}

func (p *Parser) peek() lexer.Token {
	if p.pos+1 >= len(p.tokens) {
		return lexer.Token{Type: lexer.EOF}
	}
	return p.tokens[p.pos+1]
}

func (p *Parser) advance() lexer.Token {
	tok := p.current()
	p.pos++
	return tok
}

func (p *Parser) expect(t lexer.TokenType) lexer.Token {
	tok := p.current()
	if tok.Type != t {
		p.errors = append(p.errors, fmt.Sprintf(
			"line %d: expected %s, got %s (%q)",
			tok.Line, t, tok.Type, tok.Literal,
		))
	}
	p.advance()
	return tok
}

func (p *Parser) currentPrecedence() int {
	if prec, ok := precedences[p.current().Type]; ok {
		return prec
	}
	return LOWEST
}

// Parse parses the token stream into a Program AST
func (p *Parser) Parse() *Program {
	program := &Program{}
	for p.current().Type != lexer.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
	}
	return program
}

func (p *Parser) parseStatement() Statement {
	switch p.current().Type {
	case lexer.LET:
		return p.parseLetStatement()
	case lexer.FN:
		return p.parseFnStatement()
	case lexer.IF:
		return p.parseIfStatement()
	case lexer.RETURN:
		return p.parseReturnStatement()
	case lexer.PRINT:
		return p.parsePrintStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// let x: int = 10
func (p *Parser) parseLetStatement() *LetStatement {
	stmt := &LetStatement{Token: p.advance()} // consume LET

	nameToken := p.expect(lexer.IDENT)
	stmt.Name = nameToken.Literal

	p.expect(lexer.COLON) // consume :

	stmt.TypeName = p.parseTypeName()

	p.expect(lexer.ASSIGN) // consume =

	stmt.Value = p.parseExpression(LOWEST)

	return stmt
}

// fn add(a: int, b: int) -> int { ... }
func (p *Parser) parseFnStatement() *FnStatement {
	stmt := &FnStatement{Token: p.advance()} // consume FN

	nameToken := p.expect(lexer.IDENT)
	stmt.Name = nameToken.Literal

	p.expect(lexer.LPAREN)
	stmt.Params = p.parseParams()
	p.expect(lexer.RPAREN)

	// Optional return type: -> int
	if p.current().Type == lexer.ARROW {
		p.advance() // consume ->
		stmt.ReturnType = p.parseTypeName()
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseParams() []Param {
	var params []Param

	if p.current().Type == lexer.RPAREN {
		return params
	}

	for {
		name := p.expect(lexer.IDENT).Literal
		p.expect(lexer.COLON)
		typeName := p.parseTypeName()
		params = append(params, Param{Name: name, TypeName: typeName})

		if p.current().Type != lexer.COMMA {
			break
		}
		p.advance() // consume comma
	}

	return params
}

func (p *Parser) parseTypeName() string {
	tok := p.current()
	switch tok.Type {
	case lexer.TYPE_INT:
		p.advance()
		return "int"
	case lexer.TYPE_STRING:
		p.advance()
		return "string"
	case lexer.TYPE_BOOL:
		p.advance()
		return "bool"
	default:
		p.errors = append(p.errors, fmt.Sprintf(
			"line %d: expected type (int, string, bool), got %s (%q)",
			tok.Line, tok.Type, tok.Literal,
		))
		p.advance()
		return "unknown"
	}
}

// if cond { ... } else { ... }
func (p *Parser) parseIfStatement() *IfStatement {
	stmt := &IfStatement{Token: p.advance()} // consume IF

	stmt.Condition = p.parseExpression(LOWEST)
	stmt.Consequence = p.parseBlockStatement()

	if p.current().Type == lexer.ELSE {
		p.advance() // consume ELSE
		stmt.Alternative = p.parseBlockStatement()
	}

	return stmt
}

// return expr
func (p *Parser) parseReturnStatement() *ReturnStatement {
	stmt := &ReturnStatement{Token: p.advance()} // consume RETURN
	stmt.Value = p.parseExpression(LOWEST)
	return stmt
}

// print(expr)
func (p *Parser) parsePrintStatement() *PrintStatement {
	stmt := &PrintStatement{Token: p.advance()} // consume PRINT
	p.expect(lexer.LPAREN)
	stmt.Value = p.parseExpression(LOWEST)
	p.expect(lexer.RPAREN)
	return stmt
}

func (p *Parser) parseExpressionStatement() *ExpressionStatement {
	stmt := &ExpressionStatement{Token: p.current()}
	stmt.Expression = p.parseExpression(LOWEST)
	return stmt
}

// { stmt1; stmt2; ... }
func (p *Parser) parseBlockStatement() *BlockStatement {
	block := &BlockStatement{Token: p.current()}
	p.expect(lexer.LBRACE)

	for p.current().Type != lexer.RBRACE && p.current().Type != lexer.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
	}

	p.expect(lexer.RBRACE)
	return block
}

// Pratt parser for expressions
func (p *Parser) parseExpression(precedence int) Expression {
	left := p.parsePrefixExpression()

	for precedence < p.currentPrecedence() {
		if p.current().Type == lexer.LPAREN {
			// Function call
			if ident, ok := left.(*Identifier); ok {
				left = p.parseCallExpression(ident.Value)
			} else {
				break
			}
		} else {
			// Binary operator
			left = p.parseBinaryExpression(left)
		}
	}

	return left
}

func (p *Parser) parsePrefixExpression() Expression {
	tok := p.current()

	switch tok.Type {
	case lexer.INT:
		p.advance()
		val, err := strconv.ParseInt(tok.Literal, 10, 64)
		if err != nil {
			p.errors = append(p.errors, fmt.Sprintf("line %d: invalid integer: %s", tok.Line, tok.Literal))
			return nil
		}
		return &IntegerLiteral{Token: tok, Value: val}

	case lexer.STRING:
		p.advance()
		return &StringLiteral{Token: tok, Value: tok.Literal}

	case lexer.TRUE:
		p.advance()
		return &BoolLiteral{Token: tok, Value: true}

	case lexer.FALSE:
		p.advance()
		return &BoolLiteral{Token: tok, Value: false}

	case lexer.IDENT:
		p.advance()
		return &Identifier{Token: tok, Value: tok.Literal}

	case lexer.LPAREN:
		p.advance() // consume (
		expr := p.parseExpression(LOWEST)
		p.expect(lexer.RPAREN) // consume )
		return expr

	case lexer.MINUS, lexer.NOT:
		p.advance()
		return &UnaryExpression{
			Token:    tok,
			Operator: tok.Literal,
			Right:    p.parseExpression(PREFIX),
		}

	default:
		p.errors = append(p.errors, fmt.Sprintf(
			"line %d: unexpected token: %s (%q)", tok.Line, tok.Type, tok.Literal,
		))
		p.advance()
		return nil
	}
}

func (p *Parser) parseBinaryExpression(left Expression) Expression {
	tok := p.advance() // consume operator
	prec := precedences[tok.Type]
	right := p.parseExpression(prec)
	return &BinaryExpression{
		Token:    tok,
		Left:     left,
		Operator: tok.Literal,
		Right:    right,
	}
}

func (p *Parser) parseCallExpression(fnName string) *CallExpression {
	tok := p.current()
	p.expect(lexer.LPAREN)

	var args []Expression
	if p.current().Type != lexer.RPAREN {
		args = append(args, p.parseExpression(LOWEST))
		for p.current().Type == lexer.COMMA {
			p.advance()
			args = append(args, p.parseExpression(LOWEST))
		}
	}
	p.expect(lexer.RPAREN)

	return &CallExpression{Token: tok, Function: fnName, Arguments: args}
}
