package typechecker

import (
	"fmt"

	"github.com/bBlazewavE/arc/parser"
)

// Type represents an Arc type
type Type string

const (
	TypeInt     Type = "int"
	TypeString  Type = "string"
	TypeBool    Type = "bool"
	TypeVoid    Type = "void"
	TypeUnknown Type = "unknown"
)

// FuncType stores function signature info
type FuncType struct {
	Params     []Type
	ReturnType Type
}

// Checker performs static type checking on the AST
type Checker struct {
	variables map[string]Type
	functions map[string]FuncType
	errors    []string
}

// New creates a new type checker
func New() *Checker {
	return &Checker{
		variables: make(map[string]Type),
		functions: make(map[string]FuncType),
	}
}

// Errors returns type checking errors
func (c *Checker) Errors() []string {
	return c.errors
}

// Check performs type checking on a program
func (c *Checker) Check(program *parser.Program) {
	for _, stmt := range program.Statements {
		c.checkStatement(stmt)
	}
}

func (c *Checker) checkStatement(stmt parser.Statement) {
	switch s := stmt.(type) {
	case *parser.LetStatement:
		c.checkLetStatement(s)
	case *parser.FnStatement:
		c.checkFnStatement(s)
	case *parser.IfStatement:
		c.checkIfStatement(s)
	case *parser.ReturnStatement:
		c.checkExpression(s.Value)
	case *parser.PrintStatement:
		c.checkExpression(s.Value)
	case *parser.ExpressionStatement:
		c.checkExpression(s.Expression)
	}
}

func (c *Checker) checkLetStatement(s *parser.LetStatement) {
	declaredType := Type(s.TypeName)
	actualType := c.checkExpression(s.Value)

	if actualType != TypeUnknown && actualType != declaredType {
		c.errors = append(c.errors, fmt.Sprintf(
			"type mismatch: variable %q declared as %s but assigned %s",
			s.Name, declaredType, actualType,
		))
	}

	c.variables[s.Name] = declaredType
}

func (c *Checker) checkFnStatement(s *parser.FnStatement) {
	// Register function signature
	var paramTypes []Type
	for _, p := range s.Params {
		paramTypes = append(paramTypes, Type(p.TypeName))
	}

	returnType := TypeVoid
	if s.ReturnType != "" {
		returnType = Type(s.ReturnType)
	}

	c.functions[s.Name] = FuncType{
		Params:     paramTypes,
		ReturnType: returnType,
	}

	// Create a new scope for function body
	prevVars := make(map[string]Type)
	for k, v := range c.variables {
		prevVars[k] = v
	}

	// Add params to scope
	for _, p := range s.Params {
		c.variables[p.Name] = Type(p.TypeName)
	}

	// Check body
	for _, stmt := range s.Body.Statements {
		c.checkStatement(stmt)
	}

	// Restore outer scope
	c.variables = prevVars
}

func (c *Checker) checkIfStatement(s *parser.IfStatement) {
	condType := c.checkExpression(s.Condition)
	if condType != TypeBool && condType != TypeUnknown {
		c.errors = append(c.errors, fmt.Sprintf(
			"if condition must be bool, got %s", condType,
		))
	}

	for _, stmt := range s.Consequence.Statements {
		c.checkStatement(stmt)
	}

	if s.Alternative != nil {
		for _, stmt := range s.Alternative.Statements {
			c.checkStatement(stmt)
		}
	}
}

func (c *Checker) checkExpression(expr parser.Expression) Type {
	if expr == nil {
		return TypeUnknown
	}

	switch e := expr.(type) {
	case *parser.IntegerLiteral:
		return TypeInt

	case *parser.StringLiteral:
		return TypeString

	case *parser.BoolLiteral:
		return TypeBool

	case *parser.Identifier:
		if t, ok := c.variables[e.Value]; ok {
			return t
		}
		c.errors = append(c.errors, fmt.Sprintf("undefined variable: %q", e.Value))
		return TypeUnknown

	case *parser.BinaryExpression:
		return c.checkBinaryExpression(e)

	case *parser.UnaryExpression:
		rightType := c.checkExpression(e.Right)
		if e.Operator == "!" && rightType != TypeBool {
			c.errors = append(c.errors, fmt.Sprintf("! operator requires bool, got %s", rightType))
		}
		if e.Operator == "-" && rightType != TypeInt {
			c.errors = append(c.errors, fmt.Sprintf("- operator requires int, got %s", rightType))
		}
		return rightType

	case *parser.CallExpression:
		return c.checkCallExpression(e)
	}

	return TypeUnknown
}

func (c *Checker) checkBinaryExpression(e *parser.BinaryExpression) Type {
	leftType := c.checkExpression(e.Left)
	rightType := c.checkExpression(e.Right)

	// Comparison operators return bool
	switch e.Operator {
	case "==", "!=", "<", ">", "<=", ">=":
		if leftType != rightType && leftType != TypeUnknown && rightType != TypeUnknown {
			c.errors = append(c.errors, fmt.Sprintf(
				"cannot compare %s and %s", leftType, rightType,
			))
		}
		return TypeBool

	case "&&", "||":
		if leftType != TypeBool && leftType != TypeUnknown {
			c.errors = append(c.errors, fmt.Sprintf(
				"%s operator requires bool operands, got %s", e.Operator, leftType,
			))
		}
		return TypeBool

	case "+":
		// Allow string concatenation
		if leftType == TypeString && rightType == TypeString {
			return TypeString
		}
		if leftType != TypeInt && leftType != TypeUnknown {
			c.errors = append(c.errors, fmt.Sprintf(
				"+ operator requires int or string, got %s and %s", leftType, rightType,
			))
		}
		if leftType != rightType && leftType != TypeUnknown && rightType != TypeUnknown {
			c.errors = append(c.errors, fmt.Sprintf(
				"type mismatch: cannot use + with %s and %s", leftType, rightType,
			))
		}
		return leftType

	case "-", "*", "/":
		if (leftType != TypeInt && leftType != TypeUnknown) ||
			(rightType != TypeInt && rightType != TypeUnknown) {
			c.errors = append(c.errors, fmt.Sprintf(
				"%s operator requires int operands, got %s and %s",
				e.Operator, leftType, rightType,
			))
		}
		return TypeInt
	}

	return TypeUnknown
}

func (c *Checker) checkCallExpression(e *parser.CallExpression) Type {
	fn, ok := c.functions[e.Function]
	if !ok {
		c.errors = append(c.errors, fmt.Sprintf("undefined function: %q", e.Function))
		return TypeUnknown
	}

	if len(e.Arguments) != len(fn.Params) {
		c.errors = append(c.errors, fmt.Sprintf(
			"function %q expects %d arguments, got %d",
			e.Function, len(fn.Params), len(e.Arguments),
		))
		return fn.ReturnType
	}

	for i, arg := range e.Arguments {
		argType := c.checkExpression(arg)
		if argType != fn.Params[i] && argType != TypeUnknown {
			c.errors = append(c.errors, fmt.Sprintf(
				"argument %d of %q: expected %s, got %s",
				i+1, e.Function, fn.Params[i], argType,
			))
		}
	}

	return fn.ReturnType
}
