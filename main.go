package main

import (
	"DofusNoobsIdentifierOffline/boot"
	"fmt"
	"log"
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

	//output := make(map[int]string)
	//for _, quest := range quests.Data {
	//	loc := internal.GetLocationFromTarget(titles, quest.Name["fr"])
	//	output[quest.ID] = loc
	//}
	//
	//internal.WriteToFile("output.json", output, false)
	//internal.WriteToFile("output_formatted.json", output, true)
}
