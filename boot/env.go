package boot

import (
	"os"
)

var (
	DofusApiUrl   string
	DofusNoobsUrl string
)

func LoadEnv() {
	DofusApiUrl = getEnvOrDefault("DOFUS_API_URL", "https://api.dofusdb.fr")
	DofusNoobsUrl = getEnvOrDefault("DOFUS_NOOBS_URL", "https://www.dofuspourlesnoobs.com")
}

func getEnvOrDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
