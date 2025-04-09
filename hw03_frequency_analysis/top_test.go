package hw03frequencyanalysis

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Change to true if needed.
var taskWithAsteriskIsCompleted = false

var text = `Как видите, он  спускается  по  лестнице  вслед  за  своим
	другом   Кристофером   Робином,   головой   вниз,  пересчитывая
	ступеньки собственным затылком:  бум-бум-бум.  Другого  способа
	сходить  с  лестницы  он  пока  не  знает.  Иногда ему, правда,
		кажется, что можно бы найти какой-то другой способ, если бы  он
	только   мог   на  минутку  перестать  бумкать  и  как  следует
	сосредоточиться. Но увы - сосредоточиться-то ему и некогда.
		Как бы то ни было, вот он уже спустился  и  готов  с  вами
	познакомиться.
	- Винни-Пух. Очень приятно!
		Вас,  вероятно,  удивляет, почему его так странно зовут, а
	если вы знаете английский, то вы удивитесь еще больше.
		Это необыкновенное имя подарил ему Кристофер  Робин.  Надо
	вам  сказать,  что  когда-то Кристофер Робин был знаком с одним
	лебедем на пруду, которого он звал Пухом. Для лебедя  это  было
	очень   подходящее  имя,  потому  что  если  ты  зовешь  лебедя
	громко: "Пу-ух! Пу-ух!"- а он  не  откликается,  то  ты  всегда
	можешь  сделать вид, что ты просто понарошку стрелял; а если ты
	звал его тихо, то все подумают, что ты  просто  подул  себе  на
	нос.  Лебедь  потом  куда-то делся, а имя осталось, и Кристофер
	Робин решил отдать его своему медвежонку, чтобы оно не  пропало
	зря.
		А  Винни - так звали самую лучшую, самую добрую медведицу
	в  зоологическом  саду,  которую  очень-очень  любил  Кристофер
	Робин.  А  она  очень-очень  любила  его. Ее ли назвали Винни в
	честь Пуха, или Пуха назвали в ее честь - теперь уже никто  не
	знает,  даже папа Кристофера Робина. Когда-то он знал, а теперь
	забыл.
		Словом, теперь мишку зовут Винни-Пух, и вы знаете почему.
		Иногда Винни-Пух любит вечерком во что-нибудь поиграть,  а
	иногда,  особенно  когда  папа  дома,  он больше любит тихонько
	посидеть у огня и послушать какую-нибудь интересную сказку.
		В этот вечер...`

func TestTop10(t *testing.T) {
	t.Run("no words in empty string", func(t *testing.T) {
		require.Len(t, Top10(""), 0)
	})

	t.Run("positive test", func(t *testing.T) {
		if taskWithAsteriskIsCompleted {
			expected := []string{
				"а",         // 8
				"он",        // 8
				"и",         // 6
				"ты",        // 5
				"что",       // 5
				"в",         // 4
				"его",       // 4
				"если",      // 4
				"кристофер", // 4
				"не",        // 4
			}
			require.Equal(t, expected, Top10(text))
		} else {
			expected := []string{
				"он",        // 8
				"а",         // 6
				"и",         // 6
				"ты",        // 5
				"что",       // 5
				"-",         // 4
				"Кристофер", // 4
				"если",      // 4
				"не",        // 4
				"то",        // 4
			}
			require.Equal(t, expected, Top10(text))
		}
	})
}

// ------------------------------------------------------------------------------------------------------------------------
// тест для примера из readme задания

var textTop7 = "cat and dog, one dog,two cats and one man"

func TestTop7(t *testing.T) {
	t.Run("no words in empty string", func(t *testing.T) {
		require.Len(t, Top10(""), 0)
	})
	t.Run("positive test", func(t *testing.T) {
		expected := []string{
			"and",     // 2
			"one",     // 2
			"cat",     // 1
			"cats",    // 1
			"dog,",    // 1
			"dog,two", // 1
			"man",     // 1
		}
		require.Equal(t, expected, Top10(textTop7))
	})
}

// ------------------------------------------------------------------------------------------------------------------------
// тесты для вспомогательных функций

func TestDetermineNumberOfOccurrences(t *testing.T) {
	in := []string{"cat", "and", "dog,two", "dog,", "one", "dog,two", "cats", "and", "one", "dog,two", "man"}
	outExpected := map[string]int{"cat": 1, "and": 2, "cats": 1, "dog,": 1, "dog,two": 3, "man": 1, "one": 2}
	t.Run("positive test determineNumberOfOccurrences func", func(t *testing.T) {
		outActual := determineNumberOfOccurrences(in)
		require.Equal(t, outExpected, outActual)
	})
}

func TestGroupByOccurrence(t *testing.T) {
	in := map[string]int{"cat": 1, "and": 2, "cats": 1, "dog,": 1, "dog,two": 3, "man": 1, "one": 2}
	maxExpected := 3
	mapExpected := map[int][]string{1: {"man", "cat", "cats", "dog,"}, 2: {"one", "and"}, 3: {"dog,two"}}
	t.Run("positive test groupByOccurrence func", func(t *testing.T) {
		maxActual, mapActual := groupByOccurrence(in)
		require.Equal(t, maxExpected, maxActual)
		require.Equal(t, mapExpected, mapActual)
	})
}

func TestSortGroupByOccurrence(t *testing.T) {
	in := map[int][]string{1: {"man", "cat", "cats", "dog,"}, 2: {"one", "and"}, 3: {"dog,two"}}
	outExpected := map[int][]string{1: {"cat", "cats", "dog,", "man"}, 2: {"and", "one"}, 3: {"dog,two"}}
	t.Run("positive test sortGroupedOccurrence func", func(t *testing.T) {
		sortGroupedOccurrence(in)
		outActual := in
		require.Equal(t, outExpected, outActual)
	})
}

func TestBuildOutSlice(t *testing.T) {
	inMax := 3
	in := map[int][]string{1: {"cat", "cats", "dog,", "man"}, 2: {"and", "one"}, 3: {"dog,two"}}
	outExpected := []string{"dog,two", "and", "one", "cat", "cats", "dog,", "man"}
	t.Run("positive test buildOutSlice func", func(t *testing.T) {
		sortGroupedOccurrence(in)
		outActual := buildOutSlice(inMax, in)
		require.Equal(t, outExpected, outActual)
	})
}
