package readline

import (
	"github.com/chzyer/readline"
)

// MalReadline wraps around inner Readline struct.
type MalReadline struct {
	rl *readline.Instance
}

func NewReadline(prompt string) (*MalReadline, error) {
	rl, err := readline.New(prompt)
	if err != nil {
		return nil, err
	}
	return &MalReadline{rl: rl}, nil
}

func (mrl MalReadline) Read() (string, error) {
	return mrl.rl.Readline()
}

func (mrl MalReadline) Close() error {
	return mrl.rl.Close()
}
