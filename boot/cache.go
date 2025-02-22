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
	if internal.FileExists(domain.TitlesFile) {
		fmt.Println("Loading titlesFile from cache")
		bodyFile, err := os.ReadFile(domain.TitlesFile)
		if err != nil {
			log.Fatalf("LoadTitles: %v", err)
		}
		if err = json.Unmarshal(bodyFile, &titles); err != nil {
			return nil, err
		}
	}

	var urlsFiltered []domain.DofusNoobsRemoteSitemapUrl
	initMaps := make(map[string]string)
	for _, url := range urls {
		if _, ok := titles[url.Loc]; !ok {
			urlsFiltered = append(urlsFiltered, url)
		} else {
			initMaps[url.Loc] = titles[url.Loc]
		}
	}

	titles, err := HttpClient.GetPageTitleDofusNoobs(initMaps, urlsFiltered)
	if err != nil {
		log.Fatalf("GetPageTitleDofusNoobs: %v", err)
	}
	internal.WriteToFile(domain.TitlesFile, titles, true, false)
	return titles, nil
}

func LoadSitemap() (*domain.DofusNoobsRemoteSitemap, error) {
	var sitemap *domain.DofusNoobsRemoteSitemap
	if internal.FileExists(domain.SitemapFile) {
		fmt.Println("Loading sitemap from cache")
		bodyFile, err := os.ReadFile(domain.SitemapFile)
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

	internal.WriteToFile(domain.SitemapFile, sitemap, true, false)
	return sitemap, nil
}

func LoadQuests() (*domain.DofusDbSearchQuest, error) {
	var quests *domain.DofusDbSearchQuest
	if internal.FileExists(domain.QuestsFile) {
		fmt.Println("Loading questsFile from cache")
		bodyFile, err := os.ReadFile(domain.QuestsFile)
		if err != nil {
			log.Fatalf("LoadQuests: %v", err)
		}
		if err = json.Unmarshal(bodyFile, &quests); err != nil {
			return nil, err
		}
		return quests, nil
	}

	quests, err := HttpClient.RequestDofusApiAllQuests()
	if err != nil {
		log.Fatalf("RequestDofusApiAllQuests: %v", err)
	}
	internal.WriteToFile(domain.QuestsFile, quests, true, false)
	return quests, nil
}
