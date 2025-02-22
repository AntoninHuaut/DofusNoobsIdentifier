package internal

import (
	"fmt"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"math"
	"regexp"
	"strings"
	"unicode"
)

var (
	locationCache                 = make(map[string]string)
	regexAlphaNumericalHyphenOnly = regexp.MustCompile(`[^a-zA-Z0-9-]+`)
)

func GetLocationFromTarget(titles map[string]string, target string) string {
	if loc, ok := locationCache[target]; ok {
		return loc
	}

	targetFormatted := strings.TrimSpace(strings.ToLower(target))
	minDistance, closestLoc, closestTitle := findClosestStringLevenshtein(targetFormatted, titles)
	if minDistance > 0 {
		fmt.Printf("%d - %s - %s - %s\n", minDistance, target, closestTitle, closestLoc)
	}
	return fmt.Sprintf("%d - %s", minDistance, closestLoc)
}

func findClosestStringLevenshtein(targetFormatted string, titles map[string]string) (int, string, string) {
	minDistance := math.MaxInt64
	closestLoc := ""
	closestTitle := ""
	for url, title := range titles {
		titleFormatted := strings.TrimSpace(strings.ToLower(title))
		distance := levenshteinDistance(titleFormatted, targetFormatted)
		if distance < minDistance {
			minDistance = distance
			closestTitle = title
			closestLoc = url
		}
	}

	return minDistance, closestLoc, closestTitle
}

func levenshteinDistance(s1, s2 string) int {
	lenS1 := len(s1)
	lenS2 := len(s2)

	if lenS1 == 0 {
		return lenS2
	}
	if lenS2 == 0 {
		return lenS1
	}

	prevRow := make([]int, lenS2+1)
	currRow := make([]int, lenS2+1)

	for j := 0; j <= lenS2; j++ {
		prevRow[j] = j
	}

	for i := 1; i <= lenS1; i++ {
		currRow[0] = i
		for j := 1; j <= lenS2; j++ {
			cost := 0
			if s1[i-1] != s2[j-1] {
				cost = 1
			}
			currRow[j] = minInt(prevRow[j]+1, minInt(currRow[j-1]+1, prevRow[j-1]+cost))
		}
		prevRow, currRow = currRow, prevRow
	}

	return prevRow[lenS2]
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func convertSlugDofusPourLesNoobs(slug string) []string {
	slug = strings.TrimSpace(slug)
	slug = strings.ToLower(slug)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "'", "-")
	slug = strings.ReplaceAll(slug, "--", "-")
	slug = strings.ReplaceAll(slug, "œ", "oe")
	slugBis := strings.ReplaceAll(slug, "é", "eacute")

	postEdit := func(s string) string {
		s = removeAccents(s)
		s = regexAlphaNumericalHyphenOnly.ReplaceAllString(s, "")
		s = strings.ReplaceAll(s, "--", "-")

		return s
	}

	slug = postEdit(slug)
	slugBis = postEdit(slugBis)

	return []string{slug, slugBis}
}

func removeAccents(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, err := transform.String(t, s)
	if err != nil {
		return s
	}
	return result
}
