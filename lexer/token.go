package lexer

// TokenType represents the type of a token
type TokenType string

const (
	// Literals
	INT    TokenType = "INT"    // 42
	STRING TokenType = "STRING" // "hello"
	IDENT  TokenType = "IDENT"  // variable names

	// Keywords
	LET    TokenType = "LET"    // let
	FN     TokenType = "FN"     // fn
	IF     TokenType = "IF"     // if
	ELSE   TokenType = "ELSE"   // else
	RETURN TokenType = "RETURN" // return
	PRINT  TokenType = "PRINT"  // print
	TRUE   TokenType = "TRUE"   // true
	FALSE  TokenType = "FALSE"  // false

	// Types
	TYPE_INT    TokenType = "TYPE_INT"    // int
	TYPE_STRING TokenType = "TYPE_STRING" // string
	TYPE_BOOL   TokenType = "TYPE_BOOL"   // bool

	// Operators
	PLUS     TokenType = "PLUS"     // +
	MINUS    TokenType = "MINUS"    // -
	STAR     TokenType = "STAR"     // *
	SLASH    TokenType = "SLASH"    // /
	ASSIGN   TokenType = "ASSIGN"   // =
	EQ       TokenType = "EQ"       // ==
	NEQ      TokenType = "NEQ"      // !=
	LT       TokenType = "LT"       // <
	GT       TokenType = "GT"       // >
	LTE      TokenType = "LTE"      // <=
	GTE      TokenType = "GTE"      // >=
	AND      TokenType = "AND"      // &&
	OR       TokenType = "OR"       // ||
	NOT      TokenType = "NOT"      // !

	// Delimiters
	COLON    TokenType = "COLON"    // :
	ARROW    TokenType = "ARROW"    // ->
	LPAREN   TokenType = "LPAREN"   // (
	RPAREN   TokenType = "RPAREN"   // )
	LBRACE   TokenType = "LBRACE"   // {
	RBRACE   TokenType = "RBRACE"   // }
	COMMA    TokenType = "COMMA"    // ,

	// Special
	EOF     TokenType = "EOF"
	ILLEGAL TokenType = "ILLEGAL"
)

// Token represents a single token from the source code
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

// keywords maps keyword strings to their token types
var keywords = map[string]TokenType{
	"let":    LET,
	"fn":     FN,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"print":  PRINT,
	"true":   TRUE,
	"false":  FALSE,
	"int":    TYPE_INT,
	"string": TYPE_STRING,
	"bool":   TYPE_BOOL,
}

// LookupIdent checks if an identifier is a keyword
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
