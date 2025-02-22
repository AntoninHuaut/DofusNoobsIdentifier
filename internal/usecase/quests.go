package usecase

import (
	"DofusNoobsIdentifier/domain"
	"DofusNoobsIdentifier/internal/client"
	"encoding/json"
	"fmt"
	"math"
	"net/url"
	"regexp"
	"strings"
)

var (
	locationCache                 = make(map[string]string)
	regexAlphaNumericalHyphenOnly = regexp.MustCompile(`[^a-zA-Z0-9-]+`)
)

type Quests interface {
	HandleQuests(id int) (dofusNoobsLoc *string, err error)
}

type quests struct {
	httpClient client.HttpClient
	sitemap    *domain.DofusNoobsRemoteSitemap
}

func NewQuests(httpClient client.HttpClient, sitemap *domain.DofusNoobsRemoteSitemap) Quests {
	return &quests{httpClient: httpClient, sitemap: sitemap}
}

func (h *quests) HandleQuests(id int) (*string, error) {
	params := url.Values{}
	params.Add(domain.LangParam, domain.LangParamDefault)
	body, err := h.httpClient.RequestDofusApi(fmt.Sprintf("%s/%d?%s", domain.QuestsPath, id, params.Encode()))
	if err != nil {
		return nil, err
	}

	quest := &domain.DofusDbQuestLight{}
	if err = json.Unmarshal(body, &quest); err != nil {
		return nil, err
	}

	slug := quest.Slug[domain.LangParamDefault]
	if slug == "" {
		return nil, fmt.Errorf("quest slug not found for lang %s", domain.LangParamDefault)
	}

	slug = convertSlugDofusPourLesNoobs(slug)
	loc := getDofusPourLesNoobsLoc(h.sitemap, slug)

	return &loc, nil
}

func getDofusPourLesNoobsLoc(sitemap *domain.DofusNoobsRemoteSitemap, slug string) string {
	if loc, ok := locationCache[slug]; ok {
		return loc
	}

	minDistance := math.MaxInt64
	closestLoc := ""
	for _, s := range sitemap.Urls {
		distance := levenshteinDistance(s.Slug, slug)
		if distance < minDistance {
			minDistance = distance
			closestLoc = s.Loc
		}
	}
	locationCache[slug] = closestLoc
	return closestLoc
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

func convertSlugDofusPourLesNoobs(slug string) string {
	slug = strings.TrimSpace(slug)
	slug = strings.ToLower(slug)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = regexAlphaNumericalHyphenOnly.ReplaceAllString(slug, "")
	return slug
}
