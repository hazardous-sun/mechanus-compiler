# Mechanus Compiler

> "Mechanus is unnatural — bottom-to-top, left-to-right."  
> — The Architect

## 📜 Overview

**Mechanus** is an esoteric programming language with a compiler implemented in Go. It embraces an unconventional syntax
where code flows in reverse — literally from **bottom to top** and **right to left**.

This repository contains:

- A compiler for the Mechanus language.
- Several example `.mecha` programs.
- A formal grammar definition (derivation tree).
- A script to compile all examples in one go.

### 📌 Project Status

- ✅ **Lexer**: Fully implemented — tokenizes input source code.
- ✅ **Syntax Analyzer**: Fully implemented — validates syntax using a recursive-descent parser.
- 🔄 **Semantic Analyzer**: *In progress* — under development to validate meaning and context.
- ⏳ **Code Generation**: *Coming soon* — will emit target code from validated AST.

---

## 🧷 Keywords

Mechanus has a collection of keywords and types that define its unconventional syntax. Here's a categorized reference:

### 🔨 Language Constructs

| Keyword     | Description                                     |
|-------------|-------------------------------------------------|
| `Construct` | Declares a module                               |
| `Architect` | Declares a function                             |
| `Integrate` | Equivalent to `return` in traditional languages |

---

### 🔁 Control Flow

| Keyword | Description       |
|---------|-------------------|
| `if`    | Conditional check |
| `else`  | Else block        |
| `elif`  | Else-if branch    |
| `for`   | Loop construct    |

---

### 📦 Data Types

| Type        | Description                          |
|-------------|--------------------------------------|
| `Nil`       | Empty or no value                    |
| `Gear`      | Integer values                       |
| `Tensor`    | Floating-point numbers               |
| `State`     | User-defined state-like values       |
| `Monodrone` | Single character string (1 char max) |
| `Omnidrone` | Regular string literal               |

---

### 📤 Built-in Functions

| Function  | Description           |
|-----------|-----------------------|
| `Send`    | Sends a value/message |
| `Receive` | Receives a value      |

---

### ⚙️ Operators & Symbols

| Symbol            | Description                                  |
|-------------------|----------------------------------------------|
| `=:`              | Declaration (type + variable)                |
| `=`               | Assignment                                   |
| `+ - * / %`       | Arithmetic operators                         |
| `== <> <= >= < >` | Comparison operators                         |
| `(` `)`           | Parentheses                                  |
| `{` `}`           | Block delimiters                             |
| `:` `,`           | Type/parameter delimiters                    |
| `'` `"`           | String delimiters (`Monodrone`, `Omnidrone`) |
| `//`              | Single-line comment                          |

---

## 📂 Project Structure

```
.
├── docs/
│   ├── derivation_tree.md        # Formal grammar of Mechanus
│   └── examples/                 # .mecha input and output example files
├── src/                          # Compiler source code
│   ├── cmd/                      # Entry point (`main.go`)
│   ├── error/                    # Error handling and logging utilities
│   └── models/                   # Lexer and parser implementations
├── run.sh                        # Script to run compiler on all examples
├── go.mod                        # Go module definition
└── README.md                     # This file
``` 