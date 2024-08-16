package main

import (
	"eventloop/interpreter/interpreter/evaluator"
	"eventloop/interpreter/interpreter/lexer"
	"eventloop/interpreter/interpreter/object"
	"eventloop/interpreter/interpreter/parser"
	"eventloop/interpreter/interpreter/repl"
	"io"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	out := os.Stdout
	input, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic("meow")
	}

	l := lexer.New(string(input))
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	if len(p.Errors()) != 0 {
		printParserErrors(os.Stdout, p.Errors())
		return
	}

	evaluated := evaluator.Eval(program, env)
	if evaluated != nil {
		io.WriteString(out, evaluated.Inspect())
		io.WriteString(out, "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
