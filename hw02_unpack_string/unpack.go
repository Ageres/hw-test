package hw02unpackstring

import (
	"errors"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

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
	firstItemRef := BuildSymbolItem(0, inRunes)

	if firstItemRef.Type == IsDigit {
		return "", ErrInvalidString
	}
	// если переданная строка содержит только один символ
	if inSize == 1 {
		if firstItemRef.Type == IsSlash { // если передан символ слеш, то вернуть ошибку
			return "", ErrInvalidString
		}
		return string(firstItemRef.Item), nil // если передан прочий символ, то вернуть строку из одного переданного символа
	}

	// третий этап - анализ  символов со второго по предпоследний
	// слайс для хранения объектов анализируемых символов
	items := make([]SymbolItem, inSize)
	items[0] = *firstItemRef

	outThirdStage, err := processThirdStage(inSize, inRunes, items)
	if err != nil {
		return "", err
	}

	// четвертый этап - анализ последнего символа
	lastIiem := items[inSize-1]
	outFourthStage, err := processFourthStage(&lastIiem)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	sb.WriteString(outThirdStage)
	sb.WriteString(outFourthStage)

	return sb.String(), nil
}

// выполнение третьего этапа.
func processThirdStage(inSize int, inRunes []rune, items []SymbolItem) (string, error) {
	var sb strings.Builder
	for i := range inSize - 1 {
		itemRef := items[i]                          // анализируемый символ
		nextItemRef := BuildSymbolItem(i+1, inRunes) // следующий символ
		if err := nextItemRef.ParseIfDigit(); err != nil {
			return "", err
		}
		items[i+1] = *nextItemRef
		// отсекаем ошибку цифр, идущих подряд, при условии, что текущий символ - цифра не экранированая слэшем
		if itemRef.Type == IsDigit && !itemRef.IsSlashed && nextItemRef.Type == IsDigit {
			return "", ErrInvalidString
		}
		// отсекаем ошибку экранирования символов, не являющихся слешем или цифрой
		if itemRef.IsSlashed && itemRef.Type == IsOther {
			return "", ErrInvalidString
		}
		// обработка, если текущий символ является цифрой или слешем и при этом экранирован
		if (itemRef.Type == IsDigit || itemRef.Type == IsSlash) && itemRef.IsSlashed {
			sb.WriteString(itemRef.RepeatWith(nextItemRef))
		}
		// обработка, если текущий символ не является цифрой или слешем  и при этом не экранирован
		if itemRef.Type == IsOther && !itemRef.IsSlashed {
			sb.WriteString(itemRef.RepeatWith(nextItemRef))
		}
	}

	return sb.String(), nil
}

// выполнение четвертого этапа.
func processFourthStage(lastItemRef *SymbolItem) (string, error) {
	var sb strings.Builder
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
