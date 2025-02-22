package main

import (
	"DofusNoobsIdentifier/boot"
	"fmt"
	"log"
)

func main() {
	boot.LoadEnv()
	boot.LoadClient()
	if err := boot.LoadSitemap(); err != nil {
		log.Fatalf("LoadSitemap: %v", err)
	}

	webServer := boot.LoadWebserver()
	if err := webServer.Run(fmt.Sprintf(":%d", boot.WebServerPort)); err != nil {
		log.Fatalf("WebServer run: %v", err)
	}
}
