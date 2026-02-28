// Hello, Arc! 🔥
// A simple statically-typed programming language

let name: string = "Dimple"
let age: int = 28
let is_engineer: bool = true

print(name)
print(age)

// Arithmetic
let x: int = 10
let y: int = 20
let sum: int = x + y
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

fn greet(who: string) -> string {
    return "Welcome to Arc, " + who
}

let result: int = add(x, y)
print(result)

let msg: string = greet(name)
print(msg)

// Nested expressions
let complex: int = add(x * 2, y + 5)
print(complex)

// Boolean logic
if is_engineer && age > 20 {
    print("Staff Engineer incoming!")
}
