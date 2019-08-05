package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/chajath/malgo/go/printer"
	"github.com/chajath/malgo/go/reader"
	"github.com/chajath/malgo/go/types"

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
		repOut, err := rep(in)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(repOut)
		}
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
		repOut, err := rep(in)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(repOut)
		}
	}
}

func read(in string) (types.MalType, error) {
	m, err := reader.ReadStr(in)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func eval(in types.MalType) types.MalType {
	return in
}

func print(in types.MalType) (string, error) {
	return printer.PrStr(in)
}

func rep(in string) (string, error) {
	readOut, err := read(in)
	if err != nil {
		return "", err
	}
	evalOut := eval(readOut)
	printOut, err := print(evalOut)
	if err != nil {
		return "", err
	}
	return printOut, nil
}
