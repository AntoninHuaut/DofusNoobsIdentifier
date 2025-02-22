package internal

import (
	"fmt"
	"math"
	"strings"
)

var (
	locationCache = make(map[string]string)
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
