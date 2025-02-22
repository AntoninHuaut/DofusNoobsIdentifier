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

	//loc := internal.GetLocationFromTarget(sitemap, quest.Name["fr"])
	//if quest.ID == 1349 {
	//	fmt.Println("Loc:", loc)
	//	fmt.Println("Name:", quest.Name["fr"])
	//}
	//output[quest.ID] = loc
	//}

	//writeToFile("output.json", output, false)
	//writeToFile("output_formatted.json", output, true)
}
