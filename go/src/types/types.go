package types

// MalType is the top-level interface.
type MalType interface{}

// MalList contains a list of MalTypes.
type MalList struct {
	List []MalType
	MalType
}

// NewMalList constructs MalList.
func NewMalList(list []MalType) *MalList {
	return &MalList{List: list}
}

// MalSymbol is a string representation of a symbol.
type MalSymbol struct {
	Symbol string
	MalType
}

// NewMalSymbol constructs MalSymbol
func NewMalSymbol(symbol string) *MalSymbol {
	return &MalSymbol{Symbol: symbol}
}

// MalNumber contains a number atom.
type MalNumber struct {
	Number int
	MalType
}

// NewMalNumber constructs MalNumber
func NewMalNumber(number int) *MalNumber {
	return &MalNumber{Number: number}
}
