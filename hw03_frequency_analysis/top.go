package hw03frequencyanalysis

import (
	"slices"
	"sort"
	"strings"
)

const mostFrequentCount = 10

func Top10(s string) []string {
	// map[string]int
	wordFrequentMap := stringToMap(s)
	// []int
	maxFrequencyCounts := prepareMaxFrequencyCounts(wordFrequentMap)
	// map[int][]string
	frequencyCountWordMap := prepareFrequencyCountWordMap(maxFrequencyCounts, wordFrequentMap)

	return getSortedMostFrequentWords(frequencyCountWordMap)
}

func getSortedMostFrequentWords(frequencyCountWordMap map[int][]string) []string {
	maxFrequencyCounts := make([]int, 0, len(frequencyCountWordMap))
	for count := range frequencyCountWordMap {
		maxFrequencyCounts = append(maxFrequencyCounts, count)
	}
	sort.Slice(maxFrequencyCounts, func(i, j int) bool {
		return maxFrequencyCounts[i] > maxFrequencyCounts[j]
	})

	result := make([]string, 0, mostFrequentCount)
	for _, count := range maxFrequencyCounts {
		if len(result) > mostFrequentCount {
			break
		}

		words := frequencyCountWordMap[count]
		slices.Sort(words)
		for _, w := range words {
			if len(result) > mostFrequentCount {
				break
			}
			result = append(result, w)
		}
	}

	return result
}

func prepareFrequencyCountWordMap(maxFrequencyCounts []int, wordsMap map[string]int) map[int][]string {
	result := make(map[int][]string, len(maxFrequencyCounts))
	for _, count := range maxFrequencyCounts {
		result[count] = make([]string, 0)
	}

	for w, count := range wordsMap {
		if _, ok := result[count]; ok {
			result[count] = append(result[count], w)
		}
	}

	return result
}

func prepareMaxFrequencyCounts(wordsMap map[string]int) []int {
	frequency := make([]int, 0, len(wordsMap))
	for _, count := range wordsMap {
		frequency = append(frequency, count)
	}
	sort.Slice(frequency, func(i, j int) bool {
		return frequency[i] > frequency[j]
	})

	counts := len(frequency)
	if counts > mostFrequentCount {
		counts = mostFrequentCount
	}

	return frequency[0:counts]
}

func stringToMap(s string) map[string]int {
	words := strings.Fields(s)
	result := make(map[string]int)
	for _, w := range words {
		result[w]++
	}

	return result
}
