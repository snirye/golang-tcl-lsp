package main

import (
	"fmt"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

// AST represents the abstract syntax tree for a TCL program.
type AST struct {
	Statements []Statement
}

// Statement represents a statement in a TCL program.
type Statement interface {
	Statement()
}

// Command represents a command in a TCL program.
type Command struct {
	Name      string
	Arguments []Expression
}

func (c *Command) Statement() {}

// Expression represents an expression in a TCL program.
type Expression interface {
	Expression()
}

type UnquotedExpression struct {
	Value string `@Ident`
}

func (UnquotedExpression) Expression() {}

type Number struct {
	Number float64 `@Float`
}

func (Number) Expression() {}

type BracesGroup struct {
	Value string `"{" @Ident "}"`
}

func (BracesGroup) Expression() {}

type DoubleQuotesGroup struct {
	String string `@String`
}

func (DoubleQuotesGroup) Expression() {}

type SquareBracketsGroup struct {
	Value string `"[" @Ident "]"`
}

func (SquareBracketsGroup) Expression() {}

func main() {
	// Parse a TCL program.
	input := `
set foo 123
set bar "hello world"
`

	// see this example:
	// https://github.com/openllb/hlb/blob/9f1194235a5f3bd3ab123239b5e517344a293f31/parser/ast/ast.go
	var lexer2 = lexer.Must(lexer.Rules{
		"Root": {
			{`BracesGroup`, `{`, lexer.Push("BracesGroup")},
		},
		"BracesGroup": {
			{`BracesGroup`, `{`, lexer.Push("BracesGroup")},
			{"BracesGroupEnd", `}`, lexer.Pop()},
			{"Ident", `(.|\n)+`, nil},
		},
	})

	parser, err := participle.Build[AST](
		participle.Unquote("String"),
		participle.Union[Expression](UnquotedExpression{}, Number{}, BracesGroup{}, DoubleQuotesGroup{}, SquareBracketsGroup{}),
	)

	// Access the parsed data.
	for _, statement := range ast.Statements {
		switch statement := statement.(type) {
		case *Command:
			fmt.Println("Command:", statement.Name, statement.Arguments)
		case *Expression:
			fmt.Println("Expression:", statement)
		}
	}
}
