package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type WordCount struct {
	Word  string
	Count int
}

func Top10(text string) []string {
	textSlice := strings.Fields(text)
	counts := make(map[string]int)

	for _, word := range textSlice {
		counts[word] = counts[word] + 1 //nolint
	}

	keys := make([]string, 0, len(counts))
	for k := range counts {
		keys = append(keys, k)
	}
	wordCounts := make([]WordCount, 0, len(keys))
	for _, word := range keys {
		wordCounts = append(wordCounts, WordCount{word, counts[word]})
	}
	// Сортировка слайса по алфавиту и частоте
	sort.SliceStable(wordCounts, func(i, j int) bool {
		if wordCounts[i].Count == wordCounts[j].Count {
			// Если частоты равны, сортируем по слову (лексикографически)
			return wordCounts[i].Word < wordCounts[j].Word
		}
		// Если частоты разные, сортируем по частоте (по убыванию)
		return wordCounts[i].Count > wordCounts[j].Count
	})
	topN := 10
	if len(wordCounts) == 0 {
		return nil
	}
	resultStruct := wordCounts[:topN]
	res := make([]string, 0)
	for _, v := range resultStruct {
		res = append(res, v.Word)
	}
	return res
}
