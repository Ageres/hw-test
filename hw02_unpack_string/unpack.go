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

	// анализируемая строка содержит 1 символ
	if inSize == 1 {
		return stringSize1(in)
	}

	// анализируемая строка содержит 2 символа
	if inSize == 2 {
		return stringSize2(in)
	}

	var sb strings.Builder

	var prePreItem rune = rune(in[0])
	var prePreItemIsDigit bool = unicode.IsDigit(rune(in[0]))
	var prePreItemIsSlash bool = rune(in[0]) == 92
	var isWritePrePreItem bool = false

	var preItem rune = rune(in[1])
	var preItemIsDigit bool = unicode.IsDigit(rune(in[1]))
	var preItemIsSlash bool = rune(in[1]) == 92
	var isWritePreItem bool = false

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
		fmt.Println("isWritePrePreItem:", isWritePrePreItem)

		fmt.Println("++++++++++++++++++++++++++++")

		fmt.Println("preItem: ", string(preItem))
		fmt.Println("preItemIsDigit: ", preItemIsDigit)
		fmt.Println("preItemIsSlash: ", preItemIsSlash)
		fmt.Println("isWritePreItem:", isWritePreItem)

		fmt.Println("++++++++++++++++++++++++++++")

		fmt.Println("item: ", string(item))
		fmt.Println("itemRune: ", item)
		fmt.Println("itemIsDigit:", itemIsDigit)
		fmt.Println("itemIsSlash:", itemIsSlash)

		fmt.Println("++++++++++++++++++++++++++++")

		//---------------------------------------------------

		//----------------------------------------------

		prePreItem = preItem
		preItem = item
		prePreItemIsDigit = preItemIsDigit
		preItemIsDigit = itemIsDigit

		fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", i, "<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
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
	*/
}

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
