package hw03_frequency_analysis //nolint:golint,stylecheck
import (
	"regexp"
	"strings"
)

// Create a dict with words and its frequency.
func getDict(txt string) map[string]int {
	wordDict := map[string]int{}
	words := regexp.MustCompile(`[,. !	\s]+`).Split(txt, -1)
	for _, w := range words {
		// skip empty and dash strings
		if w != "" && w != "-" {
			wordDict[strings.ToLower(w)]++
		}
	}

	return wordDict
}

// Find the index in the result slice with the smallest frequency.
func findIdxOfMin(wordDict map[string]int, topWords []string) int {
	resIdx := 0
	for i, w := range topWords {
		if w == "" {
			return i
		}
		if wordDict[w] < wordDict[topWords[resIdx]] {
			resIdx = i
		}
	}

	return resIdx
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func Top10(txt string) []string {
	if len(txt) == 0 {
		return []string{}
	}

	wordDict := getDict(txt)
	topWords := make([]string, min(len(wordDict), 10))
	// index with the least common word in topWords
	minIdx := 0
	for w, frequency := range wordDict {
		if topWords[minIdx] == "" || frequency > wordDict[topWords[minIdx]] {
			topWords[minIdx] = w
			minIdx = findIdxOfMin(wordDict, topWords)
		}
	}

	return topWords
}
