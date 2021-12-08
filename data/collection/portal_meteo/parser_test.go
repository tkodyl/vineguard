package portal_meteo

import (
	"log"
	"testing"
)

func TestPrepareLines(t *testing.T) {
	const expectedRecordCount = 3
	rawCsv := "DATE;TIME;outsideTemperature;outsideHumidity;rainfall;leafWetness1;leafWetness2;solarRadiation\n18/11/2021;00:00;5.1;94;0;0;0;M\n18/11/2021;00:30;5.3;94;0;0;0;M\n18/11/2021;01:00;5.4;94;0.2;0;0;M\n"
	lines := DeleteHeaderLine(rawCsv)
	if len(lines) == expectedRecordCount && lines[0] != "18/11/2021;00:00;5.1;94;0;0;0;M" {
		t.Fatalf("Prepared lines incorrect")
	}
}

func TestToRecordsShouldNotIssueErrorAndGiveCorrectRecords(t *testing.T) {
	const expectedRecordCount = 3
	date := "2021-11-18T00:00:00Z"
	expectedFirstRecord := WeatherRecord{Date: date, Temperature: 5.1, Humidity: 94, Rain: 0}
	rawCsv := []string{"18/11/2021;00:00;5.1;94;0;0;0;M", "18/11/2021;00:30;5.3;94;0;0;0;M", "18/11/2021;01:00;5.4;94;0.2;0;0;M"}
	result, err := ToRecords(rawCsv)
	if err != nil {
		t.Fatalf("error occured during parsing !!")
	}
	if len(result) != expectedRecordCount {
		t.Fatalf("Inproper lenght of result slice")
	}
	if *result[0] != expectedFirstRecord && result[1].Date == date && result[2].Date == date {
		t.Fatalf("incorrectly decoded data")
	}
}

func TestToRecordsWhenSecondRecordsFieldsAreEmpty(t *testing.T) {
	secondWeatherRecord := WeatherRecord{Date: "2021-11-18T00:30:00Z", Temperature: float32(5.3), Humidity: uint8(94), Rain: float32(0)}
	incorrectRawCsv := []string{"18/11/2021;00:00;5.3;94;0;0;0;M", "18/11/2021;00:30;;;;;;", "18/11/2021;01:00;5.4;94;0.2;0;0;M"}
	result, _ := ToRecords(incorrectRawCsv)
	if *result[1] != secondWeatherRecord {
		t.Fatalf("Mismatch of second object")
	}
}

func TestToRecordsWhenSecondFieldsAreIncorrect(t *testing.T) {
	expectedSecondRecord := WeatherRecord{Date: "2021-11-18T00:30:00Z", Temperature: float32(5.3), Humidity: uint8(94), Rain: float32(0)}
	incorrectRawCsv := []string{"18/11/2021;00:00;5.3;94;0;0;0;M", "18/11/2021;00:30;M;M;M;M;M;M", "18/11/2021;01:00;5.4;94;0.2;0;0;M"}
	result, _ := ToRecords(incorrectRawCsv)
	if *result[1] != expectedSecondRecord {
		t.Fatalf("Mismatch of second object")
	}
}

func TestToRecordsShouldDefaultWhenIssueInFirstRecord(t *testing.T) {
	expectedFirstRecord := WeatherRecord{Date: "2021-11-18T00:00:00Z", Temperature: float32(0), Humidity: uint8(0), Rain: float32(0)}
	incorrectRawCsv := []string{"18/11/2021;00:00;;;;;;", "18/11/2021;00:30;5.3;94;0;0;0;M", "18/11/2021;01:00;5.4;94;0.2;0;0;M"}
	records, _ := ToRecords(incorrectRawCsv)
	log.Println(records)
	if *records[0] != expectedFirstRecord {
		t.Fatalf("First record should hold zeros in temperature, humidity and rain when missing")
	}
}
