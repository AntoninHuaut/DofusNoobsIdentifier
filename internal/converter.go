package internal

import (
	"fmt"
	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"regexp"
	"strings"
)

type locationValue struct {
	Location   string
	Similarity string
}

const (
	dofusTouchSuffix         = "(Dofus Touch)"
	ExactSimilarity          = "exact"
	PrefixOrSuffixSimilarity = "prefixOrSuffix"
	NoneSimilarity           = ""
)

var (
	locationCache            = make(map[string]locationValue)
	stripAlignmentOrderTitle = regexp.MustCompile(`(alignement|ordre) \d+ :\s*`)
)

func FormatGeneral(general string) string {
	general = strings.ToLower(strings.TrimSpace(general))
	general = strings.ReplaceAll(general, "œ", "oe")
	general = strings.ReplaceAll(general, "\u00a0", " ") // NBSP char
	return general
}

func FormatLog(bestSimilarity string, id string, targetFormatted string, closestTitle string, closestLoc string) string {
	return fmt.Sprintf("%14s | %10s |%60s | %60s (-> %s)\n", bestSimilarity, id, targetFormatted, closestTitle, closestLoc)
}

func formatTitle(title string) string {
	title = stripAlignmentOrderTitle.ReplaceAllString(title, "")
	return title
}

func GetLocationFromTarget(id int, titles map[string]string, targetKey, targetFormatted string) (string, string, string) {
	if loc, ok := locationCache[targetKey]; ok {
		return loc.Location, loc.Similarity, ""
	}

	bestSimilarity, closestLoc, closestTitle := findClosestString(targetFormatted, titles)
	locationCache[targetKey] = locationValue{Location: closestLoc, Similarity: bestSimilarity}
	return closestLoc, bestSimilarity, FormatLog(bestSimilarity, fmt.Sprintf("%d", id), targetFormatted, closestTitle, closestLoc)
}

func findClosestString(targetFormatted string, titles map[string]string) (string, string, string) {
	bestSimilarityType := NoneSimilarity
	bestSimilarity := .0

	closestLoc := ""
	closestTitle := ""

	for url, title := range titles {
		if strings.HasSuffix(title, dofusTouchSuffix) {
			continue
		}

		var similarity float64
		var similarityType string

		titleFormatted := formatTitle(FormatGeneral(title))
		if titleFormatted == targetFormatted {
			similarity = 1
			similarityType = ExactSimilarity
		} else if strings.HasPrefix(titleFormatted, targetFormatted) || strings.HasSuffix(titleFormatted, targetFormatted) {
			similarity = 0.8
			similarityType = PrefixOrSuffixSimilarity
		} else {
			similarity = strutil.Similarity(titleFormatted, targetFormatted, metrics.NewJaccard())
			similarityType = NoneSimilarity
		}

		if similarity > bestSimilarity || (similarity == bestSimilarity && strings.Contains(titleFormatted, "(partie 1)")) {
			bestSimilarity = similarity
			bestSimilarityType = similarityType
			closestLoc = url
			closestTitle = titleFormatted
		}
	}

	if bestSimilarityType != "" {
		return bestSimilarityType, closestLoc, closestTitle
	}

	return fmt.Sprintf("%f", bestSimilarity), closestLoc, closestTitle
}
