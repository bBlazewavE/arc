package main

import (
	"fmt"
	"os"

	"github.com/bBlazewavE/arc/interpreter"
	"github.com/bBlazewavE/arc/lexer"
	"github.com/bBlazewavE/arc/parser"
	"github.com/bBlazewavE/arc/typechecker"
)

const version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Arc Programming Language v" + version)
		fmt.Println("Usage: arc <filename.arc>")
		fmt.Println()
		fmt.Println("Example:")
		fmt.Println("  arc hello.arc")
		os.Exit(0)
	}

	filename := os.Args[1]
	source, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err)
		os.Exit(1)
	}

	// Phase 1: Lexing
	l := lexer.New(string(source))
	tokens := l.Tokenize()

	// Phase 2: Parsing
	p := parser.New(tokens)
	program := p.Parse()

	if len(p.Errors()) > 0 {
		fmt.Fprintln(os.Stderr, "Parse errors:")
		for _, e := range p.Errors() {
			fmt.Fprintln(os.Stderr, "  "+e)
		}
		os.Exit(1)
	}

	// Phase 3: Type Checking (static!)
	tc := typechecker.New()
	tc.Check(program)

	if len(tc.Errors()) > 0 {
		fmt.Fprintln(os.Stderr, "Type errors:")
		for _, e := range tc.Errors() {
			fmt.Fprintln(os.Stderr, "  "+e)
		}
		os.Exit(1)
	}

	// Phase 4: Execution
	interp := interpreter.New()
	interp.Run(program)
}
