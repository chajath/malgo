package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/chajath/malgo/go/env"
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

type numEvaluator struct {
	biNumFunc biNumFunc
	types.MalType
}

type evaluator interface {
	eval(types.MalType, *env.MalEnv) (types.MalType, error)
}

func newEnv() env.MalEnv {
	root := env.New(nil)
	root.Set("+", &numEvaluator{biNumFunc: plusFunc})
	root.Set("-", &numEvaluator{biNumFunc: minusFunc})
	root.Set("*", &numEvaluator{biNumFunc: multiFunc})
	root.Set("/", &numEvaluator{biNumFunc: divFunc})

	return root
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

func (ne *numEvaluator) eval(rest types.MalType, env env.MalEnv) (types.MalType, error) {
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

func eval(in types.MalType, cEnv env.MalEnv) (types.MalType, error) {
	switch in.(type) {
	case *types.MalSymbol:
		sym := in.(*types.MalSymbol)
		return cEnv.Get(sym.Symbol)
	case *types.MalList:
		l := in.(*types.MalList)
		// Return empty list as is.
		if len(l.List) == 0 {
			return l, nil
		}

		// Get head symbol.
		head := l.List[0]

		headSym, ok := head.(*types.MalSymbol)
		if !ok {
			return nil, fmt.Errorf("eval: expecting symbol for list head, given %+v", head)
		}

		// Handling specials.
		switch headSym.Symbol {
		case "def!":
			if len(l.List) != 3 {
				return nil, fmt.Errorf("eval: wrong number of arguments for def!: %+v", l)
			}
			key, ok := l.List[1].(*types.MalSymbol)
			if !ok {
				return nil, fmt.Errorf("eval: key to def! is not symbol: %+v", l.List[1])
			}
			value, err := eval(l.List[2], cEnv)
			if err != nil {
				return nil, err
			}
			cEnv.Set(key.Symbol, value)
			return value, nil
		case "let*":
			if len(l.List) != 3 {
				return nil, fmt.Errorf("eval: wrong number of arguments for let*: %+v", l)
			}
			argList, ok := l.List[1].(*types.MalList)
			if !ok {
				return nil, fmt.Errorf("eval: first argument to let* is not a list: %+v", l.List[1])
			}
			if len(argList.List)%2 != 0 {
				return nil, fmt.Errorf("eval: number of elements in binding list is not even: %+v", argList)
			}
			localEnv := env.New(cEnv)
			for i := 0; i < len(argList.List); i += 2 {
				key, ok := argList.List[i].(*types.MalSymbol)
				if !ok {
					return nil, fmt.Errorf("eval: %+v is not a symbol", key)
				}
				value, err := eval(argList.List[i+1], localEnv)
				if err != nil {
					return nil, err
				}
				localEnv.Set(key.Symbol, value)
			}
			return eval(l.List[2], localEnv)
		}

		// Look up built-in evaluator to call.
		headValue, err := cEnv.Get(headSym.Symbol)
		if err != nil {
			return nil, err
		}
		evaluator, ok := headValue.(*numEvaluator)

		// Call evaluator with the rest of the list.
		return evaluator.eval(types.NewMalList(l.List[1:]), cEnv)
	}

	// Other atoms are no-op at eval.
	return in, nil
}

func print(in types.MalType) (string, error) {
	return printer.PrStr(in)
}

func rep(in string, env env.MalEnv) (string, error) {
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
