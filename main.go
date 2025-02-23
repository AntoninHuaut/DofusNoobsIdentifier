package main

import (
	"DofusNoobsIdentifierOffline/boot"
	"DofusNoobsIdentifierOffline/domain"
	"DofusNoobsIdentifierOffline/internal"
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"
)

const (
	questsKey   = "quest"
	dungeonsKey = "dungeon"
)

func main() {
	boot.LoadEnv()
	boot.LoadClient()
	sitemap, err := boot.LoadSitemap()
	if err != nil {
		log.Fatalf("LoadSitemap: %v", err)
	}
	fmt.Println("Total urls:", len(sitemap.Urls))

	quests, err := boot.LoadDofusDbResources[domain.DofusDbQuestLight](questsKey, domain.QuestsFile, boot.HttpClient.RequestDofusApiAllQuests)
	if err != nil {
		log.Fatalf("LoadQuests: %v", err)
	}
	fmt.Println("Total quests:", len(quests.Data))

	dungeons, err := boot.LoadDofusDbResources[domain.DofusDbDungeonLight](dungeonsKey, domain.DungeonsFile, boot.HttpClient.RequestDofusApiAllDungeons)
	if err != nil {
		log.Fatalf("LoadDungeons: %v", err)
	}
	fmt.Println("Total dungeons:", len(dungeons.Data))

	titles, err := boot.LoadTitles(sitemap.Urls)
	if err != nil {
		log.Fatalf("LoadTitles: %v", err)
	}
	fmt.Println("Total titles:", len(titles))

	logs := make(map[string][]string)         // map[resourceType][]logs
	output := make(map[string]map[int]string) // map[resourceType]map[resourceId]dofusNoobsUrl
	initKey := func(key string) {
		logs[key] = make([]string, 0)
		output[key] = make(map[int]string)
	}
	initKey(questsKey)
	initKey(dungeonsKey)

	resolveKey := func(key string, name string, id int) {
		location, similarity, locLog := internal.GetLocationFromTarget(titles, name)

		similarityNb, convertErr := strconv.ParseFloat(similarity, 64)
		if locLog != "" && ((convertErr == nil && similarityNb < 0.9) || similarity == "prefixOrSuffix") {
			logs[key] = append(logs[key], locLog)
		}
		output[key][id] = location
	}
	logKey := func(key string) {
		slices.Sort(logs[key])
		internal.WriteToFile(fmt.Sprintf("logs_%s.txt", key), strings.Join(logs[key], ""), false, true)
	}

	for _, quest := range quests.Data {
		resolveKey(questsKey, quest.Name["fr"], quest.ID)
	}
	logKey(questsKey)

	for _, dungeon := range dungeons.Data {
		resolveKey(dungeonsKey, dungeon.Name["fr"], dungeon.ID)
	}
	logKey(dungeonsKey)

	internal.WriteToFile("mapping.json", output, false, false)
	internal.WriteToFile("mapping_formatted.json", output, true, false)
}
