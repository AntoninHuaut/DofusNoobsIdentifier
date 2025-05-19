package internal

import (
	"DofusNoobsIdentifierOffline/domain"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

const (
	dofusDbParallelRequests    = 4
	dofusNoobsParallelRequests = 1
)

var (
	wpTitleRegex = regexp.MustCompile(`<h2 class="wsite-content-title"[^>]*>.*?<strong.*?>(.*?)</strong>.*?</h2>`)
	zwsp         = string(rune(0x200B))
)

type HttpClient interface {
	RequestDofusApiAllDungeons() (*domain.DofusDbSearchResource[domain.DofusDbDungeonLight], error)
	RequestDofusApiAllQuests() (*domain.DofusDbSearchResource[domain.DofusDbQuestLight], error)
	RequestDofusSitemap() ([]byte, error)
	GetPageTitleDofusNoobs(initMaps map[string]string, url []domain.DofusNoobsRemoteSitemapUrl) (map[string]string, error)
}

type httpClient struct {
	dofusApiUrl   string
	dofusNoobsUrl string
}

func NewHttpClient(dofusApiUrl string, dofusNoobsUrl string) HttpClient {
	return &httpClient{dofusApiUrl: dofusApiUrl, dofusNoobsUrl: dofusNoobsUrl}
}

func (h *httpClient) GetPageTitleDofusNoobs(requestTitles map[string]string, urls []domain.DofusNoobsRemoteSitemapUrl) (map[string]string, error) {
	type urlMapped struct {
		URL   string
		Title string
	}

	resultCh := make(chan urlMapped, dofusNoobsParallelRequests)
	errCh := make(chan error, dofusNoobsParallelRequests)

	responseTitles := make(map[string]string)
	for k, v := range requestTitles {
		responseTitles[k] = v
	}

	makeRequest := func(url string) {
		body, err := request(url)
		if err != nil {
			errCh <- err
			return
		}

		titles := wpTitleRegex.FindStringSubmatch(string(body))
		if len(titles) < 2 {
			fmt.Printf("Title not found for %s\n", url)
			resultCh <- urlMapped{URL: url, Title: ""}
			return
		}
		title := html.UnescapeString(titles[1])
		title = strings.ReplaceAll(strings.ReplaceAll(title, "<span>", ""), "</span>", "")
		title = strings.ReplaceAll(title, zwsp, "")
		title = strings.ReplaceAll(strings.TrimSpace(title), zwsp, "")
		resultCh <- urlMapped{URL: url, Title: title}
	}

	currentIndex := 0
	for currentIndex < len(urls) {
		waitCounter := 0
		time.Sleep(500 * time.Millisecond)
		fmt.Printf("Requesting title for %s (%d)\n", urls[currentIndex].Loc, len(urls))
		for i := 0; i < dofusNoobsParallelRequests; i++ {
			if currentIndex+i >= len(urls) {
				break
			}
			waitCounter++
			go makeRequest(urls[currentIndex+i].Loc)
		}

		for i := 0; i < waitCounter; i++ {
			select {
			case mappedRes := <-resultCh:
				currentIndex++
				responseTitles[mappedRes.URL] = mappedRes.Title

				// Write to file to avoid losing data
				jsonMap, err := json.MarshalIndent(responseTitles, "", "  ")
				if err != nil {
					return nil, err
				}
				if err = os.WriteFile(GetStorageFilePath(domain.TitlesFile), jsonMap, 0644); err != nil {
					return nil, err
				}
			case err := <-errCh:
				return nil, err
			}
		}
	}

	close(resultCh)
	close(errCh)

	return responseTitles, nil
}

func (h *httpClient) RequestDofusApiAllDungeons() (*domain.DofusDbSearchResource[domain.DofusDbDungeonLight], error) {
	return RequestDofusApiAllResources[domain.DofusDbDungeonLight](h.dofusApiUrl, "dungeons")
}

func (h *httpClient) RequestDofusApiAllQuests() (*domain.DofusDbSearchResource[domain.DofusDbQuestLight], error) {
	return RequestDofusApiAllResources[domain.DofusDbQuestLight](h.dofusApiUrl, "quests")
}

func RequestDofusApiAllResources[T domain.HasID](dofusApiUrl string, path string) (*domain.DofusDbSearchResource[T], error) {
	resources := &domain.DofusDbSearchResource[T]{
		Limit: 50,
	}

	resultCh := make(chan *domain.DofusDbSearchResource[T], dofusDbParallelRequests)
	errCh := make(chan error, dofusDbParallelRequests)

	makeRequest := func(skip int) {
		fmt.Printf("Requesting resources: skip=%d limit=%d total=%d\n", skip, resources.Limit, resources.Total)
		body, err := request(fmt.Sprintf("%s/%s?$limit=%d&$skip=%d", dofusApiUrl, path, resources.Limit, skip))
		if err != nil {
			errCh <- err
			return
		}

		var subResources *domain.DofusDbSearchResource[T]
		if err = json.Unmarshal(body, &subResources); err != nil {
			errCh <- err
			return
		}

		resultCh <- subResources
	}

	for resources.Total == 0 || resources.Limit+resources.Skip <= resources.Total {
		waitCounter := 0
		for i := 0; i < dofusDbParallelRequests; i++ {
			offsetSkip := resources.Skip + i*resources.Limit
			if resources.Total != 0 && offsetSkip >= resources.Total {
				break
			}

			waitCounter++
			go makeRequest(offsetSkip)
		}

		for i := 0; i < waitCounter; i++ {
			select {
			case subResources := <-resultCh:
				resources.Total = subResources.Total
				resources.Limit = subResources.Limit
				resources.Skip += subResources.Limit
				resources.Data = append(resources.Data, subResources.Data...)
			case err := <-errCh:
				return nil, err
			}
		}
	}

	close(resultCh)
	close(errCh)

	sort.Slice(resources.Data, func(i, j int) bool {
		return resources.Data[i].GetID() < resources.Data[j].GetID()
	})

	return resources, nil
}

func (h *httpClient) RequestDofusSitemap() ([]byte, error) {
	return request(fmt.Sprintf("%s/sitemap.xml", h.dofusNoobsUrl))
}

func request(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK && resp.StatusCode > http.StatusPermanentRedirect {
		return nil, fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}
