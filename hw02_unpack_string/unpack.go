package hw02unpackstring

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(in string) (string, error) {
	// Place your code here.

	fmt.Println("-------------------------------------------------------")

	inSize := utf8.RuneCountInString(in)
	fmt.Println("inSize:", inSize)

	// анализируемая строка содержит 0 символов
	if inSize == 0 {
		return "", nil
	}

	// если первый символ строки содержит цифру, то вернуть ошибку
	if unicode.IsDigit(rune(in[0])) {
		return "", ErrInvalidString
	}

	// анализируемая строка содержит 1 символ
	if inSize == 1 {
		return stringSize1(in)
	}

	// анализируемая строка содержит 2 символа
	if inSize == 2 {
		return stringSize2(in)
	}

	var sb strings.Builder

	var prePreItem rune = rune(in[0]) // предпредпоследний символ (i - 2)
	var prePreItemIsDigit bool = unicode.IsDigit(rune(in[0]))
	var prePreItemIsSlash bool = rune(in[0]) == 92
	var prePreItemIsWritten bool = false

	var preItem rune = rune(in[1]) // предпоследний символ (i - 1)
	var preItemIsDigit bool = unicode.IsDigit(rune(in[1]))
	var preItemIsSlash bool = rune(in[1]) == 92
	var preItemIsWritten bool = false

	for i := 2; i < inSize; i++ {
		fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", i, ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
		item := rune(in[i])
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

		// если предпредпоследний символ не записан и если предпоследний символ некая цифра x
		// то записать предпредпоследний символ x раз
		// пометить предпредпоследний символ как записанный и предпоследний символ как не записанный
		if !prePreItemIsWritten && preItemIsDigit {
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

		/*
			if !prePreItemIsWritten {
				if preItemIsDigit {
					if itemIsDigit {
						return "", ErrInvalidString
					} else {
						preItemInt, err := strconv.Atoi(string(item))
						if err != nil {
							return "", err
						}
						sb.WriteString(strings.Repeat(string(prePreItem), preItemInt))
					}
				}
			}
		*/

		/*
			if prePreItemIsWritten {
				if preItemIsWritten {

				} else {
					if prePreItemIsDigit {

					} else {

					}
				}
			} else {

			}*/

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
	sb.WriteRune(rune(in[inSize-1]))

	return sb.String(), nil
}

/*
		inSize := utf8.RuneCountInString(in)
		fmt.Println("inSize:", inSize)

		// анализируемая строка содержит 0 символов
		if inSize == 0 {
			return "", nil
		}

		var sb strings.Builder
		isWritePreItem := false

		in0 := rune(in[0])
		in0IsDigit := unicode.IsDigit(in0)
		if in0IsDigit {
			return "", ErrInvalidString
		} else {
			sb.WriteRune(in0)
		}

		// анализируемая строка содержит 1 символ
		if inSize == 1 {
			return sb.String(), nil
		}

		in1 := rune(in[1])
		in1IsDigit := unicode.IsDigit(in1)

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

		// анализируемая строка содержит 2 символа
		if inSize == 2 {
			return sb.String(), nil
		}

		// анализируемая строка содержит больше 2-х символов
		var prePreItem rune = rune(in[0])
		var prePreItemIsDigit bool = unicode.IsDigit(rune(in[0]))
		var prePreItemIsSlash bool = rune(in[0]) == 92

		var preItem rune = rune(in[1])
		var preItemIsDigit bool = unicode.IsDigit(rune(in[1]))
		var preItemIsSlash bool = rune(in[1]) == 92

		isWritePreItem := true

		for i := 2; i < inSize; i++ {
			fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", i, ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
			item := rune(in[i])
			itemIsDigit := unicode.IsDigit(item)
			itemIsSlash := (item == 92)

			//---------------------------

			fmt.Println("prePreItem: ", string(prePreItem))
			fmt.Println("prePreItemIsDigit: ", prePreItemIsDigit)
			fmt.Println("prePreItemIsSlash: ", prePreItemIsSlash)

			fmt.Println("++++++++++++++++++++++++++++")

			fmt.Println("preItem: ", string(preItem))
			fmt.Println("preItemIsDigit: ", preItemIsDigit)
			fmt.Println("preItemIsSlash: ", preItemIsSlash)

			fmt.Println("++++++++++++++++++++++++++++")

			fmt.Println("item: ", string(item))
			fmt.Println("itemRune: ", item)
			fmt.Println("itemIsDigit:", itemIsDigit)
			fmt.Println("itemIsSlash:", itemIsSlash)


				if prePreItemIsDigit == true && preItemIsDigit == true {
					return "", ErrInvalidString
				}


			if preItemIsDigit && itemIsDigit {
				return "", ErrInvalidString
			}

			if !preItemIsDigit {
				if itemIsDigit {
					itemInt, err := strconv.Atoi(string(item))
					if err != nil {
						return "", err
					}
					sb.WriteString(strings.Repeat(string(preItem), itemInt))
				} else {
					sb.WriteString(string(preItem))
				}

			}

			//fmt.Println("++++++++++++++++++++++++++++")
			prePreItem = preItem
			preItem = item
			prePreItemIsDigit = preItemIsDigit
			preItemIsDigit = itemIsDigit

				fmt.Println("prePreItem: ", string(prePreItem))
				fmt.Println("prePreItemIsDigit: ", prePreItemIsDigit)
				fmt.Println("preItem: ", string(preItem))
				fmt.Println("preItemIsDigit: ", preItemIsDigit)

			fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", i, "<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
		}
		fmt.Println("-------------------------------------------------------")
		return sb.String(), nil

}*/

func stringSize1(in string) (string, error) {
	in0 := rune(in[0])
	in0IsDigit := unicode.IsDigit(in0)
	if in0IsDigit {
		return "", ErrInvalidString
	} else {
		return string(in0), nil
	}
}

func stringSize2(in string) (string, error) {
	var sb strings.Builder
	in0 := rune(in[0])
	in0IsDigit := unicode.IsDigit(in0)
	in1 := rune(in[1])
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
	return sb.String(), nil
}
