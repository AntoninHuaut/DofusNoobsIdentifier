package boot

import (
	"fmt"
	"os"
	"strconv"
)

var (
	DofusApiUrl   string
	DofusNoobsUrl string
	WebServerPort int
)

func LoadEnv() {
	DofusApiUrl = getEnvOrDefault("DOFUS_API_URL", "https://api.dofusdb.fr")
	DofusNoobsUrl = getEnvOrDefault("DOFUS_NOOBS_URL", "https://www.dofuspourlesnoobs.com")
	WebServerPort = getEnvIntOrDefault("WEB_SERVER_PORT", 8080)
}

func getEnvOrDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvIntOrDefault(key string, defaultValue int) int {
	value := getEnvOrDefault(key, fmt.Sprintf("%d", defaultValue))
	valueInt, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return valueInt
}
