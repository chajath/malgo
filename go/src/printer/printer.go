package printer

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/chajath/malgo/go/types"
)

// PrStr pretty prints the given MalType
func PrStr(m types.MalType) (string, error) {
	switch m.(type) {
	case *types.MalSymbol:
		return m.(*types.MalSymbol).Symbol, nil
	case *types.MalNumber:
		return strconv.Itoa(m.(*types.MalNumber).Number), nil
	case *types.MalList:
		var sb strings.Builder
		sb.WriteString("(")
		sep := ""
		for _, mElem := range m.(*types.MalList).List {
			sb.WriteString(sep)
			pr, err := PrStr(mElem)
			if err != nil {
				return "", err
			}
			sb.WriteString(pr)
			sep = " "
		}
		sb.WriteString(")")
		return sb.String(), nil
	case *types.MalTrue:
		return "true", nil
	case *types.MalFalse:
		return "false", nil
	case *types.MalNil:
		return "nil", nil
	}
	return "", fmt.Errorf("printer: unknown MalType %+v", m)
}
