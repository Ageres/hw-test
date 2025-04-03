package hw02unpackstring

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(in string) (string, error) {
	// Place your code here.

	fmt.Println("-------------------------------------------------------")

	fmt.Println("in:    ", in)

	var sb strings.Builder
	inRunes := []rune(in)

	//inSize := utf8.RuneCountInString(in)
	inSize := len(inRunes)
	fmt.Println("inSize:", inSize)

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

		fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", i, ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")

		previousItem := inRunes[i-1]
		//previousItemIsDigit := unicode.IsDigit(previousItem)
		//previousItemIsSlash := (previousItem == 92)

		item := inRunes[i]
		itemIsDigit := unicode.IsDigit(item)
		itemIsSlash := (item == 92)

		nextItem := inRunes[i+1]
		nextItemIsDigit := unicode.IsDigit(nextItem)
		//nextItemIsSlash := (nextItem == 92)

		//--------------------------------------------------------------------

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

		//--------------------------------------------------------------------

		/*
			fmt.Println("++++++++++++++++++++++++++++")
			fmt.Println("previousItem:        ", string(previousItem))
			fmt.Println("previousItemIsDigit: ", previousItemIsDigit)
			fmt.Println("previousItemIsSlash: ", previousItemIsSlash)
		*/

		fmt.Println("++++++++++++++++++++++++++++")
		fmt.Println("item:                ", string(item))
		fmt.Println("itemIsDigit:         ", itemIsDigit)
		fmt.Println("itemIsSlash:         ", itemIsSlash)
		//fmt.Println("countPreviousSlash:  ", countPreviousSlash)
		fmt.Println("itemIsSlashed:       ", itemIsSlashed)
		fmt.Println("++++++++++++++++++++++++++++")

		/*
			fmt.Println("nextItem:            ", string(nextItem))
			fmt.Println("nextItemIsDigit:     ", nextItemIsDigit)
			fmt.Println("nextItemIsSlash:     ", nextItemIsSlash)
			fmt.Println("++++++++++++++++++++++++++++")
		*/

		// отсекаем ошибку цифр, идущих подряд, при условии, что первая цифра не экранирована слэшем
		if !itemIsSlashed && itemIsDigit && nextItemIsDigit {
			return "", ErrInvalidString
		}

		if itemIsDigit || itemIsSlash {
			if itemIsSlashed {
				if nextItemIsDigit {
					nextItemInt, err := strconv.Atoi(string(nextItem))
					fmt.Println("-----------001---------", nextItemInt)
					fmt.Println("previousItem:        ", string(previousItem))
					fmt.Println("item:                ", string(item))
					if err != nil {
						return "", err
					}
					sb.WriteString(strings.Repeat(string(item), nextItemInt))
					/*
						for range nextItemInt {
							//sb.WriteRune(previousItem)
							sb.WriteRune(item)
						}
					*/
				} else {
					//sb.WriteRune(previousItem)
					sb.WriteRune(item)
				}
			}
		}

		if !itemIsDigit && !itemIsSlash {
			if itemIsSlashed {
				if nextItemIsDigit {
					nextItemInt, err := strconv.Atoi(string(nextItem))
					fmt.Println("-----------002---------", nextItemInt)
					fmt.Println("previousItem:        ", string(previousItem))
					fmt.Println("item:                ", string(item))
					if err != nil {
						return "", err
					}
					//sb.WriteString(strings.Repeat(string(item), nextItemInt))

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
					fmt.Println("-----------003---------", nextItemInt)
					fmt.Println("previousItem:        ", string(previousItem))
					fmt.Println("item:                ", string(item))
					if err != nil {
						return "", err
					}
					//sb.WriteString(strings.Repeat(string(item), nextItemInt))

					for range nextItemInt {
						sb.WriteRune(item)
					}

				} else {
					sb.WriteRune(item)
				}
			}
		}

		fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", i, "<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
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
