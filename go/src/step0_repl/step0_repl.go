package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/chajath/malgo/go/readline"
)

var rl bool

func main() {
	flag.BoolVar(&rl, "rl", false, "Use readline")
	flag.Parse()
	if rl {
		mainRl()
		return
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("user> ")
	for scanner.Scan() {
		in := scanner.Text()
		repOut := rep(in)
		fmt.Println(repOut)
		fmt.Print("user> ")
	}
}

func mainRl() {
	rl, err := readline.NewReadline("user> ")
	if err != nil {
		return
	}

	defer rl.Close()

	for {
		in, err := rl.Read()
		if err != nil {
			break
		}
		repOut := rep(in)
		fmt.Println(repOut)
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
