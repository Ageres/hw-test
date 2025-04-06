package hw02unpackstring

import (
	"strconv"
	"unicode"
)

var SymbolSlash = []rune(`\`)[0]

// перечисление с типами анализируемого символа.
type Type int

const (
	IsDigit Type = iota + 1 // анализируемый символ это цифра
	IsSlash                 // анализируемый символ это слеш
	IsOther                 // анализируемый символ это не цифра и не слеш
)

// структура с параметрами анализируемого символа.
type SymbolItem struct {
	Item      rune // анализируемый символ
	Type      Type // тип анализируемого символа
	IsSlashed bool // анализируемый символ экранирован
	ValueInt  int  // численное значение анализируемого символа, если это цифра
}

func BuildSymbolItem(itemNumber int, inRunes []rune) *SymbolItem {
	si := SymbolItem{}
	si.Item = inRunes[itemNumber]
	switch {
	case unicode.IsDigit(si.Item):
		si.Type = IsDigit
	case si.Item == SymbolSlash:
		si.Type = IsSlash
	default:
		si.Type = IsOther
	}
	si.IsSlashed = defineIfItemIsSlashed(itemNumber, inRunes) // определение экранирован ли текущий символ

	return &si
}

func (si *SymbolItem) ParseIfDigit() error {
	if si.Type == IsDigit {
		valueInt, err := strconv.Atoi(string(si.Item))
		if err != nil {
			return err
		}
		si.ValueInt = valueInt
	}
	return nil
}

// определение экранирован ли символ.
func defineIfItemIsSlashed(itemNumber int, inRunes []rune) bool {
	// подсчет количества предыдущих символов слеш, следующих подряд
	countPreviousSlash := 0
	for j := itemNumber - 1; j >= 0; j-- {
		sItem := inRunes[j]
		sItemIsSlash := (sItem == SymbolSlash)
		if sItemIsSlash {
			countPreviousSlash++
		} else {
			break
		}
	}

	return !(countPreviousSlash%2 == 0)
}
