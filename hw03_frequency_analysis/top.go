package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

//-----------------------------------------------------------------------------------------------------------

func Top10(in string) []string {
	// Place your code here.
	if in == "" { // отсекаем исключение - пустая строка
		return nil
	}
	inArray := strings.Fields(in) // разделяем по отступам
	occurrenceMap := determineNumberOfOccurrences(inArray)
	maxOccurrence, groupedOccurrenceMap := groupByOccurrence(occurrenceMap)
	sortGroupedOccurrence(groupedOccurrenceMap)
	out := buildOutSlice(maxOccurrence, groupedOccurrenceMap)
	return out
}

//----------------------------------------------------------------------------------------------------
// вспомогательные функции

// определяем количество вхождений по каждому слову ([слово]число вхождений).
func determineNumberOfOccurrences(inArray []string) map[string]int {
	occurrenceMap := map[string]int{}
	for _, item := range inArray {
		occurrenceMap[item]++
	}
	return occurrenceMap
}

// строим карту вхождений ([число вхождений]слово)
// определяем максимальное число вхождений для использование в алгоритме сортировки по алфавиту.
func groupByOccurrence(occurrenceMap map[string]int) (int, map[int][]string) {
	maxOccurrence := 0
	groupedOccurrenceMap := map[int][]string{}
	for key1, value1 := range occurrenceMap {
		if value1 > maxOccurrence {
			maxOccurrence = value1
		}
		key2 := value1
		value2 := groupedOccurrenceMap[key2]
		value2 = append(value2, key1)
		groupedOccurrenceMap[key2] = value2
	}
	return maxOccurrence, groupedOccurrenceMap
}

// сортируем сгрупированные слова в алфавином порядке.
func sortGroupedOccurrence(groupedOccurrenceMap map[int][]string) {
	for key2, value2 := range groupedOccurrenceMap {
		groupedOccurrenceMap[key2] = value2
		sort.Slice(value2, func(i, j int) bool {
			return value2[i] < value2[j]
		})
	}
}

// выстраиваем сгруппированые слова в одну последовательность, в порядке от максимальных вхождений к минимальным
// ограничиваем длину последовательности 10-ю словами.
const OutSizeMax = 10

func buildOutSlice(maxOccurrence int, groupedOccurrenceMap map[int][]string) []string {
	count := 0
	out := make([]string, 0, OutSizeMax)
	for i := maxOccurrence; i > 0; i-- {
		value2 := groupedOccurrenceMap[i]
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
	return out
}
