package infrastructure

import (
	"slices"
	"strconv"
	"strings"
	"unicode"
)

const TitleSimilarityThreshold = 0.80

func normalizeProgramTitle(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	if s == "" {
		return ""
	}
	var b strings.Builder
	b.Grow(len(s))
	lastSpace := false
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			lastSpace = false
		} else if unicode.IsSpace(r) || r == '-' || r == '_' || r == '/' || r == '.' || r == ',' {
			if !lastSpace && b.Len() > 0 {
				b.WriteByte(' ')
				lastSpace = true
			}
		}
	}
	return strings.TrimSpace(b.String())
}

func extractNumericTokens(norm string) []string {
	fields := strings.Fields(norm)
	var nums []string
	for _, f := range fields {
		if _, err := strconv.Atoi(f); err == nil {
			nums = append(nums, f)
		}
	}
	return nums
}

func charBigramDice(a, b []rune) float64 {
	if len(a) < 2 || len(b) < 2 {
		if string(a) == string(b) {
			return 1
		}
		return 0
	}
	counts := make(map[string]int, len(a))
	for i := 0; i < len(a)-1; i++ {
		counts[string(a[i:i+2])]++
	}
	overlap := 0
	for i := 0; i < len(b)-1; i++ {
		bg := string(b[i : i+2])
		if counts[bg] > 0 {
			overlap++
			counts[bg]--
		}
	}
	return float64(2*overlap) / float64((len(a)-1)+(len(b)-1))
}

func levenshteinRunes(a, b []rune) int {
	if len(a) == 0 {
		return len(b)
	}
	if len(b) == 0 {
		return len(a)
	}
	prev := make([]int, len(b)+1)
	curr := make([]int, len(b)+1)
	for j := range prev {
		prev[j] = j
	}
	for i := 1; i <= len(a); i++ {
		curr[0] = i
		for j := 1; j <= len(b); j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			del := prev[j] + 1
			ins := curr[j-1] + 1
			sub := prev[j-1] + cost
			m := del
			if ins < m {
				m = ins
			}
			if sub < m {
				m = sub
			}
			curr[j] = m
		}
		prev, curr = curr, prev
	}
	return prev[len(b)]
}

// ComputeTitleSimilarity returns a similarity score in [0, 1] between two program titles.
func ComputeTitleSimilarity(a, b string) float64 {
	na := normalizeProgramTitle(a)
	nb := normalizeProgramTitle(b)
	if na == "" || nb == "" {
		return 0
	}
	if na == nb {
		return 1
	}

	ra := []rune(na)
	rb := []rune(nb)
	// Avoid false positives on very short identifiers (e.g. "tes1" vs "tes2").
	if len(ra) < 4 || len(rb) < 4 {
		return 0
	}

	// Distinct numbered items (e.g. "Rapat Kerja 1" vs "Rapat Kerja 2") are intentionally distinct prokers.
	numsA := extractNumericTokens(na)
	numsB := extractNumericTokens(nb)
	if len(numsA) > 0 && len(numsB) > 0 && !slices.Equal(numsA, numsB) {
		return 0
	}

	maxLen := len(ra)
	if len(rb) > maxLen {
		maxLen = len(rb)
	}
	levSim := 1.0 - float64(levenshteinRunes(ra, rb))/float64(maxLen)
	diceSim := charBigramDice(ra, rb)
	return 0.5*levSim + 0.5*diceSim
}

func IsSimilarProgramTitle(a, b string) bool {
	return ComputeTitleSimilarity(a, b) >= TitleSimilarityThreshold
}
