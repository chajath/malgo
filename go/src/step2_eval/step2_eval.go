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
	env := newEnv()
	for scanner.Scan() {
		in := scanner.Text()
		repOut, err := rep(in, env)
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

	env := newEnv()
	for {
		in, err := rl.Read()
		if err != nil {
			break
		}
		repOut, err := rep(in, env)
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

type biNumFunc func(int, int) int

type env struct {
	replEnv map[string]evaller
}

type evaller interface {
	eval(types.MalType, *env) (types.MalType, error)
}

type numEvaller struct {
	biNumFunc biNumFunc
}

func newEnv() *env {
	replEnv := map[string]evaller{
		"+": &numEvaller{biNumFunc: plusFunc},
		"-": &numEvaller{biNumFunc: minusFunc},
		"*": &numEvaller{biNumFunc: multiFunc},
		"/": &numEvaller{biNumFunc: divFunc}}
	return &env{replEnv: replEnv}
}

func plusFunc(a, b int) int {
	return a + b
}

func minusFunc(a, b int) int {
	return a - b
}

func multiFunc(a, b int) int {
	return a * b
}

func divFunc(a, b int) int {
	return a / b
}

func (ne *numEvaller) eval(rest types.MalType, env *env) (types.MalType, error) {
	l, ok := rest.(*types.MalList)
	if !ok {
		return nil, fmt.Errorf("eval: expecting list, given %+v", rest)
	}
	firstEval, err := eval(l.List[0], env)
	if err != nil {
		return nil, err
	}
	first, ok := firstEval.(*types.MalNumber)
	if !ok {
		return nil, fmt.Errorf("eval: expecting number, given %+v", l.List[0])
	}
	resultNum := first.Number
	for _, m := range l.List[1:] {
		localResult, err := eval(m, env)
		if err != nil {
			return nil, err
		}
		localNumber, ok := localResult.(*types.MalNumber)
		if !ok {
			return nil, fmt.Errorf("eval: expecting number, given %+v", localResult)
		}
		resultNum = ne.biNumFunc(resultNum, localNumber.Number)
	}
	return types.NewMalNumber(resultNum), nil
}

func eval(in types.MalType, env *env) (types.MalType, error) {
	switch in.(type) {
	case *types.MalList:
		l := in.(*types.MalList)
		// Return empty list as is.
		if len(l.List) == 0 {
			return l, nil
		}

		// Evaluate head.
		head, err := eval(l.List[0], env)
		if err != nil {
			return nil, err
		}

		headSym, ok := head.(*types.MalSymbol)
		if !ok {
			return nil, fmt.Errorf("eval: expecting symbol for list head, given %+v", head)
		}

		// Look up built-in evaller to call.
		evaller, ok := env.replEnv[headSym.Symbol]
		if !ok {
			return nil, fmt.Errorf("eval: built-in function not found for %+v", headSym)
		}

		// Call evaller with the rest of the list.
		return evaller.eval(types.NewMalList(l.List[1:]), env)
	}

	// Other atoms are no-op at eval.
	return in, nil
}

func print(in types.MalType) (string, error) {
	return printer.PrStr(in)
}

func rep(in string, env *env) (string, error) {
	readOut, err := read(in)
	if err != nil {
		return "", err
	}
	evalOut, err := eval(readOut, env)
	if err != nil {
		return "", err
	}
	printOut, err := print(evalOut)
	if err != nil {
		return "", err
	}
	return printOut, nil
}
