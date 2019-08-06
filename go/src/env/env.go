package env

import (
	"fmt"

	"github.com/chajath/malgo/go/types"
)

type malEnv struct {
	data  map[string]types.MalType
	outer MalEnv
}

// New creates a new MalEnv with outer reference.
func New(outer MalEnv) MalEnv {
	return &malEnv{outer: outer, data: make(map[string]types.MalType)}
}

// MalEnv holds the reference to local environment and pointer to the outer environment.
type MalEnv interface {
	Set(string, types.MalType)
	Find(string) types.MalType
	Get(string) (types.MalType, error)
}

func (env *malEnv) Set(symbol string, m types.MalType) {
	env.data[symbol] = m
}

func (env *malEnv) Find(symbol string) types.MalType {
	lookUp, ok := env.data[symbol]
	if ok {
		return lookUp
	}
	if env.outer != nil {
		return env.outer.Find(symbol)
	}
	return nil
}

func (env *malEnv) Get(symbol string) (types.MalType, error) {
	lookUp := env.Find(symbol)
	if lookUp != nil {
		return lookUp, nil
	}
	return nil, fmt.Errorf("env: symbol %+v not found", symbol)
}
