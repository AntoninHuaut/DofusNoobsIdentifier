package internal

import (
	"DofusNoobsIdentifierOffline/domain"
	"fmt"
	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"strings"
)

type locationValue struct {
	Location   string
	Similarity float64
}

var (
	locationCache = make(map[string]locationValue)
)

func formatTitle(title string) string {
	title = strings.ToLower(strings.TrimSpace(title))
	title = strings.ReplaceAll(title, "œ", "oe")
	return title
}

func formatTag(quest domain.DofusDbQuestLight, target string) string {
	if strings.HasSuffix(target, "touriste") {
		target = strings.ReplaceAll(target, "touriste", "touriste/amateur/spécialiste/expert")
	} else if strings.HasSuffix(target, "amateur") {
		target = strings.ReplaceAll(target, "amateur", "touriste/amateur/spécialiste/expert")
	} else if strings.HasSuffix(target, "spécialiste") {
		target = strings.ReplaceAll(target, "spécialiste", "touriste/amateur/spécialiste/expert")
	} else if strings.HasSuffix(target, "expert") {
		target = strings.ReplaceAll(target, "expert", "touriste/amateur/spécialiste/expert")
	}

	if quest.IsAlignment() {
		if strings.HasPrefix(target, "on recherche ") {
			return target
		}

		lvl := quest.GetAlignmentLevel()
		if lvl > 0 {
			return fmt.Sprintf("alignement %d : %s", lvl+1, target)
		} else {
			return fmt.Sprintf("alignement : %s", target)
		}
	}
	return target
}

func skipQuest(quest domain.DofusDbQuestLight, target string) bool {
	return strings.HasPrefix(target, "offrande à ") || strings.HasPrefix(target, "chasse au dopeul ")
}

func GetLocationFromTarget(titles map[string]string, quest domain.DofusDbQuestLight) (string, float64, string) {
	targetKey := quest.Name["fr"]
	targetFormatted := formatTag(quest, formatTitle(targetKey))
	if skipQuest(quest, targetFormatted) {
		return "[SKIPPED] Offrande ou Chasse au Dopeul", 0, ""
	}

	if loc, ok := locationCache[targetKey]; ok {
		return loc.Location, loc.Similarity, ""
	}

	bestSimilarity, closestLoc, closestTitle := findClosestString(targetFormatted, titles)
	locationCache[targetKey] = locationValue{Location: closestLoc, Similarity: bestSimilarity}
	return closestLoc, bestSimilarity, fmt.Sprintf("%f | %60s | %60s (-> %s)\n", bestSimilarity, targetFormatted, closestTitle, closestLoc)
}

func findClosestString(targetFormatted string, titles map[string]string) (float64, string, string) {
	bestSimilarity := .0

	closestLoc := ""
	closestTitle := ""

	for url, title := range titles {
		titleFormatted := formatTitle(title)
		if titleFormatted == targetFormatted {
			return 1, url, title
		}

		similarity := strutil.Similarity(titleFormatted, targetFormatted, metrics.NewJaccard())
		if similarity > bestSimilarity {
			bestSimilarity = similarity
			closestLoc = url
			closestTitle = titleFormatted
		}
	}

	return bestSimilarity, closestLoc, closestTitle
}
