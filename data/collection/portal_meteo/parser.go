package portal_meteo

import (
	"log"
	"strconv"
	"strings"
)

type WeatherRecord struct {
	Date        string  `json:"date"`
	Temperature float32 `json:"temperature"`
	Humidity    uint8   `json:"humidity"`
	Rain        float32 `json:"rain"`
}

func DeleteHeaderLine(rawCsvLines string) []string {
	return strings.Fields(rawCsvLines)[1:]
}

func ToRecords(csvLines []string) ([]*WeatherRecord, error) {
	const dateIndex = 0
	const timeIndex = 1
	var weatherRecords []*WeatherRecord
	var lineItems []string
	var prevRecord WeatherRecord
	for _, line := range csvLines {
		lineItems = strings.Split(line, ";")
		var rec WeatherRecord
		rec.Date = FormatDate(lineItems[dateIndex], lineItems[timeIndex])
		temperature, err := getTemperature(lineItems)
		if err != nil {
			log.Printf("Cannot parse temperature line %s: %s, defaulting ...", line, err.Error())
			temperature = prevRecord.Temperature
		}
		rec.Temperature = temperature

		humidity, err := getHumidity(lineItems)
		if err != nil {
			log.Printf("Cannot parse humidity: %s: %s, defaulting ...", line, err.Error())
			humidity = prevRecord.Humidity
		}
		rec.Humidity = humidity

		rain, err := getRain(lineItems)
		if err != nil {
			log.Printf("Cannot parse rain: %s: %s, defaulting ...", line, err.Error())
			rain = prevRecord.Rain
		}
		rec.Rain = rain

		prevRecord = rec
		weatherRecords = append(weatherRecords, &rec)
	}
	return weatherRecords, nil
}

func getTemperature(lineItems []string) (float32, error) {
	const temperatureIndex = 2
	parsedTemperature, err := strconv.ParseFloat(lineItems[temperatureIndex], 32)
	if err != nil {
		return 0, err
	}
	return float32(parsedTemperature), nil
}

func getHumidity(lineItems []string) (uint8, error) {
	const humidityIndex = 3
	parsedHumidity, err := strconv.ParseUint(lineItems[humidityIndex], 10, 8)
	if err != nil {
		return 0, err
	}
	return uint8(parsedHumidity), nil
}

func getRain(lineItems []string) (float32, error) {
	const rainIndex = 4
	parsedRain, err := strconv.ParseFloat(lineItems[rainIndex], 8)
	if err != nil {
		return 0, err
	}
	return float32(parsedRain), nil
}
