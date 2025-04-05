package hw02unpackstring

import (
	"errors"
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
	// третий этап - анализ  символов со второго по предпоследний
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
		itemRef := BuildSymbolItem(i, inRunes)                // анализируемый символ
		nextItemRef, err := BuildNextSymBolItem(inRunes[i+1]) // следующий символ
		if err != nil {
			return "", err
		}
		// отсекаем ошибку цифр, идущих подряд, при условии, что текущий символ - цифра не экранированая слэшем
		if itemRef.Type == IsDigit && !itemRef.IsSlashed && nextItemRef.IsDigit {
			return "", ErrInvalidString
		}
		// отсекаем ошибку экранирования символов, не являющихся слешем или цифрой
		if itemRef.IsSlashed && itemRef.Type == IsOther {
			return "", ErrInvalidString
		}
		// обработка, если текущий символ является цифрой или слешем и при этом экранирован
		if (itemRef.Type == IsDigit || itemRef.Type == IsSlash) && itemRef.IsSlashed {
			if nextItemRef.IsDigit { // если следующий символ некая цифра x, то записать текущий символ x раз
				sb.WriteString(strings.Repeat(string(itemRef.Item), nextItemRef.ValueInt))
			} else { // если следующий символ не цифра, то записать текущий символ 1 раз
				sb.WriteRune(itemRef.Item)
			}
		}
		// обработка, если текущий символ не является цифрой или слешем  и при этом не экранирован
		if itemRef.Type == IsOther && !itemRef.IsSlashed {
			if nextItemRef.IsDigit { // если следующий символ некая цифра x, то записать текущий символ x раз
				sb.WriteString(strings.Repeat(string(itemRef.Item), nextItemRef.ValueInt))
			} else { // если следующий символ не цифра, то записать текущий символ 1 раз
				sb.WriteRune(itemRef.Item)
			}
		}
	}

	return sb.String(), nil
}

// выполнение четвертого этапа.
func processFourthStage(inSize int, inRunes []rune) (string, error) {
	var sb strings.Builder
	lastItemRef := BuildSymbolItem(inSize-1, inRunes)
	// обработка, если последний символ экранирован
	if lastItemRef.IsSlashed {
		if lastItemRef.Type == IsDigit || lastItemRef.Type == IsSlash { // если последний символ цифра/слеш, то записать его
			sb.WriteRune(lastItemRef.Item)
		} else { // если последний символ не является цифрой или слешем, то вернуть ошибку
			return "", ErrInvalidString
		}
	}
	// обработка, если текущий символ не экранирован
	if !lastItemRef.IsSlashed {
		if lastItemRef.Type == IsOther { // если последний символ не является цифрой или слешем, то записать его
			sb.WriteRune(lastItemRef.Item)
		} else if lastItemRef.Type == IsSlash { // если последний символ является слешем, то вернуть ошибку
			return "", ErrInvalidString
		}
	}

	return sb.String(), nil
}
