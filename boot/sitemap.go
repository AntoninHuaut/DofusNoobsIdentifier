package boot

import (
	"DofusNoobsIdentifier/domain"
	"encoding/xml"
)

var (
	Sitemap *domain.DofusNoobsRemoteSitemap
)

func LoadSitemap() error {
	rawXml, err := httpClient.RequestDofusSitemap()
	if err != nil {
		return err
	}

	if err = xml.Unmarshal(rawXml, &Sitemap); err != nil {
		return err
	}

	for i := range Sitemap.Urls {
		Sitemap.Urls[i].FillSlug(DofusNoobsUrl)
	}

	return nil
}
