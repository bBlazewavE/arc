# Arc ⚡

A simple statically-typed programming language built in Go.

## Features

- **Static type checking** — catches type errors before runtime
- **Types:** `int`, `string`, `bool`
- **Functions** with typed parameters and return types
- **If/else** control flow
- **Arithmetic** and **boolean** operators
- **String concatenation**
- **Comments** with `//`

## Quick Start

```bash
go build -o arc .
./arc examples/hello.arc
```

## Syntax

```arc
// Variables with types
let name: string = "Dimple"
let age: int = 28
let active: bool = true

// Arithmetic
let sum: int = 10 + 20
print(sum)

// String concatenation
let greeting: string = "Hello, " + name
print(greeting)

// If/else
if age > 25 {
    print("experienced!")
} else {
    print("still learning!")
}

// Functions
fn add(a: int, b: int) -> int {
    return a + b
}

let result: int = add(10, 20)
print(result)
```

## Architecture

```
Source Code (.arc)
    ↓
[Lexer]        → Tokens
    ↓
[Parser]       → AST (Abstract Syntax Tree)
    ↓
[Type Checker] → Static type validation
    ↓
[Interpreter]  → Execution
```

## Type System

Arc catches type errors at compile time:

```arc
let x: int = "hello"     // ❌ Type error: declared as int but assigned string
let y: int = 10 + "hi"   // ❌ Type error: cannot use + with int and string
```

## Project Structure

```
arc/
├── main.go              # Entry point
├── lexer/
│   ├── token.go         # Token types and keywords
│   └── lexer.go         # Tokenizer
├── parser/
│   ├── ast.go           # AST node definitions
│   └── parser.go        # Pratt parser
├── typechecker/
│   └── checker.go       # Static type checker
├── interpreter/
│   └── interpreter.go   # Tree-walking interpreter
└── examples/
    └── hello.arc        # Example program
```

## Built by

Dimple Mathew — as a learning project for understanding compilers and Go.
