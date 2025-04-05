package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")
var SymbolSlash = []rune(`\`)[0]

/*
 * Распаковка переданной строки в 4 этапа:
 * Первый этап    - проверяется есть ли символы в переданной строке, если нет - возвращается пустая строка;
 * Второй этап    - проверяется достаточность условий и формат распаковки первого символа на основе анализа вида
 *                  и значения первого и второго символов (число/слеш/прочий символ);
 * Третий этап    - проверяется достаточность условий и формат распаковки последовательно от второго до предпоследнего
 *	                символа, решение для анализируемого символа принимается на основании собственного вида и значения,
 *                  а так же видов и значений предыдущих и последующего символов;
 *				    при необходимости первый символ так же добавляется в распаковку;
 * Четвертый этап - проверяется достаточность условий и формат распаковки последнего символа на основании собственного
 *                  вида и значения, а так же видов и значений предыдущих символов.
 */
func Unpack(in string) (string, error) {
	// Place your code here.
	inRunes := []rune(in)
	inSize := len(inRunes)
	// первый этап - анализируемая строка содержит 0 символов
	if inSize == 0 {
		return "", nil
	}
	// второй этап - анализ первого символа
	firstItem := inRunes[0]
	// если первый символ строки содержит цифру, то вернуть ошибку
	if unicode.IsDigit(firstItem) {
		return "", ErrInvalidString
	}
	// если переданная строка содержит только один символ и это слеш, то вернуть ошибку
	if inSize == 1 && firstItem == 92 {
		return "", ErrInvalidString
	}
	// третий этап - анализ со второго по предпоследний символов
	outThirdStage, err := processThirdStage(inSize, inRunes)
	if err != nil {
		return "", err
	}
	// четвертый этап - анализ последнего символа
	outFourthStage, err := processFourthStage(inSize, inRunes)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	sb.WriteString(outThirdStage)
	sb.WriteString(outFourthStage)

	return sb.String(), nil
}

// выполнение третьего этапа.
func processThirdStage(inSize int, inRunes []rune) (string, error) {
	var sb strings.Builder
	for i := range inSize - 1 {
		sItem := BuildSymBolItem(i, inRunes)
		nextItem := inRunes[i+1]                     // следующий символ
		nextItemIsDigit := unicode.IsDigit(nextItem) // является ли следующий символ цифрой
		// отсекаем ошибку цифр, идущих подряд, при условии, что текущий символ - цифра не экранированая слэшем
		if sItem.Type == IsDigit && !sItem.IsSlashed && nextItemIsDigit {
			return "", ErrInvalidString
		}
		// отсекаем ошибку экранирования символов, не являющихся слешем или цифрой
		if sItem.IsSlashed && sItem.Type == IsOther {
			return "", ErrInvalidString
		}
		// обработка, если текущий символ является цифрой или слешем и при этом экранирован
		if (sItem.Type == IsDigit || sItem.Type == IsSlash) && sItem.IsSlashed {
			if nextItemIsDigit { // если следующий символ некая цифра x, то записать текущий символ x раз
				nextItemInt, err := strconv.Atoi(string(nextItem))
				if err != nil {
					return "", err
				}
				sb.WriteString(strings.Repeat(string(sItem.Item), nextItemInt))
			} else { // если следующий символ не цифра, то записать текущий символ 1 раз
				sb.WriteRune(sItem.Item)
			}
		}
		// обработка, если текущий символ не является цифрой или слешем  и при этом не экранирован
		if sItem.Type == IsOther && !sItem.IsSlashed {
			if nextItemIsDigit { // если следующий символ некая цифра x, то записать текущий символ x раз
				nextItemInt, err := strconv.Atoi(string(nextItem))
				if err != nil {
					return "", err
				}
				for range nextItemInt {
					sb.WriteRune(sItem.Item)
				}
			} else { // если следующий символ не цифра, то записать текущий символ 1 раз
				sb.WriteRune(sItem.Item)
			}
		}
	}

	return sb.String(), nil
}

// выполнение четвертого этапа.
func processFourthStage(inSize int, inRunes []rune) (string, error) {
	var sb strings.Builder
	lastSItem := BuildSymBolItem(inSize-1, inRunes)
	// обработка, если последний символ экранирован
	if lastSItem.IsSlashed {
		if lastSItem.Type == IsDigit || lastSItem.Type == IsSlash { // если последний символ цифра или слеш, то записать его
			sb.WriteRune(lastSItem.Item)
		} else { // если последний символ не является цифрой или слешем, то вернуть ошибку
			return "", ErrInvalidString
		}
	}
	// обработка, если текущий символ не экранирован
	if !lastSItem.IsSlashed {
		if lastSItem.Type == IsOther { // если последний символ не является цифрой или слешем, то записать его
			sb.WriteRune(lastSItem.Item)
		} else if lastSItem.Type == IsSlash { // если последний символ является слешем, то вернуть ошибку
			return "", ErrInvalidString
		}
	}

	return sb.String(), nil
}

func BuildSymBolItem(itemNumber int, inRunes []rune) SymbolItem {
	si := SymbolItem{}
	si.buildSymbolItem(itemNumber, inRunes)
	return si
}

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
}

func (s *SymbolItem) buildSymbolItem(itemNumber int, inRunes []rune) {
	s.Item = inRunes[itemNumber]
	switch {
	case unicode.IsDigit(s.Item):
		s.Type = IsDigit
	case s.Item == SymbolSlash:
		s.Type = IsSlash
	default:
		s.Type = IsOther
	}
	s.IsSlashed = defineIfItemIsSlashed(itemNumber, inRunes) // определение экранирован ли текущий символ
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
