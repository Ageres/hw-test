package hw03frequencyanalysis

import (
	"cmp"
	"fmt"
	"slices"
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
	wordItems := buildWordItems(occurrenceMap)
	sortedWordItems := sortWordItem(wordItems)
	out := buildOutSlice(sortedWordItems)
	return out
}

//----------------------------------------------------------------------------------------------------
// вспомогательные функции

type WordItem struct {
	Occurrence int    // число вхождений
	Word       string // слово
}

// определяем количество вхождений по каждому слову ([слово]число вхождений).
func determineNumberOfOccurrences(inArray []string) map[string]int {
	occurrenceMap := map[string]int{}
	for _, item := range inArray {
		occurrenceMap[item]++
	}
	return occurrenceMap
}

// строим слайс объектов WordItem {число вхождений, слово}.
func buildWordItems(occurrenceMap map[string]int) []WordItem {
	var wordItems []WordItem
	for word, occurrence := range occurrenceMap {
		wordItem := WordItem{
			Occurrence: occurrence,
			Word:       word,
		}
		wordItems = append(wordItems, wordItem)
	}
	return wordItems
}

// сортируем объекты WordItem по вхождению и алфавитному порядку.
func sortWordItem(wordItems []WordItem) []WordItem {
	slices.SortFunc(wordItems, func(a, b WordItem) int {
		return cmp.Or(
			cmp.Compare(b.Occurrence, a.Occurrence),
			cmp.Compare(a.Word, b.Word),
		)
	})
	for i, w := range wordItems {
		fmt.Println("---------i[", i, "] = ", w)
	}
	return wordItems
}

// выстраиваем сгруппированые слова в одну последовательность, в порядке от максимальных вхождений к минимальным
// ограничиваем длину последовательности 10-ю словами.
const OutSizeMax = 10

func buildOutSlice(wordItems []WordItem) []string {
	maxItem := min(len(wordItems), OutSizeMax)
	out := make([]string, 0, maxItem)
	for i := range maxItem {
		wordItem := wordItems[i]
		out = append(out, wordItem.Word)
	}
	return out
}
