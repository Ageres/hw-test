package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
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

	//--------------------------------
	// Первый этап
	// анализируемая строка содержит 0 символов

	if inSize == 0 {
		return "", nil
	}

	//--------------------------------
	// Второй этап
	// анализ первого символа
	firstItem := inRunes[0]

	// если первый символ строки содержит цифру, то вернуть ошибку
	if unicode.IsDigit(firstItem) {
		return "", ErrInvalidString
	}

	// если переданная строка содержит только один символ и это слеш
	// то вернуть ошибку
	if inSize == 1 && firstItem == 92 {
		return "", ErrInvalidString
	}

	//--------------------------------
	// Третий этап
	// анализ со второго по предпоследний символов
	var sb strings.Builder

	outTS, err := processThirdStage(inSize, inRunes)
	if err != nil {
		return "", err
	}
	sb.WriteString(outTS)

	//--------------------------------
	// Четвертый этап
	// анализ последнего символа
	outFS, err := processFourthStage(inSize, inRunes)
	if err != nil {
		return "", err
	}
	sb.WriteString(outFS)

	return sb.String(), nil
}

// выполнение третьего этапа.
func processThirdStage(inSize int, inRunes []rune) (string, error) {
	var sb strings.Builder
	for i := 0; i < inSize-1; i++ {
		item := inRunes[i]                                 // текущий анализируемый символ
		itemIsDigit := unicode.IsDigit(item)               // является ли текущий элемент цифрой
		itemIsSlash := (item == 92)                        // является ли текущий элемент слешем
		itemIsOther := !itemIsDigit && !itemIsSlash        // является ли текущий элемент прочим символом
		itemIsSlashed := defineIfItemIsSlashed(i, inRunes) // определение экранирован ли текущий символ

		nextItem := inRunes[i+1]                     // следующий символ
		nextItemIsDigit := unicode.IsDigit(nextItem) // является ли следующий символ цифрой

		// отсекаем ошибку цифр, идущих подряд, при условии, что текущий символ - цифра не экранированая слэшем
		if !itemIsSlashed && itemIsDigit && nextItemIsDigit {
			return "", ErrInvalidString
		}

		// отсекаем ошибку экранирования символов, не являющихся слешем или цифрой
		if itemIsSlashed && itemIsOther {
			return "", ErrInvalidString
		}

		// обработка, если текущий символ является цифрой или слешем и при этом экранирован
		if (itemIsDigit || itemIsSlash) && itemIsSlashed {
			if nextItemIsDigit { // если следующий символ некая цифра x, то записать текущий символ x раз
				nextItemInt, err := strconv.Atoi(string(nextItem))
				if err != nil {
					return "", err
				}
				sb.WriteString(strings.Repeat(string(item), nextItemInt))
			} else { // если следующий символ не цифра, то записать текущий символ 1 раз
				sb.WriteRune(item)
			}
		}

		// обработка, если текущий символ не является цифрой или слешем
		if itemIsOther {
			outTSMFOS, err := processThirdStageModuleForOtherSymbolType(
				item,
				nextItem,
				itemIsSlashed,
				nextItemIsDigit,
			)
			if err != nil {
				return "", err
			}
			sb.WriteString(outTSMFOS)
		}
	}
	return sb.String(), nil
}

// обработка, если текущий символ не является цифрой или слешем.
func processThirdStageModuleForOtherSymbolType(
	item, nextItem rune,
	itemIsSlashed, nextItemIsDigit bool,
) (string, error) {
	var sb strings.Builder
	// обработка, если текущий символ не экранирован
	if !itemIsSlashed {
		if nextItemIsDigit { // если следующий символ некая цифра x, то записать текущий символ x раз
			nextItemInt, err := strconv.Atoi(string(nextItem))
			if err != nil {
				return "", err
			}
			for range nextItemInt {
				sb.WriteRune(item)
			}
		} else { // если следующий символ не цифра, то записать текущий символ 1 раз
			sb.WriteRune(item)
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

	// определение экранирован ли последний символ
	lastItemIsSlashed := defineIfItemIsSlashed(inSize-1, inRunes)

	if lastItemIsSlashed {
		if lastItemIsDigit || lastItemIsSlash {
			sb.WriteRune(lastItem)
		} else {
			return "", ErrInvalidString
		}
	}

	if !lastItemIsSlashed {
		if !lastItemIsDigit && !lastItemIsSlash {
			sb.WriteRune(lastItem)
		} else if lastItemIsSlash {
			return "", ErrInvalidString
		}
	}
	return sb.String(), nil
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
