package interpreter

import (
	"fmt"

	"github.com/bBlazewavE/arc/parser"
)

// Value represents a runtime value
type Value struct {
	Type     string // "int", "string", "bool"
	IntVal   int64
	StrVal   string
	BoolVal  bool
}

func (v Value) String() string {
	switch v.Type {
	case "int":
		return fmt.Sprintf("%d", v.IntVal)
	case "string":
		return v.StrVal
	case "bool":
		if v.BoolVal {
			return "true"
		}
		return "false"
	}
	return "nil"
}

// Function stores a function definition
type Function struct {
	Params []parser.Param
	Body   *parser.BlockStatement
}

// Environment holds variables and functions
type Environment struct {
	variables map[string]Value
	functions map[string]Function
	parent    *Environment
}

func NewEnvironment(parent *Environment) *Environment {
	return &Environment{
		variables: make(map[string]Value),
		functions: make(map[string]Function),
		parent:    parent,
	}
}

func (env *Environment) Get(name string) (Value, bool) {
	val, ok := env.variables[name]
	if !ok && env.parent != nil {
		return env.parent.Get(name)
	}
	return val, ok
}

func (env *Environment) Set(name string, val Value) {
	env.variables[name] = val
}

func (env *Environment) GetFunc(name string) (Function, bool) {
	fn, ok := env.functions[name]
	if !ok && env.parent != nil {
		return env.parent.GetFunc(name)
	}
	return fn, ok
}

// ReturnSignal is used to unwind the stack on return
type ReturnSignal struct {
	Value Value
}

// Interpreter executes the AST
type Interpreter struct {
	env    *Environment
	output []string // captured output for testing
}

// New creates a new Interpreter
func New() *Interpreter {
	return &Interpreter{
		env: NewEnvironment(nil),
	}
}

// Output returns captured print output
func (i *Interpreter) Output() []string {
	return i.output
}

// Run executes a program
func (i *Interpreter) Run(program *parser.Program) {
	for _, stmt := range program.Statements {
		result := i.execStatement(stmt)
		if _, ok := result.(*ReturnSignal); ok {
			return
		}
	}
}

func (i *Interpreter) execStatement(stmt parser.Statement) interface{} {
	switch s := stmt.(type) {
	case *parser.LetStatement:
		val := i.evalExpression(s.Value)
		i.env.Set(s.Name, val)

	case *parser.FnStatement:
		i.env.functions[s.Name] = Function{
			Params: s.Params,
			Body:   s.Body,
		}

	case *parser.PrintStatement:
		val := i.evalExpression(s.Value)
		line := val.String()
		fmt.Println(line)
		i.output = append(i.output, line)

	case *parser.IfStatement:
		cond := i.evalExpression(s.Condition)
		if cond.BoolVal {
			return i.execBlock(s.Consequence)
		} else if s.Alternative != nil {
			return i.execBlock(s.Alternative)
		}

	case *parser.ReturnStatement:
		val := i.evalExpression(s.Value)
		return &ReturnSignal{Value: val}

	case *parser.ExpressionStatement:
		i.evalExpression(s.Expression)
	}

	return nil
}

func (i *Interpreter) execBlock(block *parser.BlockStatement) interface{} {
	for _, stmt := range block.Statements {
		result := i.execStatement(stmt)
		if ret, ok := result.(*ReturnSignal); ok {
			return ret
		}
	}
	return nil
}

func (i *Interpreter) evalExpression(expr parser.Expression) Value {
	if expr == nil {
		return Value{Type: "int", IntVal: 0}
	}

	switch e := expr.(type) {
	case *parser.IntegerLiteral:
		return Value{Type: "int", IntVal: e.Value}

	case *parser.StringLiteral:
		return Value{Type: "string", StrVal: e.Value}

	case *parser.BoolLiteral:
		return Value{Type: "bool", BoolVal: e.Value}

	case *parser.Identifier:
		val, ok := i.env.Get(e.Value)
		if !ok {
			fmt.Printf("Runtime error: undefined variable %q\n", e.Value)
			return Value{}
		}
		return val

	case *parser.BinaryExpression:
		return i.evalBinary(e)

	case *parser.UnaryExpression:
		right := i.evalExpression(e.Right)
		if e.Operator == "!" {
			return Value{Type: "bool", BoolVal: !right.BoolVal}
		}
		if e.Operator == "-" {
			return Value{Type: "int", IntVal: -right.IntVal}
		}

	case *parser.CallExpression:
		return i.evalCall(e)
	}

	return Value{}
}

func (i *Interpreter) evalBinary(e *parser.BinaryExpression) Value {
	left := i.evalExpression(e.Left)
	right := i.evalExpression(e.Right)

	switch e.Operator {
	// Arithmetic (int)
	case "+":
		if left.Type == "string" {
			return Value{Type: "string", StrVal: left.StrVal + right.StrVal}
		}
		return Value{Type: "int", IntVal: left.IntVal + right.IntVal}
	case "-":
		return Value{Type: "int", IntVal: left.IntVal - right.IntVal}
	case "*":
		return Value{Type: "int", IntVal: left.IntVal * right.IntVal}
	case "/":
		if right.IntVal == 0 {
			fmt.Println("Runtime error: division by zero")
			return Value{Type: "int", IntVal: 0}
		}
		return Value{Type: "int", IntVal: left.IntVal / right.IntVal}

	// Comparison
	case "==":
		return Value{Type: "bool", BoolVal: i.valuesEqual(left, right)}
	case "!=":
		return Value{Type: "bool", BoolVal: !i.valuesEqual(left, right)}
	case "<":
		return Value{Type: "bool", BoolVal: left.IntVal < right.IntVal}
	case ">":
		return Value{Type: "bool", BoolVal: left.IntVal > right.IntVal}
	case "<=":
		return Value{Type: "bool", BoolVal: left.IntVal <= right.IntVal}
	case ">=":
		return Value{Type: "bool", BoolVal: left.IntVal >= right.IntVal}

	// Logical
	case "&&":
		return Value{Type: "bool", BoolVal: left.BoolVal && right.BoolVal}
	case "||":
		return Value{Type: "bool", BoolVal: left.BoolVal || right.BoolVal}
	}

	return Value{}
}

func (i *Interpreter) valuesEqual(a, b Value) bool {
	if a.Type != b.Type {
		return false
	}
	switch a.Type {
	case "int":
		return a.IntVal == b.IntVal
	case "string":
		return a.StrVal == b.StrVal
	case "bool":
		return a.BoolVal == b.BoolVal
	}
	return false
}

func (i *Interpreter) evalCall(e *parser.CallExpression) Value {
	fn, ok := i.env.GetFunc(e.Function)
	if !ok {
		fmt.Printf("Runtime error: undefined function %q\n", e.Function)
		return Value{}
	}

	// Create new scope for function
	funcEnv := NewEnvironment(i.env)

	// Bind arguments to parameters
	for idx, param := range fn.Params {
		if idx < len(e.Arguments) {
			val := i.evalExpression(e.Arguments[idx])
			funcEnv.variables[param.Name] = val
		}
	}

	// Execute function body in new scope
	prevEnv := i.env
	i.env = funcEnv

	var returnVal Value
	for _, stmt := range fn.Body.Statements {
		result := i.execStatement(stmt)
		if ret, ok := result.(*ReturnSignal); ok {
			returnVal = ret.Value
			break
		}
	}

	i.env = prevEnv
	return returnVal
}
