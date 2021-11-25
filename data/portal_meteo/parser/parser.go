package parser

import (
	"log"
	"strconv"
	"strings"
)

type WeatherRecord struct {
	Date        string  `json:"date"`
	Time        string  `json:"time"`
	Temperature float32 `json:"temperature"`
	Humidity    uint8   `json:"humidity"`
	Rain        float32 `json:"rain"`
}

func PrepareLines(rawCsvLines string) []string {
	return strings.Fields(rawCsvLines)[1:]
}

func ToRecords(csvLines []string) ([]WeatherRecord, error) {
	var weatherRecords []WeatherRecord
	var lineItems []string
	for _, line := range csvLines {
		lineItems = strings.Split(line, ";")
		var rec WeatherRecord
		rec.Date = lineItems[0]
		rec.Time = lineItems[1]
		parsedTemperature, err := strconv.ParseFloat(lineItems[2], 32)
		if err != nil {
			log.Println("Cannot parse temperature")
			return nil, err
		}
		rec.Temperature = float32(parsedTemperature)

		parsedHumidity, err := strconv.ParseUint(lineItems[3], 10, 8)
		if err != nil {
			log.Println("Cannot parse humidity")
			return nil, err
		}
		rec.Humidity = uint8(parsedHumidity)

		parsedRain, err := strconv.ParseFloat(lineItems[4], 8)
		if err != nil {
			log.Println("Cannot parse rain")
			return nil, err
		}
		rec.Rain = float32(parsedRain)
		weatherRecords = append(weatherRecords, rec)
	}
	return weatherRecords, nil
}
