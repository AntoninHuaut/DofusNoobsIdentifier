package internal

import (
	"DofusNoobsIdentifierOffline/domain"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func GetStorageFilePath(fileName string) string {
	return fmt.Sprintf("%s/%s", domain.StorageDir, fileName)
}

func createStorageDirIfNotExists() {
	if _, statsErr := os.Stat(domain.StorageDir); os.IsNotExist(statsErr) {
		if mkdirErr := os.Mkdir(domain.StorageDir, 0755); mkdirErr != nil {
			log.Fatalf("Failed to create storage directory: %v", mkdirErr)
		}
	}
}

func IsFileExists(fileName string) bool {
	_, err := os.Stat(fileName)
	return !os.IsNotExist(err)
}

func WriteToFile(filename string, data any, formatted bool, raw bool) {
	createStorageDirIfNotExists()

	var outputByte []byte

	if raw {
		outputByte = []byte(data.(string))
	} else {
		var err error
		if formatted {
			outputByte, err = json.MarshalIndent(data, "", "  ")
		} else {
			outputByte, err = json.Marshal(data)
		}
		if err != nil {
			log.Fatalf("json.Marshal: %v", err)
		}
	}

	if err := os.WriteFile(filename, outputByte, 0644); err != nil {
		log.Fatalf("os.WriteFile: %v", err)
	}
}
