package boot

import (
	"DofusNoobsIdentifierOffline/domain"
	"DofusNoobsIdentifierOffline/internal"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"os"
)

func LoadTitles(urls []domain.DofusNoobsRemoteSitemapUrl) (map[string]string, error) {
	var titles map[string]string
	cacheFile := internal.GetStorageFilePath(domain.TitlesFile)
	if internal.IsFileExists(cacheFile) {
		fmt.Println("Loading titles from cache")
		bodyFile, err := os.ReadFile(cacheFile)
		if err != nil {
			log.Fatalf("LoadTitles: %v", err)
		}
		if err = json.Unmarshal(bodyFile, &titles); err != nil {
			return nil, err
		}
	}

	var missingUrls []domain.DofusNoobsRemoteSitemapUrl
	initTitles := make(map[string]string)
	for _, url := range urls {
		if _, ok := titles[url.Loc]; !ok {
			missingUrls = append(missingUrls, url)
		} else {
			initTitles[url.Loc] = titles[url.Loc]
		}
	}

	if len(missingUrls) > 0 {
		fmt.Printf("Fetching %d missing titles entries\n", len(missingUrls))
	}

	titles, err := HttpClient.GetPageTitleDofusNoobs(initTitles, missingUrls)
	if err != nil {
		log.Fatalf("GetPageTitleDofusNoobs: %v", err)
	}
	internal.WriteToFile(cacheFile, titles, true, false)
	return titles, nil
}

func LoadSitemap() (*domain.DofusNoobsRemoteSitemap, error) {
	var sitemap *domain.DofusNoobsRemoteSitemap
	cacheFile := internal.GetStorageFilePath(domain.SitemapFile)
	if internal.IsFileExists(cacheFile) {
		fmt.Println("Loading sitemap from cache")
		bodyFile, err := os.ReadFile(cacheFile)
		if err != nil {
			log.Fatalf("LoadSitemap: %v", err)
		}
		if err = json.Unmarshal(bodyFile, &sitemap); err != nil {
			return nil, err
		}
		return sitemap, nil
	}

	rawXml, err := HttpClient.RequestDofusSitemap()
	if err != nil {
		return nil, err
	}

	if err = xml.Unmarshal(rawXml, &sitemap); err != nil {
		return nil, err
	}

	internal.WriteToFile(cacheFile, sitemap, true, false)
	return sitemap, nil
}

func LoadDofusDbResources[T any](resourceType domain.TypeKey, cacheFileName string, getFunc func() (*domain.DofusDbSearchResource[T], error)) (*domain.DofusDbSearchResource[T], error) {
	var resources *domain.DofusDbSearchResource[T]
	cacheFile := internal.GetStorageFilePath(cacheFileName)
	if internal.IsFileExists(cacheFile) {
		fmt.Printf("Loading %s from cache\n", resourceType)
		bodyFile, err := os.ReadFile(cacheFile)
		if err != nil {
			log.Fatalf("%s: %v", resourceType, err)
		}
		if err = json.Unmarshal(bodyFile, &resources); err != nil {
			return nil, err
		}
		return resources, nil
	}

	resources, err := getFunc()
	if err != nil {
		log.Fatalf("RequestDofusApiAllResources (%s): %v", resourceType, err)
	}
	internal.WriteToFile(cacheFile, resources, true, false)
	return resources, nil
}
