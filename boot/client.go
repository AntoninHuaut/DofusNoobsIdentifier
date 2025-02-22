package boot

import (
	"DofusNoobsIdentifierOffline/internal"
)

var (
	HttpClient internal.HttpClient
)

func LoadClient() {
	HttpClient = internal.NewHttpClient(DofusApiUrl, DofusNoobsUrl)
}
