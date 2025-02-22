package domain

import (
	"encoding/xml"
)

type DofusNoobsRemoteSitemap struct {
	XMLName xml.Name                     `xml:"urlset"`
	Urls    []DofusNoobsRemoteSitemapUrl `xml:"url"`
}

type DofusNoobsRemoteSitemapUrl struct {
	Loc     string `xml:"loc"`
	LastMod string `xml:"lastmod"`
}
