package main

import (
	"fmt"
	"github.com/tkodyl/vineguard/configuration"
	"github.com/tkodyl/vineguard/data/portal_meteo/collector"
	"log"
	"strings"
)

func main() {
	config := configuration.GetConfig()
	collector := collector.NewCollector(&config)
	fileContent, _ := collector.GetDataFromPortalMeteo()
	log.Println(fileContent)
	fmt.Println(strings.Fields(fileContent)[1:])
}
