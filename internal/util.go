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

func WriteToFile(filename string, data any, formatted bool, raw bool) {
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
