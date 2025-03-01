package boot

import (
	"DofusNoobsIdentifierOffline/domain"
	"gopkg.in/yaml.v3"
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

func LoadConf() (*domain.ManualFile, error) {
	var manualConf *domain.ManualFile
	file, err := os.ReadFile("./domain/manual.yaml")
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(file, &manualConf)
	if err != nil {
		return nil, err
	}

	return manualConf, nil
}

func getEnvOrDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
