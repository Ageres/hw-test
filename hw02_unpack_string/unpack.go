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
	//inRunes := []rune(in)

	var sb strings.Builder

	var prePreItem rune
	var prePreItemIsDigit bool
	var prePreItemIsLetter bool
	//var prePreItemIsSlash bool

	var preItem rune
	var preItemIsDigit bool
	var preItemIsLetter bool
	//var preItemIsSlash bool

	var itemIsDigit bool
	var itemIsLetter bool
	var itemIsSlash bool

	for i, item := range in {
		fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", i, ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
		itemIsDigit = unicode.IsDigit(item)
		itemIsLetter = unicode.IsLetter(item)
		itemIsSlash = (item == 92)

		//---------------------------

		fmt.Println("prePreItem: ", string(prePreItem))
		fmt.Println("prePreItemIsDigit: ", prePreItemIsDigit)
		fmt.Println("prePreItemIsLetter: ", prePreItemIsLetter)

		fmt.Println("++++++++++++++++++++++++++++")

		fmt.Println("preItem: ", string(preItem))
		fmt.Println("preItemIsDigit: ", preItemIsDigit)
		fmt.Println("preItemIsLetter: ", preItemIsLetter)

		fmt.Println("++++++++++++++++++++++++++++")

		fmt.Println("item: ", string(item))
		fmt.Println("itemRune: ", item)
		fmt.Println("itemIsDigit:", itemIsDigit)
		fmt.Println("itemIsLetter:", itemIsLetter)
		fmt.Println("itemIsSlash:", itemIsSlash)

		/*
			if i == 0 && itemIsDigit {
				return "", ErrInvalidString
			}
		*/

		/*
			if prePreItemIsDigit == true && preItemIsDigit == true {
				return "", ErrInvalidString
			}
		*/

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
		/*
			fmt.Println("prePreItem: ", string(prePreItem))
			fmt.Println("prePreItemIsDigit: ", prePreItemIsDigit)
			fmt.Println("preItem: ", string(preItem))
			fmt.Println("preItemIsDigit: ", preItemIsDigit)
		*/
		fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", i, "<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
	}
	fmt.Println("-------------------------------------------------------")
	return sb.String(), nil
}
