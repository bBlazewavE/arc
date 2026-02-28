package parser

import "github.com/bBlazewavE/arc/lexer"

// Node is the base interface for all AST nodes
type Node interface {
	TokenLiteral() string
}

// Statement nodes
type Statement interface {
	Node
	statementNode()
}

// Expression nodes
type Expression interface {
	Node
	expressionNode()
}

// Program is the root node — a list of statements
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

// LetStatement: let x: int = 10
type LetStatement struct {
	Token    lexer.Token // the LET token
	Name     string
	TypeName string // "int", "string", "bool"
	Value    Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

// ReturnStatement: return expr
type ReturnStatement struct {
	Token lexer.Token
	Value Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

// PrintStatement: print(expr)
type PrintStatement struct {
	Token lexer.Token
	Value Expression
}

func (ps *PrintStatement) statementNode()       {}
func (ps *PrintStatement) TokenLiteral() string { return ps.Token.Literal }

// ExpressionStatement wraps an expression as a statement
type ExpressionStatement struct {
	Token      lexer.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

// IfStatement: if cond { ... } else { ... }
type IfStatement struct {
	Token       lexer.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement // nil if no else
}

func (is *IfStatement) statementNode()       {}
func (is *IfStatement) TokenLiteral() string { return is.Token.Literal }

// BlockStatement: { stmt1; stmt2; ... }
type BlockStatement struct {
	Token      lexer.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }

// FnStatement: fn add(a: int, b: int) -> int { ... }
type FnStatement struct {
	Token      lexer.Token
	Name       string
	Params     []Param
	ReturnType string // "" if no return type
	Body       *BlockStatement
}

type Param struct {
	Name     string
	TypeName string
}

func (fs *FnStatement) statementNode()       {}
func (fs *FnStatement) TokenLiteral() string { return fs.Token.Literal }

// --- Expressions ---

// IntegerLiteral: 42
type IntegerLiteral struct {
	Token lexer.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }

// StringLiteral: "hello"
type StringLiteral struct {
	Token lexer.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }

// BoolLiteral: true / false
type BoolLiteral struct {
	Token lexer.Token
	Value bool
}

func (bl *BoolLiteral) expressionNode()      {}
func (bl *BoolLiteral) TokenLiteral() string { return bl.Token.Literal }

// Identifier: x, name, myVar
type Identifier struct {
	Token lexer.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

// BinaryExpression: left OP right (e.g., x + 5)
type BinaryExpression struct {
	Token    lexer.Token // the operator token
	Left     Expression
	Operator string
	Right    Expression
}

func (be *BinaryExpression) expressionNode()      {}
func (be *BinaryExpression) TokenLiteral() string { return be.Token.Literal }

// UnaryExpression: !expr, -expr
type UnaryExpression struct {
	Token    lexer.Token
	Operator string
	Right    Expression
}

func (ue *UnaryExpression) expressionNode()      {}
func (ue *UnaryExpression) TokenLiteral() string { return ue.Token.Literal }

// CallExpression: add(1, 2)
type CallExpression struct {
	Token     lexer.Token
	Function  string
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
