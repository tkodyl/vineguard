package main

import (
	"github.com/tkodyl/vineguard/configuration"
	"github.com/tkodyl/vineguard/data/collection/portal_meteo"
	"github.com/tkodyl/vineguard/data/storage/elasticsearch"
	"log"
)

func main() {
	config := configuration.GetConfig()
	collector := portal_meteo.NewCollector(&config)
	fileContent, _ := collector.FetchData()
	records, err := portal_meteo.ToRecords(portal_meteo.DeleteHeaderLine(fileContent))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Numbers of weather records: %d", len(records))
	indexer := elasticsearch.Indexer{Config: config}
	indexer.Do(records)
}
