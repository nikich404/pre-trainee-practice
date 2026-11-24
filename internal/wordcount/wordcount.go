package wordcount

import (
	"sort"
	"strings"
	"unicode"
)

type Pair struct {
	Word  string
	Count int
}

// Tokens splits s into words using the tokenizer rules below.
func Tokens(s string) []string {
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsDigit(r) && !unicode.IsLetter(r)
	})
	for i, p := range parts {
		parts[i] = strings.ToLower(p)
	}
	return parts
}

// Count returns frequency of each token.
func Count(tokens []string) map[string]int {
	freq := make(map[string]int, len(tokens))
	for _, t := range tokens {
		freq[t]++
	}
	return freq
}

// TopN returns up to n pairs, highest count first.
// Ties: alphabetical ascending by Word.
// If n <= 0, return an empty slice.
func TopN(freq map[string]int, n int) []Pair {
	pairs := make([]Pair, 0, len(freq))
	if n <= 0 {
		return []Pair{}
	}

	for w, c := range freq {
		pairs = append(pairs, Pair{
			Word:  w,
			Count: c,
		})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].Count != pairs[j].Count {
			return pairs[i].Count > pairs[j].Count // больший Count раньше
		}
		return pairs[i].Word < pairs[j].Word
	})

	if n < len(pairs) {
		pairs = pairs[:n]
	}
	return pairs
}

// DropStop removes tokens that appear in stop (lowercase keys).
func DropStop(tokens []string, stop map[string]struct{}) []string {
	out := make([]string, 0, len(tokens))

	for _, v := range tokens {
		if _, ok := stop[v]; ok {
			continue
		}
		stop[v] = struct{}{}
		out = append(out, v)
	}
	return out
}
