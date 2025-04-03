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
		fmt.Println("------1------")
		secondItem := inRunes[1]
		if unicode.IsDigit(secondItem) {
			fmt.Println("------2------")
			secondItemInt, err := strconv.Atoi(string(secondItem))
			if err != nil {
				return "", err
			}
			sb.WriteString(strings.Repeat(string(firstItem), secondItemInt))
		} else {
			sb.WriteRune(firstItem)
		}
	}

	/*
		var prePreItem rune = inRunes[0] // предпредпоследний символ (i - 2)
		var prePreItemIsDigit bool = unicode.IsDigit(prePreItem)
		var prePreItemIsSlash bool = prePreItem == 92
		var prePreItemIsWritten bool = false

		var preItem rune = inRunes[1] // предпоследний символ (i - 1)
		var preItemIsDigit bool = unicode.IsDigit(preItem)
		var preItemIsSlash bool = preItem == 92
		var preItemIsWritten bool = false

		//outRune := make([]rune, inSize)

		//sbLen := 0
	*/
	for i := 1; i < inSize-1; i++ {

		fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", i, ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")

		previousItem := inRunes[i]
		previousItemIsDigit := unicode.IsDigit(previousItem)
		previousItemIsSlash := (previousItem == 92)

		item := inRunes[i]
		itemIsDigit := unicode.IsDigit(item)
		itemIsSlash := (item == 92)

		nextItem := inRunes[i+1]
		nextItemIsDigit := unicode.IsDigit(nextItem)
		nextItemIsSlash := (nextItem == 92)

		fmt.Println("++++++++++++++++++++++++++++")

		fmt.Println("previousItem:        ", string(previousItem))
		fmt.Println("previousItemIsDigit: ", previousItemIsDigit)
		fmt.Println("previousItemIsSlash: ", previousItemIsSlash)

		fmt.Println("++++++++++++++++++++++++++++")

		fmt.Println("item:                ", string(item))
		fmt.Println("itemIsDigit:         ", itemIsDigit)
		fmt.Println("itemIsSlash:         ", itemIsSlash)

		fmt.Println("++++++++++++++++++++++++++++")

		fmt.Println("nextItem:            ", string(nextItem))
		fmt.Println("nextItemIsDigit:     ", nextItemIsDigit)
		fmt.Println("nextItemIsSlash:     ", nextItemIsSlash)

		fmt.Println("++++++++++++++++++++++++++++")

		// отсекаем ошибку цифр, идущих подряд, при условии, что первая цифра не экранирована слэшем
		if !previousItemIsSlash && itemIsDigit && nextItemIsDigit {
			return "", ErrInvalidString
		}

		/*
			if itemIsDigit && !previousItemIsSlash {
				if itemNextIsDigit {
					return "", ErrInvalidString
				} else {
					continue
				}
			}
		*/

		if itemIsDigit {
			// предыдущий символ может быть только слэш или не слеш и не цифра
			if previousItemIsSlash {
				if nextItemIsDigit {
					nextItemInt, err := strconv.Atoi(string(nextItem))
					if err != nil {
						return "", err
					}
					sb.WriteString(strings.Repeat(string(previousItem)+string(item), nextItemInt))
				} else if nextItemIsSlash {
					sb.WriteString(string(previousItem) + string(item))
				} else {
					sb.WriteString(string(previousItem) + string(item))
				}
			} else {
				continue
			}
		} else if itemIsSlash {
			// предыдущий символ может быть только слэш или не слеш и не цифра
			if previousItemIsSlash {

			} else {
				continue
			}
		} else {
			// предыдущий символ может быть только слэш или не слеш и не цифра
			if previousItemIsSlash {

			} else {

			}

		}

		// анализ предыдущего символа - если он слеш - то определить, нужно ли экранировать.
		// экранируется только цифра и слеш

		if itemIsSlash {
			continue
		}

		// остается item который и не слеш и не цифра
		if nextItemIsDigit {
			// записать в sb в определеное количество раз
			continue
		} else {
			// записать в sb один раз
		}

		/*
			itemIsFrst := i ==0
			if itemIsFrst {
				if itemIsDigit {
					return "", ErrInvalidString
				} else if !itemIsSlash {
					if i != inSize - 1 {
						itemNext := inRunes[i + 1]
						itemNextIsDigit := unicode.IsDigit(itemNext)
						if itemNextIsDigit {

						}
					} else {
						sb.WriteRune(item)
					}
					continue
				}
			}
		*/
		/*
			item := inRunes[i]
			if i > 0 {
			}
			previousItem := inRunes[i-1]
			item := inRunes[i]
			nextItem := inRunes[i+1]

			itemIsDigit := unicode.IsDigit(item)
			itemIsSlash := (item == 92)
			if !itemIsDigit && !itemIsSlash {
				sb.WriteRune(item)
				preItemIsWritten = true
			}
		*/
		fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", i, "<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
	}

	/*
		for item := range inRunes {
			//получить последний элемент из sb
			itemIsDigit := unicode.IsDigit(item)
			itemIsSlash := (item == 92)
		}
	*/

	/*
			for i := 2; i < inSize; i++ {
				fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", i, ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
				item := inRunes[i]
				itemIsDigit := unicode.IsDigit(item)
				itemIsSlash := (item == 92)

				//---------------------------
				fmt.Println("++++++++++++++++++++++++++++")

				fmt.Println("prePreItem: ", string(prePreItem))
				fmt.Println("prePreItemIsDigit: ", prePreItemIsDigit)
				fmt.Println("prePreItemIsSlash: ", prePreItemIsSlash)
				fmt.Println("prePreItemIsWritten:", prePreItemIsWritten)

				fmt.Println("++++++++++++++++++++++++++++")

				fmt.Println("preItem: ", string(preItem))
				fmt.Println("preItemIsDigit: ", preItemIsDigit)
				fmt.Println("preItemIsSlash: ", preItemIsSlash)
				fmt.Println("preItemIsWritten:", preItemIsWritten)

				fmt.Println("++++++++++++++++++++++++++++")

				fmt.Println("item: ", string(item))
				fmt.Println("itemRune: ", item)
				fmt.Println("itemIsDigit:", itemIsDigit)
				fmt.Println("itemIsSlash:", itemIsSlash)

				fmt.Println("++++++++++++++++++++++++++++")

				//---------------------------------------------------

				// если предпоследний символ цифра и последний символ цифра или
				// если предпредпоследний символ цифра и предпоследний символ цифра
				// то вернуть ошибку
				if (preItemIsDigit && itemIsDigit) || (prePreItemIsDigit && preItemIsDigit) {
					return "", ErrInvalidString
				}

				ifIsUsed := false // использовано хотя бы одно условие для записи

				}

				// если предпредпоследний символ не записан и если он слеш и если предпоследний символ не цифра и если последний символ некая цифра y (itemInt),
				// то записать комбинацию из препредпоследнего и предпоследнего символа y раз,
				// пометить предпредпоследний и предпоследний символ как записанные (с учетом того, что в следующей итерации предпоследний символ станет предпредпоследним
				// а последний символ станет предпоследним)
				if !prePreItemIsWritten && prePreItemIsSlash && !preItemIsDigit && itemIsDigit {
					ifIsUsed = true
					fmt.Println("------------ if 1")
					itemInt, err := strconv.Atoi(string(item))
					if err != nil {
						return "", err
					}
					sb.WriteString(strings.Repeat(string(prePreItem)+string(preItem), itemInt))
					prePreItemIsWritten = true
					preItemIsWritten = true
				}

				// если предпредпоследний символ не записан и если он не слеш и
				// если предпоследний символ не цифра и если он не слеш и
				// если последний символ не цифра,
				// то записать комбинацию из препредпоследнего и предпоследнего символа 1 раз,
				// пометить предпредпоследний и предпоследний символ как записанные (с учетом того, что в следующей итерации предпоследний символ станет предпредпоследним
				// а последний символ станет предпоследним)
				if !prePreItemIsWritten && !prePreItemIsSlash && !preItemIsDigit && !preItemIsSlash && !itemIsDigit {
					ifIsUsed = true
					fmt.Println("------------ if 1.1")
					sb.WriteRune(prePreItem)
					sb.WriteRune(preItem)
					prePreItemIsWritten = true
					preItemIsWritten = true
				}

				// если предпредпоследний символ не записан и если он не слеш
				// если предпоследний символ некая цифра x
				// то записать предпредпоследний символ x раз
				// пометить предпредпоследний символ как записанный и предпоследний символ как не записанный
				if !prePreItemIsWritten && !prePreItemIsSlash && preItemIsDigit {
					ifIsUsed = true
					fmt.Println("------------ if 2")
					preItemInt, err := strconv.Atoi(string(preItem))
					if err != nil {
						return "", err
					}
					sb.WriteString(strings.Repeat(string(prePreItem), preItemInt))
					prePreItemIsWritten = true
					preItemIsWritten = false
				}

				// если предпредпоследний символ не записан и если он не слеш и не цифра
				// если предпоследний символ (не записан и если он) не цифра и не слеш
				// то записать предпредпоследний символ 1 раз
				// пометить предпредпоследний символ как записанный и предпоследний символ как не записанный
				if !prePreItemIsWritten && !prePreItemIsSlash && !preItemIsDigit && !preItemIsSlash {
					ifIsUsed = true
					fmt.Println("------------ if 2.1")
					sb.WriteRune(prePreItem)
					prePreItemIsWritten = true
					preItemIsWritten = false
				}

				// если предпредпоследний символ записан и если предпоследний символ не записан и если предпоследний символ не цифра и если последний символ не цифра
				// то записать предпредпоследний символ 1 раз
				// пометить предпредпоследний символ как записанный и предпоследний символ как не записанный
				if prePreItemIsWritten && !preItemIsWritten && !preItemIsDigit && !itemIsDigit {
					ifIsUsed = true
					fmt.Println("------------ if 3")
					sb.WriteRune(preItem)
					prePreItemIsWritten = true
					preItemIsWritten = false
				}

				//----------------------------------------------

				if !ifIsUsed { // не использовано ни одно условие для записи
					prePreItemIsWritten = false
					preItemIsWritten = false
				}

				prePreItem = preItem
				prePreItemIsDigit = preItemIsDigit
				prePreItemIsSlash = preItemIsSlash

				preItem = item
				preItemIsDigit = itemIsDigit
				preItemIsSlash = itemIsSlash

				fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", i, "<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
			}

			// запись последнего элемента

			lastItem := inRunes[inSize-1]
			lastItemIsDigit := unicode.IsDigit(lastItem)
			fmt.Println("lastItemIsDigit:", lastItemIsDigit)

			fmt.Println("sb:", sb.String())

			if lastItemIsDigit {
				preLastitem := inRunes[inSize-2]
				lastItemInt, err := strconv.Atoi(string(lastItem))
				if err != nil {
					return "", err
				}
				sb.WriteString(strings.Repeat(string(preLastitem), lastItemInt))
			} else {
				sb.WriteRune(inRunes[inSize-1])
			}

			return sb.String(), nil
		}

		func stringSize2(in string) (string, error) {
			fmt.Println("----------------- size 2")

			inRunes := []rune(in)
			var sb strings.Builder
			in0 := inRunes[0]
			//fmt.Println("in[0]:", in[0])
			//fmt.Println("in0:", string(in0))
			in0IsDigit := unicode.IsDigit(in0)
			in1 := inRunes[1]
			//fmt.Println("in[1]:", in[1])
			//fmt.Println("in1:", string(in1))
			in1IsDigit := unicode.IsDigit(in1)

			if in0IsDigit {
				return "", ErrInvalidString
			} else {
				if in1IsDigit {
					in1Int, err := strconv.Atoi(string(in1))
					if err != nil {
						return "", err
					}
					sb.WriteString(strings.Repeat(string(in0), in1Int))
				} else {
					sb.WriteRune(in0)
					sb.WriteRune(in1)
				}
			}
	*/
	return sb.String(), nil
}

func stringSize1(in string) string {
	inRunes := []rune(in)
	in0 := inRunes[0]
	in0IsSlash := (in0 == 92)
	if in0IsSlash {
		return ""
	} else {
		return string(in0)
	}
}
