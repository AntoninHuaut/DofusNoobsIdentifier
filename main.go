package main

import (
	"DofusNoobsIdentifierOffline/boot"
	"DofusNoobsIdentifierOffline/internal"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
)

func main() {
	boot.LoadEnv()
	boot.LoadClient()
	sitemap, err := boot.LoadSitemap()
	if err != nil {
		log.Fatalf("LoadSitemap: %v", err)
	}
	fmt.Println("Total urls:", len(sitemap.Urls))

	quests, err := boot.LoadQuests()
	if err != nil {
		log.Fatalf("LoadQuests: %v", err)
	}
	fmt.Println("Total quests:", len(quests.Data))

	titles, err := boot.LoadTitles(sitemap.Urls)
	if err != nil {
		log.Fatalf("LoadTitles: %v", err)
	}
	fmt.Println("Total titles:", len(titles))

	var logs []string
	output := make(map[int]string)
	for _, quest := range quests.Data {
		location, similarity, locLog := internal.GetLocationFromTarget(titles, quest)

		similarityNb, convertErr := strconv.ParseFloat(similarity, 64)
		if locLog != "" && ((convertErr == nil && similarityNb < 0.9) || similarity == "contains") {
			logs = append(logs, locLog)
		}
		output[quest.ID] = location
	}

	sort.Strings(logs)
	fmt.Println(strings.Join(logs, ""))
	internal.WriteToFile("logs.txt", strings.Join(logs, ""), false, true)

	internal.WriteToFile("mapping.json", output, false, false)
	internal.WriteToFile("mapping_formatted.json", output, true, false)
}
