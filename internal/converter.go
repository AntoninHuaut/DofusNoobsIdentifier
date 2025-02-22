package internal

import (
	"DofusNoobsIdentifierOffline/domain"
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

var (
	locationCache            = make(map[string]locationValue)
	stripAlignmentOrderTitle = regexp.MustCompile(`(alignement|ordre) \d+ :\s*`)
)

func formatGeneral(general string) string {
	general = strings.ToLower(strings.TrimSpace(general))
	general = strings.ReplaceAll(general, "œ", "oe")
	general = strings.ReplaceAll(general, "\u00a0", " ") // NBSP char
	return general
}

func formatTarget(target string) string {
	if strings.HasSuffix(target, "touriste") {
		target = strings.ReplaceAll(target, "touriste", "touriste/amateur/spécialiste/expert")
	} else if strings.HasSuffix(target, "amateur") {
		target = strings.ReplaceAll(target, "amateur", "touriste/amateur/spécialiste/expert")
	} else if strings.HasSuffix(target, "spécialiste") {
		target = strings.ReplaceAll(target, "spécialiste", "touriste/amateur/spécialiste/expert")
	} else if strings.HasSuffix(target, "expert") {
		target = strings.ReplaceAll(target, "expert", "touriste/amateur/spécialiste/expert")
	}

	if strings.HasPrefix(target, "gare aux krokilles") {
		if strings.HasSuffix(target, "juvéniles") {
			target = strings.ReplaceAll(target, "juvéniles", "juvéniles/novices/matures/vénérables")
		} else if strings.HasSuffix(target, "novices") {
			target = strings.ReplaceAll(target, "novices", "juvéniles/novices/matures/vénérables")
		} else if strings.HasSuffix(target, "matures") {
			target = strings.ReplaceAll(target, "matures", "juvéniles/novices/matures/vénérables")
		} else if strings.HasSuffix(target, "vénérables") {
			target = strings.ReplaceAll(target, "vénérables", "juvéniles/novices/matures/vénérables")
		}
	}

	return target
}

func formatTitle(title string) string {
	title = stripAlignmentOrderTitle.ReplaceAllString(title, "")
	return title
}

func skipQuest(target string) bool {
	return strings.HasPrefix(target, "offrande à ") || strings.HasPrefix(target, "chasse au dopeul ")
}

func GetLocationFromTarget(titles map[string]string, quest domain.DofusDbQuestLight) (string, string, string) {
	targetKey := quest.Name["fr"]
	targetFormatted := formatTarget(formatGeneral(targetKey))
	if skipQuest(targetFormatted) {
		return "[SKIPPED] Offrande ou Chasse au Dopeul", "skipped", ""
	}

	if loc, ok := locationCache[targetKey]; ok {
		return loc.Location, loc.Similarity, ""
	}

	bestSimilarity, closestLoc, closestTitle := findClosestString(targetFormatted, titles)
	locationCache[targetKey] = locationValue{Location: closestLoc, Similarity: bestSimilarity}
	return closestLoc, bestSimilarity, fmt.Sprintf("%s | %60s | %60s (-> %s)\n", bestSimilarity, targetFormatted, closestTitle, closestLoc)
}

func findClosestString(targetFormatted string, titles map[string]string) (string, string, string) {
	bestSimilarityType := ""
	bestSimilarity := .0

	closestLoc := ""
	closestTitle := ""

	for url, title := range titles {
		if strings.HasSuffix(title, "(Dofus Touch)") {
			continue
		}

		var similarity float64
		var similarityType string

		titleFormatted := formatTitle(formatGeneral(title))
		if titleFormatted == targetFormatted {
			similarity = 1
			similarityType = "exact"
		} else if strings.HasPrefix(titleFormatted, targetFormatted) || strings.HasSuffix(titleFormatted, targetFormatted) {
			similarity = 0.8
			similarityType = "prefixOrSuffix"
		} else {
			similarity = strutil.Similarity(titleFormatted, targetFormatted, metrics.NewJaccard())
			similarityType = ""
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
