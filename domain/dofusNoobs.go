package domain

import (
	"encoding/xml"
	"fmt"
	"strings"
)

type DofusNoobsRemoteSitemap struct {
	XMLName xml.Name                     `xml:"urlset"`
	Urls    []DofusNoobsRemoteSitemapUrl `xml:"url"`
}

type DofusNoobsRemoteSitemapUrl struct {
	Loc     string `xml:"loc"`
	LastMod string `xml:"lastmod"`

	FakeSlug string // Injected
}

func (u *DofusNoobsRemoteSitemapUrl) FillSlug(dofusNoobsUrl string) {
	u.FakeSlug = strings.ReplaceAll(u.Loc, fmt.Sprintf("%s/", dofusNoobsUrl), "")
	u.FakeSlug = strings.ReplaceAll(u.FakeSlug, ".html", "")
}
