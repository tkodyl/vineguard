package main

import (
	"github.com/tkodyl/vineguard/configuration"
	"github.com/tkodyl/vineguard/data/portal_meteo/collector"
	"github.com/tkodyl/vineguard/data/portal_meteo/parser"
	"log"
)

func main() {
	config := configuration.GetConfig()
	collector := collector.NewCollector(&config)
	fileContent, _ := collector.GetDataFromPortalMeteo()
	log.Println(fileContent)
	records, err := parser.ToRecords(parser.PrepareLines(fileContent))
	if err != nil {
		log.Fatal(err)
	}
	log.Println(records)
}
