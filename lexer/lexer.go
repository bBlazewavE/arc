package lexer

import "fmt"

// Lexer turns source code into tokens
type Lexer struct {
	input   string
	pos     int  // current position (points to current char)
	readPos int  // next position (after current char)
	ch      byte // current character
	line    int
	column  int
}

// New creates a new Lexer
func New(input string) *Lexer {
	l := &Lexer{input: input, line: 1, column: 0}
	l.readChar() // initialize first character
	return l
}

// readChar reads the next character and advances position
func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.ch = 0 // ASCII NUL = end of input
	} else {
		l.ch = l.input[l.readPos]
	}
	l.pos = l.readPos
	l.readPos++
	l.column++
	if l.ch == '\n' {
		l.line++
		l.column = 0
	}
}

// peekChar looks at the next character without advancing
func (l *Lexer) peekChar() byte {
	if l.readPos >= len(l.input) {
		return 0
	}
	return l.input[l.readPos]
}

// NextToken returns the next token from the input
func (l *Lexer) NextToken() Token {
	l.skipWhitespace()
	l.skipComments()

	tok := Token{Line: l.line, Column: l.column}

	switch l.ch {
	case '+':
		tok = l.newToken(PLUS, "+")
	case '*':
		tok = l.newToken(STAR, "*")
	case '/':
		tok = l.newToken(SLASH, "/")
	case '(':
		tok = l.newToken(LPAREN, "(")
	case ')':
		tok = l.newToken(RPAREN, ")")
	case '{':
		tok = l.newToken(LBRACE, "{")
	case '}':
		tok = l.newToken(RBRACE, "}")
	case ',':
		tok = l.newToken(COMMA, ",")
	case ':':
		tok = l.newToken(COLON, ":")

	case '-':
		if l.peekChar() == '>' {
			l.readChar()
			tok = l.newToken(ARROW, "->")
		} else {
			tok = l.newToken(MINUS, "-")
		}

	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok = l.newToken(EQ, "==")
		} else {
			tok = l.newToken(ASSIGN, "=")
		}

	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = l.newToken(NEQ, "!=")
		} else {
			tok = l.newToken(NOT, "!")
		}

	case '<':
		if l.peekChar() == '=' {
			l.readChar()
			tok = l.newToken(LTE, "<=")
		} else {
			tok = l.newToken(LT, "<")
		}

	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			tok = l.newToken(GTE, ">=")
		} else {
			tok = l.newToken(GT, ">")
		}

	case '&':
		if l.peekChar() == '&' {
			l.readChar()
			tok = l.newToken(AND, "&&")
		} else {
			tok = l.newToken(ILLEGAL, string(l.ch))
		}

	case '|':
		if l.peekChar() == '|' {
			l.readChar()
			tok = l.newToken(OR, "||")
		} else {
			tok = l.newToken(ILLEGAL, string(l.ch))
		}

	case '"':
		tok.Type = STRING
		tok.Literal = l.readString()
		tok.Line = l.line
		tok.Column = l.column
		return tok // don't readChar again, readString already advanced

	case 0:
		tok = l.newToken(EOF, "")
		return tok

	default:
		if isLetter(l.ch) {
			literal := l.readIdentifier()
			tokenType := LookupIdent(literal)
			return Token{Type: tokenType, Literal: literal, Line: l.line, Column: l.column}
		} else if isDigit(l.ch) {
			literal := l.readNumber()
			return Token{Type: INT, Literal: literal, Line: l.line, Column: l.column}
		} else {
			tok = l.newToken(ILLEGAL, string(l.ch))
		}
	}

	l.readChar()
	return tok
}

// Tokenize returns all tokens from the input
func (l *Lexer) Tokenize() []Token {
	var tokens []Token
	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == EOF {
			break
		}
	}
	return tokens
}

func (l *Lexer) newToken(tokenType TokenType, literal string) Token {
	return Token{Type: tokenType, Literal: literal, Line: l.line, Column: l.column}
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) skipComments() {
	if l.ch == '/' && l.peekChar() == '/' {
		for l.ch != '\n' && l.ch != 0 {
			l.readChar()
		}
		l.skipWhitespace()
		l.skipComments() // handle multiple comment lines
	}
}

func (l *Lexer) readIdentifier() string {
	start := l.pos
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}
	return l.input[start:l.pos]
}

func (l *Lexer) readNumber() string {
	start := l.pos
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[start:l.pos]
}

func (l *Lexer) readString() string {
	l.readChar() // skip opening quote
	start := l.pos
	for l.ch != '"' && l.ch != 0 {
		l.readChar()
	}
	str := l.input[start:l.pos]
	if l.ch == '"' {
		l.readChar() // skip closing quote
	} else {
		fmt.Printf("Error: unterminated string at line %d\n", l.line)
	}
	return str
}

func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}
