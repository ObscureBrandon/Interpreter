package main

import (
	"eventloop/interpreter/interpreter/repl"
	"os"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
