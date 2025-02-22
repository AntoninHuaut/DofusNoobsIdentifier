package domain

import (
	"encoding/xml"
	"fmt"
	"strings"
)

type DofusDbQuestLight struct {
	ID   int               `json:"id"`
	Slug map[string]string `json:"slug"`
}

type DofusNoobsRemoteSitemap struct {
	XMLName xml.Name                     `xml:"urlset"`
	Urls    []DofusNoobsRemoteSitemapUrl `xml:"url"`
}

type DofusNoobsRemoteSitemapUrl struct {
	Loc     string `xml:"loc"`
	LastMod string `xml:"lastmod"`

	Slug string // Injected
}

func (u *DofusNoobsRemoteSitemapUrl) FillSlug(dofusNoobsUrl string) {
	u.Slug = strings.ReplaceAll(u.Loc, fmt.Sprintf("%s/", dofusNoobsUrl), "")
	u.Slug = strings.ReplaceAll(u.Slug, ".html", "")
}
