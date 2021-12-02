package main

import (
	"encoding/json"
	"github.com/tkodyl/vineguard/configuration"
	"github.com/tkodyl/vineguard/data/collection/pm"
	"log"
)

func main() {
	config := configuration.GetConfig()
	collector := pm.NewCollector(&config)
	fileContent, _ := collector.FetchData()
	log.Println(fileContent)
	records, err := pm.ToRecords(pm.PrepareLines(fileContent))
	if err != nil {
		log.Fatal(err)
	}
	log.Println(records)
	jsonData, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(jsonData))
}
