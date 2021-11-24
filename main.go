package main

import (
	"fmt"
	"github.com/tkodyl/vineguard/configuration"
	"github.com/tkodyl/vineguard/data/portal_meteo/collector"
)

func main() {
	config := configuration.GetConfig()
	collector := collector.NewCollector(config)
	fileContent, _ := collector.GetDataFromPortalMeteo()
	fmt.Println(fileContent)
}
