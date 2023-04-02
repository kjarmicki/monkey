package main

import (
	"fmt"
	"os"

	"kjarmicki.github.com/monkey/repl"
)

func main() {
	fmt.Println("Monkey REPL ready")
	repl.Start(os.Stdin, os.Stdout)
}
