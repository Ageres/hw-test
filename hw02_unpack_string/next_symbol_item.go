package hw02unpackstring

import (
	"strconv"
	"unicode"
)

// структура с параметрами символа, следующего за анализируемым символом.
type NextSymbolItem struct {
	Item     rune // следующий символ
	IsDigit  bool // следующий символ - это цифра
	ValueInt int  // численное значение следующего символа, если это цифра
}

func BuildNextSymBolItem(nextItem rune) (*NextSymbolItem, error) {
	nsi := NextSymbolItem{}
	nsi.Item = nextItem
	nsi.IsDigit = unicode.IsDigit(nextItem)
	if nsi.IsDigit {
		valueInt, err := strconv.Atoi(string(nextItem))
		if err != nil {
			return nil, err
		}
		nsi.ValueInt = valueInt
	}
	return &nsi, nil
}
