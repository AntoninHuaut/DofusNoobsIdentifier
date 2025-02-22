package internal

import (
	"encoding/json"
	"log"
	"os"
)

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func WriteToFile(filename string, data any, formatted bool) {
	var jsonOutput []byte
	var err error
	if formatted {
		jsonOutput, err = json.MarshalIndent(data, "", "  ")
	} else {
		jsonOutput, err = json.Marshal(data)
	}
	if err != nil {
		log.Fatalf("json.Marshal: %v", err)
	}

	if err = os.WriteFile(filename, jsonOutput, 0644); err != nil {
		log.Fatalf("os.WriteFile: %v", err)
	}
}
