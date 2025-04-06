package hw03frequencyanalysis

import (
	"fmt"
	"sort"
	"strings"
)

func Top10(in string) []string {
	// Place your code here.
	if in == "" {
		return []string{}
	}

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
	outMap1 := map[string]int{}
	for i := range inArray {
		item := inArray[i]
		if item == "-" {
			continue
		}
		//item = strings.Replace(item, " ", "", -1)
		//item = strings.Replace(item, ",", "", -1)
		//inArray[i] = item
		//fmt.Println(i, ":    ", item, " | ", []rune(item))

		value := outMap1[item]
		value++
		outMap1[item] = value
	}
	fmt.Println("---------------------02----------------------")

	for key, value := range outMap1 {
		fmt.Println(key, ":    ", value)

	}

	fmt.Println("---------------------03----------------------")

	max := 0
	outMap2 := map[int][]string{}
	for key1, value1 := range outMap1 {
		if value1 > max {
			max = value1
		}
		key2 := value1
		value2 := outMap2[key2]
		value2 = append(value2, key1)
		outMap2[key2] = value2
	}

	fmt.Println("---------------------04----------------------")

	for key2, value2 := range outMap2 {
		fmt.Println(key2, ": ", value2)
		outMap2[key2] = value2
		sort.Slice(value2, func(i, j int) bool {
			return value2[i] < value2[j]
		})
	}

	fmt.Println("---------------------05----------------------")
	for key2, value2 := range outMap2 {
		outMap2[key2] = value2
		fmt.Println(key2, ": ", value2)
	}
	fmt.Println("---------------------06----------------------")

	count := 0
	out := []string{}
	for i := max; i > 0; i-- {
		value2 := outMap2[i]
		if value2 == nil {
			continue
		}
		out = append(out, value2...)
		count++
		if count > 10 {
			break
		}
	}

	fmt.Println(out)
	fmt.Println("---------------------07----------------------")
	if len(out) > 10 {
		out = out[0:10]
	}
	fmt.Println(out)
	fmt.Println("---------------------08----------------------")
	fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
	//strings.Fields()
	return out
}

//----------------------------------------------------------------------------------------------------

func Top7(in string) []string {
	// Place your code here.

	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")

	fmt.Println("---------------------00----------------------")

	fmt.Println(in)
	if in == "" {
		return []string{}
	}

	fmt.Println("---------------------01---------------------")

	inArray := strings.Fields(in)
	for i := range inArray {
		fmt.Println(inArray[i])
	}
	//fmt.Println(inArray)
	fmt.Println("---------------------02----------------------")
	outMap1 := map[string]int{}
	for i := range inArray {
		item := inArray[i]
		value := outMap1[item]
		value++
		outMap1[item] = value
	}
	fmt.Println("---------------------03----------------------")

	for key, value := range outMap1 {
		fmt.Println(key, ":    ", value)

	}
	fmt.Println("---------------------04----------------------")

	max := 0
	outMap2 := map[int][]string{}
	for key1, value1 := range outMap1 {
		if value1 > max {
			max = value1
		}
		key2 := value1
		value2 := outMap2[key2]
		value2 = append(value2, key1)
		outMap2[key2] = value2
	}

	fmt.Println("---------------------05----------------------")
	for key2, value2 := range outMap2 {
		fmt.Println(key2, ": ", value2)
		outMap2[key2] = value2
		sort.Slice(value2, func(i, j int) bool {
			return value2[i] < value2[j]
		})
	}
	fmt.Println("---------------------06----------------------")
	for key2, value2 := range outMap2 {
		outMap2[key2] = value2
		fmt.Println(key2, ": ", value2)
	}
	fmt.Println("---------------------07----------------------")

	count := 0
	out := []string{}
	for i := max; i > 0; i-- {
		value2 := outMap2[i]
		if value2 == nil {
			continue
		}
		out = append(out, value2...)
		count++
		if count > 10 {
			break
		}
	}

	fmt.Println(out)
	fmt.Println("---------------------08----------------------")
	if len(out) > 10 {
		out = out[0:10]
	}
	fmt.Println(out)
	fmt.Println("---------------------09----------------------")
	fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
	return out
}
