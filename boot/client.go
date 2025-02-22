package boot

import (
	"DofusNoobsIdentifier/internal/client"
)

var (
	httpClient client.HttpClient
)

func LoadClient() {
	httpClient = client.NewHttpClient(DofusApiUrl, DofusNoobsUrl)
}
