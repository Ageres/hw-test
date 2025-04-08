package hw03frequencyanalysis

import (
	"fmt"
	"sort"
	"strings"
)

//-----------------------------------------------------------------------------------------------------------

func Top7WithOutAsterisk(in string) []string {
	// Place your code here.
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	fmt.Println("in:", in)
	fmt.Println("---------------------01----------------------")
	// отсекаем исключение - пустая строка
	if in == "" {
		return []string{}
	}
	fmt.Println("---------------------02---------------------")
	// разделяем по отступам
	inArray := strings.Fields(in)
	for i := range inArray {
		fmt.Println("inArray[", i, "]", inArray[i])
	}
	fmt.Println("---------------------03----------------------")
	// определяем количество вхождений по каждому слову ([слово]число вхождений )
	occurrenceMap := determineNumberOfOccurrences(inArray)
	fmt.Println("---------------------04----------------------")
	// строим карту группировки слов по числу вхождений ([число вхождений]слово)
	// определяем максимальное число вхождений для использования в алгоритме сортировки по алфавиту
	maxOccurrence, groupedOccurrenceMap := buildOutMap(occurrenceMap)
	fmt.Println("---------------------05----------------------")
	sortOutMap(groupedOccurrenceMap)
	fmt.Println("---------------------06----------------------")
	out := buildOutSlice(maxOccurrence, groupedOccurrenceMap)
	fmt.Println("---------------------07----------------------")
	fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
	return out
}

//-----------------------------------------------------------------------------------------------------------

func Top10(in string) []string {
	// Place your code here.
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	fmt.Println("---------------------01----------------------")
	fmt.Println(in)
	if in == "" {
		return []string{}
	}
	fmt.Println("---------------------02----------------------")
	in = removeOtherSymbols(in)
	fmt.Println("---------------------03----------------------")
	inArray := strings.Fields(in)
	for i := range inArray {
		fmt.Println(inArray[i])
	}
	fmt.Println("---------------------04----------------------")
	outMap1 := determineNumberOfOccurrences(inArray)
	fmt.Println("---------------------05----------------------")
	maxLen, outMap2 := buildOutMap(outMap1)
	fmt.Println("---------------------06----------------------")
	sortOutMap(outMap2)
	fmt.Println("---------------------07----------------------")
	out := buildOutSlice(maxLen, outMap2)
	fmt.Println("---------------------08----------------------")
	fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
	//strings.Fields()
	return out
}

//----------------------------------------------------------------------------------------------------

//----------------------------------------------------------------------------------------------------

func Top7TaskWithAsterisk(in string) []string {
	return []string{}
}

//----------------------------------------------------------------------------------------------------

func removeOtherSymbols(in string) string {
	in = strings.ToLower(in)
	in = strings.ReplaceAll(in, "_", " ")
	in = strings.ReplaceAll(in, "\n", " ")
	in = strings.ReplaceAll(in, "\t", " ")
	in = strings.ReplaceAll(in, "\"", " ")
	in = strings.ReplaceAll(in, ",", " ")
	in = strings.ReplaceAll(in, ".", " ")
	in = strings.ReplaceAll(in, "!", " ")
	in = strings.ReplaceAll(in, ":", " ")
	in = strings.ReplaceAll(in, ";", " ")
	in = strings.ReplaceAll(in, "?", " ")
	fmt.Println(in)
	return in
}

// определяем количество вхождений по каждому слову ([слово]число вхождений )
func determineNumberOfOccurrences(inArray []string) map[string]int {
	occurrenceMap := map[string]int{}
	for i := range inArray {
		item := inArray[i]
		if item == "-" {
			continue
		}
		value := occurrenceMap[item]
		value++
		occurrenceMap[item] = value
	}
	for key1, value1 := range occurrenceMap {
		fmt.Println(value1, ":", key1)
	}
	return occurrenceMap
}

// строим карту вхождений ([число вхождений]слово)
// определяем максимальное число вхождений для использование в алгоритме сортировки по алфавиту
func buildOutMap(outMap1 map[string]int) (int, map[int][]string) {
	maxLen := 0
	outMap2 := map[int][]string{}
	for key1, value1 := range outMap1 {
		if value1 > maxLen {
			maxLen = value1
		}
		key2 := value1
		value2 := outMap2[key2]
		value2 = append(value2, key1)
		outMap2[key2] = value2
	}
	for key2, value2 := range outMap2 {
		fmt.Println(key2, ": ", value2)
	}
	return maxLen, outMap2
}

func sortOutMap(outMap2 map[int][]string) {
	for key2, value2 := range outMap2 {
		outMap2[key2] = value2
		sort.Slice(value2, func(i, j int) bool {
			return value2[i] < value2[j]
		})
		fmt.Println(key2, ": ", value2)
	}
}

func buildOutSlice(maxLen int, outMap2 map[int][]string) []string {
	count := 0
	out := make([]string, 0, 10)
	for i := maxLen; i > 0; i-- {
		value2 := outMap2[i]
		if value2 == nil {
			continue
		}
		for j := range value2 {
			out = append(out, value2[j])
			count++
			if count == 10 {
				break
			}
		}
		if count == 10 {
			break
		}
	}
	fmt.Println(out)
	return out
}
