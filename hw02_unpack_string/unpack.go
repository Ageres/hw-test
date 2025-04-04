package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

/*
 * Распаковка в 4 этапа:
 * Первый этап    - проверяется есть ли символы в переданной строке, если нет - возвращается пустая строка;
 * Второй этап    - проверяется достаточность условий и формат распаковки первого символа на основании видов первого
 *                  и второго символов (число/слеш/прочий символ);
 * Третий этап    - проверяется достаточность условий и формат распаковки от второго до предпоследнего символа;
 *	                решение для анализируемого символа принимается на основании собственного вида, а так же видов
 *                  предыдущих и последующего символов;
 *				    при необходимости первый символ так же добавляется в распаковку;
 * Четвертый этап - проверяется достаточность условий и формат распаковки последнего символа на основании собственного
 *                  вида, а так же видов предыдущих символов.
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

	// анализируемая строка содержит 1 символ
	if inSize == 1 {
		if !firstItemIsSlash {
			return string(firstItem), nil
		}
		return "", nil
	}

	// если первый символ не слеш
	// то проводится анализ второго символа - если это цифра x, то записать первый символ x раз
	// иначе 1 раз
	if !firstItemIsSlash {
		secondItem := inRunes[1]
		if unicode.IsDigit(secondItem) {
			secondItemInt, err := strconv.Atoi(string(secondItem))
			if err != nil {
				return "", err
			}
			sb.WriteString(strings.Repeat(string(firstItem), secondItemInt))
		} else {
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

// выполнение третьего этапа
func processThirdStage(inSize int, inRunes []rune) (string, error) {
	var sb strings.Builder
	for i := 1; i < inSize-1; i++ {
		previousItem := inRunes[i-1]

		item := inRunes[i]
		itemIsDigit := unicode.IsDigit(item)               // является ли текущий элемент цифрой
		itemIsSlash := (item == 92)                        // является ли текущий элемент слешем
		itemIsSlashed := defineIfItemIsSlashed(i, inRunes) // определение экранирован ли текущий символ

		nextItem := inRunes[i+1]
		nextItemIsDigit := unicode.IsDigit(nextItem)

		// отсекаем ошибку цифр, идущих подряд, при условии, что первая цифра не экранирована слэшем
		if !itemIsSlashed && itemIsDigit && nextItemIsDigit {
			return "", ErrInvalidString
		}

		if itemIsDigit || itemIsSlash {
			if itemIsSlashed {
				if nextItemIsDigit {
					nextItemInt, err := strconv.Atoi(string(nextItem))
					if err != nil {
						return "", err
					}
					sb.WriteString(strings.Repeat(string(item), nextItemInt))
				} else {
					sb.WriteRune(item)
				}
			}
		}

		if !itemIsDigit && !itemIsSlash {
			if itemIsSlashed {
				if nextItemIsDigit {
					nextItemInt, err := strconv.Atoi(string(nextItem))
					if err != nil {
						return "", err
					}
					for range nextItemInt {
						sb.WriteRune(previousItem)
						sb.WriteRune(item)
					}

				} else {
					sb.WriteRune(previousItem)
					sb.WriteRune(item)
				}
			} else {
				if nextItemIsDigit {
					nextItemInt, err := strconv.Atoi(string(nextItem))
					if err != nil {
						return "", err
					}
					for range nextItemInt {
						sb.WriteRune(item)
					}
				} else {
					sb.WriteRune(item)
				}
			}
		}
	}
	return sb.String(), nil
}

func processFourthStage(inSize int, inRunes []rune) string {
	var sb strings.Builder
	// запись последнего элемента
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

// определение экранирован ли символ
func defineIfItemIsSlashed(itemNumber int, inRunes []rune) bool {
	// подсчет количества предыдущих символов слеш, следующих подряд
	countPreviousSlash := 0
	for j := itemNumber - 1; j >= 0; j-- {
		sItem := inRunes[j]
		sItemIsSlash := (sItem == 92)
		if sItemIsSlash {
			countPreviousSlash = countPreviousSlash + 1
		} else {
			break
		}
	}

	return !(countPreviousSlash%2 == 0)
}
