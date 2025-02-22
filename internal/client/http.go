package client

import (
	"fmt"
	"io"
	"net/http"
)

type HttpClient interface {
	RequestDofusApi(path string) ([]byte, error)
	RequestDofusSitemap() ([]byte, error)
}

type httpClient struct {
	dofusApiUrl   string
	dofusNoobsUrl string
}

func NewHttpClient(dofusApiUrl string, dofusNoobsUrl string) HttpClient {
	return &httpClient{dofusApiUrl: dofusApiUrl, dofusNoobsUrl: dofusNoobsUrl}
}

func (h *httpClient) RequestDofusApi(path string) ([]byte, error) {
	return request(fmt.Sprintf("%s%s", h.dofusApiUrl, path))
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
