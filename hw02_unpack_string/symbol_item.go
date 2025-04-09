package hw02unpackstring

import (
	"strconv"
	"strings"
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
	isSlashed := false
	for j := itemNumber - 1; j >= 0; j-- {
		sItem := inRunes[j]
		sItemIsSlash := (sItem == SymbolSlash)
		if sItemIsSlash {
			isSlashed = !isSlashed
		} else {
			break
		}
	}
	return isSlashed
}

func (si *SymbolItem) RepeatWith(nextRef *SymbolItem) string {
	var sb strings.Builder
	if nextRef.Type == IsDigit { // если следующий символ некая цифра x, то записать текущий символ x раз
		sb.WriteString(strings.Repeat(string(si.Item), nextRef.ValueInt))
	} else { // если следующий символ не цифра, то записать текущий символ 1 раз
		sb.WriteRune(si.Item)
	}
	return sb.String()
}
