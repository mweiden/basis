package anagrams

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func Anagrams(corpus []string) map[string][]string {
	result := make(map[string][]string)

	for _, word := range corpus {
		key := makeKey(word)
		_, found := result[key]
		if !found {
			result[key] = []string{word}
		} else {
			result[key] = append(result[key], word)
		}
	}
	return result
}

func makeKey(word string) string {
	countMap := make(map[rune]int)
	// count runes
	for _, c := range word {
		if c == ' ' {
			continue
		}
		_, found := countMap[c]
		if !found {
			countMap[c] = 1
		} else {
			countMap[c] += 1
		}
	}
	// build key tokens
	var tokens []string
	for k, v := range countMap {
		var countString string
		if v == 1 {
			countString = ""
		} else {
			countString = strconv.Itoa(v)
		}
		tokens = append(tokens, fmt.Sprintf("%c%s", k, countString))
	}
	// sort and return joined tokens
	sort.Slice(
		tokens,
		func(i, j int) bool {
			return tokens[i][0] < tokens[j][0]
		},
	)
	return strings.Join(tokens, "")
}
