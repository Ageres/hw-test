package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")
var SymbolSlash rune = []rune(`\`)[0]

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
		item := BuildSymBolItem(i, inRunes)

		//item := inRunes[i]                                 // текущий анализируемый символ
		//itemIsDigit := unicode.IsDigit(item)               // является ли текущий элемент цифрой
		//itemIsSlash := (item == 92)                        // является ли текущий элемент слешем
		//itemIsOther := !itemIsDigit && !itemIsSlash        // является ли текущий элемент прочим символом
		//itemIsSlashed := defineIfItemIsSlashed(i, inRunes) // определение экранирован ли текущий символ
		nextItem := inRunes[i+1]                     // следующий символ
		nextItemIsDigit := unicode.IsDigit(nextItem) // является ли следующий символ цифрой
		// отсекаем ошибку цифр, идущих подряд, при условии, что текущий символ - цифра не экранированая слэшем
		if item.Type == IsDigit && !item.IsSlashed && nextItemIsDigit {
			return "", ErrInvalidString
		}
		// отсекаем ошибку экранирования символов, не являющихся слешем или цифрой
		if item.IsSlashed && item.Type == IsOther {
			return "", ErrInvalidString
		}
		// обработка, если текущий символ является цифрой или слешем и при этом экранирован
		if (item.Type == IsDigit || item.Type == IsSlash) && item.IsSlashed {
			if nextItemIsDigit { // если следующий символ некая цифра x, то записать текущий символ x раз
				nextItemInt, err := strconv.Atoi(string(nextItem))
				if err != nil {
					return "", err
				}
				sb.WriteString(strings.Repeat(string(item.Item), nextItemInt))
			} else { // если следующий символ не цифра, то записать текущий символ 1 раз
				sb.WriteRune(item.Item)
			}
		}
		// обработка, если текущий символ не является цифрой или слешем  и при этом не экранирован
		if item.Type == IsOther && !item.IsSlashed {
			if nextItemIsDigit { // если следующий символ некая цифра x, то записать текущий символ x раз
				nextItemInt, err := strconv.Atoi(string(nextItem))
				if err != nil {
					return "", err
				}
				for range nextItemInt {
					sb.WriteRune(item.Item)
				}
			} else { // если следующий символ не цифра, то записать текущий символ 1 раз
				sb.WriteRune(item.Item)
			}
		}
	}

	return sb.String(), nil
}

// выполнение четвертого этапа.
func processFourthStage(inSize int, inRunes []rune) (string, error) {
	var sb strings.Builder
	lastItem := inRunes[inSize-1]
	lastItemIsDigit := unicode.IsDigit(lastItem)
	lastItemIsSlash := (lastItem == 92)
	lastItemIsSlashed := defineIfItemIsSlashed(inSize-1, inRunes) // определение экранирован ли последний символ
	// обработка, если последний символ экранирован
	if lastItemIsSlashed {
		if lastItemIsDigit || lastItemIsSlash { // если последний символ является цифрой или слешем, то записать его
			sb.WriteRune(lastItem)
		} else { // если последний символ не является цифрой или слешем, то вернуть ошибку
			return "", ErrInvalidString
		}
	}
	// обработка, если текущий символ не экранирован
	if !lastItemIsSlashed {
		if !lastItemIsDigit && !lastItemIsSlash { // если последний символ не является цифрой или слешем, то записать его
			sb.WriteRune(lastItem)
		} else if lastItemIsSlash { // если последний символ является слешем, то вернуть ошибку
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

// перечисление с типами анализируемого символа
type Type int

const (
	IsDigit Type = iota + 1 // анализируемый символ это цифра
	IsSlash                 // анализируемый символ это слеш
	IsOther                 // анализируемый символ это не цифра и не слеш
)

// структура с параметрами анализируемого символа
type SymbolItem struct {
	Item      rune // анализируемый символ
	Type      Type // тип анализируемого символа
	IsSlashed bool // анализируемый символ экранирован
}

func (s *SymbolItem) buildSymbolItem(itemNumber int, inRunes []rune) {
	s.Item = inRunes[itemNumber]
	if unicode.IsDigit(s.Item) {
		s.Type = IsDigit
	} else if s.Item == 92 {
		s.Type = IsSlash
	} else {
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
		sItemIsSlash := (sItem == 92)
		if sItemIsSlash {
			countPreviousSlash++
		} else {
			break
		}
	}

	return !(countPreviousSlash%2 == 0)
}
