package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(in string) (string, error) {
	// Place your code here.

	var sb strings.Builder
	inRunes := []rune(in)

	inSize := len(inRunes)

	// анализируемая строка содержит 0 символов
	if inSize == 0 {
		return "", nil
	}

	// анализ первого символа
	firstItem := inRunes[0]

	// если первый символ строки содержит цифру, то вернуть ошибку
	if unicode.IsDigit(firstItem) {
		return "", ErrInvalidString
	}

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

	for i := 1; i < inSize-1; i++ {
		previousItem := inRunes[i-1]

		item := inRunes[i]
		itemIsDigit := unicode.IsDigit(item)
		itemIsSlash := (item == 92)

		nextItem := inRunes[i+1]
		nextItemIsDigit := unicode.IsDigit(nextItem)

		// подсчет количества предыдущих символов слеш, следующих подряд
		countPreviousSlash := 0
		for j := i - 1; j >= 0; j-- {
			sItem := inRunes[j]
			sItemIsSlash := (sItem == 92)
			if sItemIsSlash {
				countPreviousSlash = countPreviousSlash + 1
			} else {
				break
			}
		}
		// определение экранирован ли текущий символ
		itemIsSlashed := !(countPreviousSlash%2 == 0)

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

	// запись последнего элемента
	lastItem := inRunes[inSize-1]
	lastItemIsDigit := unicode.IsDigit(lastItem)
	lastItemIsSlash := (lastItem == 92)

	// подсчет количества предыдущих символов слеш, следующих подряд
	countPreviousSlash := 0
	for j := inSize - 2; j >= 0; j-- {
		sItem := inRunes[j]
		sItemIsSlash := (sItem == 92)
		if sItemIsSlash {
			countPreviousSlash = countPreviousSlash + 1
		} else {
			break
		}
	}
	// определение экранирован ли текущий символ
	lastItemIsSlashed := !(countPreviousSlash%2 == 0)

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
	return sb.String(), nil
}
