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
	RequestDofusApiAllQuests() (*domain.DofusDbSearchQuest, error)
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

func (h *httpClient) GetPageTitleDofusNoobs(resTitles map[string]string, urls []domain.DofusNoobsRemoteSitemapUrl) (map[string]string, error) {
	type urlMapped struct {
		URL   string
		Title string
	}

	resultCh := make(chan urlMapped, dofusNoobsParallelRequests)
	errCh := make(chan error, dofusNoobsParallelRequests)

	makeRequest := func(url string) {
		body, err := request(url)
		if err != nil {
			errCh <- err
			return
		}

		titles := wpTitleRegex.FindStringSubmatch(string(body))
		if len(titles) < 2 {
			fmt.Printf("Title not found: %s\n", url)
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
		fmt.Println("Requesting titles:", urls[currentIndex], len(urls))
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
				resTitles[mappedRes.URL] = mappedRes.Title

				jsonMap, err := json.MarshalIndent(resTitles, "", "  ")
				if err != nil {
					return nil, err
				}
				if err = os.WriteFile(domain.TitlesTmpFile, jsonMap, 0644); err != nil {
					return nil, err
				}
			case err := <-errCh:
				return nil, err
			}
		}
	}

	close(resultCh)
	close(errCh)

	return resTitles, nil
}

func (h *httpClient) RequestDofusApiAllQuests() (*domain.DofusDbSearchQuest, error) {
	quests := &domain.DofusDbSearchQuest{
		Limit: 50,
	}

	resultCh := make(chan *domain.DofusDbSearchQuest, dofusDbParallelRequests)
	errCh := make(chan error, dofusDbParallelRequests)

	makeRequest := func(skip int) {
		fmt.Printf("Requesting quests: skip=%d limit=%d total=%d\n", skip, quests.Limit, quests.Total)
		body, err := request(fmt.Sprintf("%s/quests?$limit=%d&$skip=%d", h.dofusApiUrl, quests.Limit, skip))
		if err != nil {
			errCh <- err
			return
		}

		var subQuests *domain.DofusDbSearchQuest
		if err = json.Unmarshal(body, &subQuests); err != nil {
			errCh <- err
			return
		}

		resultCh <- subQuests
	}

	for quests.Total == 0 || quests.Limit+quests.Skip <= quests.Total {
		waitCounter := 0
		for i := 0; i < dofusDbParallelRequests; i++ {
			offsetSkip := quests.Skip + i*quests.Limit
			if quests.Total != 0 && offsetSkip >= quests.Total {
				break
			}

			waitCounter++
			go makeRequest(offsetSkip)
		}

		for i := 0; i < waitCounter; i++ {
			select {
			case subQuests := <-resultCh:
				quests.Total = subQuests.Total
				quests.Limit = subQuests.Limit
				quests.Skip += subQuests.Limit
				quests.Data = append(quests.Data, subQuests.Data...)
			case err := <-errCh:
				return nil, err
			}
		}
	}

	close(resultCh)
	close(errCh)

	return quests, nil
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
