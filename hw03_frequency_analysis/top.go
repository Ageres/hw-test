package hw03frequencyanalysis

import (
	"fmt"
	"strings"
)

func Top10(in string) []string {
	// Place your code here.

	in = strings.ToLower(in)
	in = strings.ReplaceAll(in, " ", "_")
	in = strings.ReplaceAll(in, "\n", "_")
	in = strings.ReplaceAll(in, "\t", "_")
	in = strings.ReplaceAll(in, "\"", "_")
	in = strings.ReplaceAll(in, ",", "_")
	in = strings.ReplaceAll(in, ".", "_")
	in = strings.ReplaceAll(in, "!", "_")

	for range 20 {
		in = strings.ReplaceAll(in, "__", "_")
	}

	fmt.Println("---------------------00----------------------")
	fmt.Println(in)
	inArray := strings.Split(in, "_")
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")

	fmt.Println("---------------------01----------------------")
	outMap := map[string]int{}
	for i := range inArray {
		item := inArray[i]
		//item = strings.Replace(item, " ", "", -1)
		//item = strings.Replace(item, ",", "", -1)
		//inArray[i] = item
		//fmt.Println(i, ":    ", item, " | ", []rune(item))

		value := outMap[item]
		value++
		outMap[item] = value
	}
	fmt.Println("---------------------02----------------------")

	var max1, max2, max3, max4, max5, max6, max7, max8, max9, max10 int
	var max1Str, max2Str, max3Str, max4Str, max5Str, max6Str, max7Str, max8Str, max9Str, max10Str string

	for key, value := range outMap {
		fmt.Println(key, ":    ", value)
		if value > max1 {
			max10 = max9
			max9 = max8
			max8 = max7
			max7 = max6
			max6 = max5
			max5 = max4
			max4 = max3
			max3 = max2
			max2 = max1
			max1 = value

			max10Str = max9Str
			max9Str = max8Str
			max8Str = max7Str
			max7Str = max6Str
			max6Str = max5Str
			max5Str = max4Str
			max4Str = max3Str
			max3Str = max2Str
			max2Str = max1Str
			max1Str = key
		}
	}
	fmt.Println("---------------------03----------------------")

	fmt.Println(max1, ":    ", max1Str)
	fmt.Println(max2, ":    ", max2Str)
	fmt.Println(max3, ":    ", max3Str)
	fmt.Println(max4, ":    ", max4Str)
	fmt.Println(max5, ":    ", max5Str)
	fmt.Println(max6, ":    ", max6Str)
	fmt.Println(max7, ":    ", max7Str)
	fmt.Println(max8, ":    ", max8Str)
	fmt.Println(max9, ":    ", max9Str)
	fmt.Println(max10, ":    ", max10Str)

	fmt.Println("---------------------04----------------------")
	out := []string{max10Str, max9Str, max8Str, max7Str, max6Str, max5Str, max4Str, max3Str, max2Str, max1Str}
	fmt.Println(out)

	fmt.Println("---------------------05----------------------")
	fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
	//strings.Fields()
	return out
}
