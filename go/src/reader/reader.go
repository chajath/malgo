package reader

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/chajath/malgo/go/types"
)

type reader struct {
	tokens []string
	pos    int
}

var re = regexp.MustCompile(`[\s,]*(~@|[\[\]{}()'` + "`" +
	`~^@]|"(?:\\.|[^\\"])*"?|;.*|[^\s\[\]{}('"` + "`" +
	`,;)]*)`)

func newReader(str string) *reader {
	tokens := tokenize(str)
	return &reader{tokens: tokens, pos: 0}
}

func tokenize(str string) []string {
	results := make([]string, 0)
	for _, group := range re.FindAllStringSubmatch(str, -1) {
		if (group[1] == "") || (group[1][0] == ';') {
			// Ignore whitespaces and comments.
			continue
		}
		results = append(results, group[1])
	}
	return results
}

// ReadStr turns the string into a MalType
func ReadStr(str string) (types.MalType, error) {
	r := newReader(str)
	return readForm(r)
}

func readForm(r *reader) (types.MalType, error) {
	p, err := r.Peek()
	if err != nil {
		return nil, err
	}
	if p == "(" {
		return readList(r)
	}
	return readAtom(r)
}

func readList(r *reader) (*types.MalList, error) {
	r.Next()
	toReturn := make([]types.MalType, 0)
	for {
		p, err := r.Peek()
		if err != nil {
			return nil, err
		}
		if p == ")" {
			break
		}
		m, err := readForm(r)
		if err != nil {
			return nil, err
		}
		toReturn = append(toReturn, m)
	}
	r.Next()
	return types.NewMalList(toReturn), nil
}

func readAtom(r *reader) (types.MalType, error) {
	n := r.Next()
	tryNum, err := strconv.Atoi(n)
	if err != nil {
		return types.NewMalSymbol(n), nil
	}
	return types.NewMalNumber(tryNum), nil
}

func (r *reader) Next() string {
	var toReturn = r.tokens[r.pos]
	r.pos = r.pos + 1
	return toReturn
}

func (r *reader) Peek() (string, error) {
	if r.pos >= len(r.tokens) {
		return "", errors.New("EOF reached")
	}
	return r.tokens[r.pos], nil
}
