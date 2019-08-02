package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	// Main REPL loop.
	fmt.Print("user> ")
	for scanner.Scan() {
		repOut := rep(scanner.Text())
		fmt.Println(repOut)
		fmt.Print("user> ")
	}
}

func read(in string) string {
	return in
}

func eval(in string) string {
	return in
}

func print(in string) string {
	return in
}

func rep(in string) string {
	readOut := read(in)
	evalOut := eval(readOut)
	printOut := print(evalOut)
	return printOut
}
