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
	HandleQuests(id int, lang string) (dofusNoobsLoc *string, err error)
}

type quests struct {
	httpClient client.HttpClient
	sitemap    *domain.DofusNoobsRemoteSitemap
}

func NewQuests(httpClient client.HttpClient, sitemap *domain.DofusNoobsRemoteSitemap) Quests {
	return &quests{httpClient: httpClient, sitemap: sitemap}
}

func (h *quests) HandleQuests(id int, lang string) (*string, error) {
	params := url.Values{}
	params.Add(domain.LangParam, lang)
	body, err := h.httpClient.RequestDofusApi(fmt.Sprintf("%s/%d?%s", domain.QuestsPath, id, params.Encode()))
	if err != nil {
		return nil, err
	}

	quest := &domain.DofusDbQuestLight{}
	if err = json.Unmarshal(body, &quest); err != nil {
		return nil, err
	}

	slug := quest.Slug[lang]
	if slug == "" {
		return nil, fmt.Errorf("quest slug not found for lang %s", lang)
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

	matrix := make([][]int, lenS1+1)
	for i := range matrix {
		matrix[i] = make([]int, lenS2+1)
	}

	for i := 0; i <= lenS1; i++ {
		matrix[i][0] = i
	}
	for j := 0; j <= lenS2; j++ {
		matrix[0][j] = j
	}

	for i := 1; i <= lenS1; i++ {
		for j := 1; j <= lenS2; j++ {
			cost := 0
			if s1[i-1] != s2[j-1] {
				cost = 1
			}
			matrix[i][j] = int(math.Min(
				float64(matrix[i-1][j]+1),
				math.Min(
					float64(matrix[i][j-1]+1),
					float64(matrix[i-1][j-1]+cost),
				),
			))
		}
	}

	return matrix[lenS1][lenS2]
}

func convertSlugDofusPourLesNoobs(slug string) string {
	slug = strings.TrimSpace(slug)
	slug = strings.ToLower(slug)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = regexAlphaNumericalHyphenOnly.ReplaceAllString(slug, "")
	return slug
}
