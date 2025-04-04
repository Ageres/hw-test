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
	var sb strings.Builder

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

	// является ли первый символ спешем
	firstItemIsSlash := (firstItem == 92)

	// обработка , если анализируемая строка содержит 1 символ
	if inSize == 1 {
		if !firstItemIsSlash { // если символ не слеш, то вернуть его в виде строки
			return string(firstItem), nil
		}
		return "", nil // иначе вернуть пустую строку
	}

	// если первый символ не слеш
	if !firstItemIsSlash {
		// то проводится анализ второго символа
		secondItem := inRunes[1]
		if unicode.IsDigit(secondItem) { // если второй символ некая цифра x, то записать первый символ x раз
			secondItemInt, err := strconv.Atoi(string(secondItem))
			if err != nil {
				return "", err
			}
			sb.WriteString(strings.Repeat(string(firstItem), secondItemInt))
		} else { // иначе записать первый символ 1 раз
			sb.WriteRune(firstItem)
		}
	}

	//--------------------------------
	// Третий этап
	outTS, err := processThirdStage(inSize, inRunes)
	if err != nil {
		return "", err
	}
	sb.WriteString(outTS)

	//--------------------------------
	// Четвертый этап
	outFS := processFourthStage(inSize, inRunes)
	sb.WriteString(outFS)

	return sb.String(), nil
}

// выполнение третьего этапа.
func processThirdStage(inSize int, inRunes []rune) (string, error) {
	var sb strings.Builder
	for i := 1; i < inSize-1; i++ {
		previousItem := inRunes[i-1] // предыдущий символ

		item := inRunes[i]                                 // текущий анализируемый символ
		itemIsDigit := unicode.IsDigit(item)               // является ли текущий элемент цифрой
		itemIsSlash := (item == 92)                        // является ли текущий элемент слешем
		itemIsOther := !itemIsDigit && !itemIsSlash        // является ли текущий элемент прочим символом
		itemIsSlashed := defineIfItemIsSlashed(i, inRunes) // определение экранирован ли текущий символ

		nextItem := inRunes[i+1] // следующий символ
		nextItemIsDigit := unicode.IsDigit(nextItem)

		// отсекаем ошибку цифр, идущих подряд, при условии, что текущий символ - цифра не экранированая слэшем
		if !itemIsSlashed && itemIsDigit && nextItemIsDigit {
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
				previousItem,
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
	previousItem, item, nextItem rune,
	itemIsSlashed, nextItemIsDigit bool,
) (string, error) {
	var sb strings.Builder
	// обработка, если текущий символ экранирован
	if itemIsSlashed {
		if nextItemIsDigit { // если следующий символ некая цифра x, то записать комбинацию из слеша и текущего символа x раз
			nextItemInt, err := strconv.Atoi(string(nextItem))
			if err != nil {
				return "", err
			}
			for range nextItemInt {
				sb.WriteRune(previousItem)
				sb.WriteRune(item)
			}
		} else { // если следующий символ не цифра, то записать комбинацию из слеша и текущего символа 1 раз
			sb.WriteRune(previousItem)
			sb.WriteRune(item)
		}
	}
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
func processFourthStage(inSize int, inRunes []rune) string {
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
			sb.WriteString("\\")
			sb.WriteRune(lastItem)
		}
	} else {
		if !lastItemIsDigit && !lastItemIsSlash {
			sb.WriteRune(lastItem)
		}
	}
	return sb.String()
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
