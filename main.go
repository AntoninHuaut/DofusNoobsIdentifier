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

func main() {
	boot.LoadEnv()
	manualConf, err := boot.LoadConf()
	if err != nil {
		log.Fatalf("LoadConf: %v", err)
	}
	fmt.Printf("Manual configuration loaded: %+v\n", manualConf)

	boot.LoadClient()
	sitemap, err := boot.LoadSitemap()
	if err != nil {
		log.Fatalf("LoadSitemap: %v", err)
	}
	fmt.Println("Total urls:", len(sitemap.Urls))

	quests, err := boot.LoadDofusDbResources[domain.DofusDbQuestLight](domain.QuestsKey, domain.QuestsFile, boot.HttpClient.RequestDofusApiAllQuests)
	if err != nil {
		log.Fatalf("LoadQuests: %v", err)
	}
	fmt.Println("Total quests:", len(quests.Data))

	dungeons, err := boot.LoadDofusDbResources[domain.DofusDbDungeonLight](domain.DungeonsKey, domain.DungeonsFile, boot.HttpClient.RequestDofusApiAllDungeons)
	if err != nil {
		log.Fatalf("LoadDungeons: %v", err)
	}
	fmt.Println("Total dungeons:", len(dungeons.Data))

	titles, err := boot.LoadTitles(sitemap.Urls)
	if err != nil {
		log.Fatalf("LoadTitles: %v", err)
	}
	fmt.Println("Total titles:", len(titles))

	logs := make(map[domain.TypeKey][]string)         // map[resourceType][]logs
	output := make(map[domain.TypeKey]map[int]string) // map[resourceType]map[resourceId]dofusNoobsUrl
	initKey := func(key domain.TypeKey) {
		logs[key] = make([]string, 0)
		output[key] = make(map[int]string)
	}
	initKey(domain.QuestsKey)
	initKey(domain.DungeonsKey)

	resolveKey := func(key domain.TypeKey, name string, id int) {
		formattedName := manualConf.RewriteTarget(key, internal.FormatGeneral(name))
		if manualConf.GetIfUnmapped(key, id) {
			return
		}
		if url, ok := manualConf.GetIfRemapped(key, id, formattedName); ok {
			output[key][id] = url
			return
		}

		location, similarity, locLog := internal.GetLocationFromTarget(id, titles, name, formattedName)

		similarityNb, convertErr := strconv.ParseFloat(similarity, 64)
		if locLog != "" && ((convertErr == nil && similarityNb < 0.9) || similarity == "prefixOrSuffix") {
			logs[key] = append(logs[key], locLog)
		}
		output[key][id] = location
	}
	logKey := func(key domain.TypeKey) {
		slices.Sort(logs[key])
		logs[key] = append([]string{internal.FormatLog("Similarity", "DofusDB ID", "DofusDB Title", "DofusNoobs Title", "DofusNoobs URL")}, logs[key]...)
		internal.WriteToFile(internal.GetStorageFilePath(fmt.Sprintf("log_%s.txt", key)), strings.Join(logs[key], ""), false, true)
	}

	for _, quest := range quests.Data {
		resolveKey(domain.QuestsKey, quest.Name["fr"], quest.ID)
	}
	logKey(domain.QuestsKey)

	for _, dungeon := range dungeons.Data {
		resolveKey(domain.DungeonsKey, dungeon.Name["fr"], dungeon.ID)
	}
	logKey(domain.DungeonsKey)

	internal.WriteToFile(internal.GetStorageFilePath("mapping.json"), output, false, false)
	internal.WriteToFile(internal.GetStorageFilePath("mapping_formatted.json"), output, true, false)
}
