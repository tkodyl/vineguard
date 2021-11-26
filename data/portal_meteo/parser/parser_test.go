package parser

import "testing"

func TestPrepareLines(t *testing.T) {
	const expectedRecordCount = 3
	rawCsv := "DATE;TIME;outsideTemperature;outsideHumidity;rainfall;leafWetness1;leafWetness2;solarRadiation\n18/11/2021;00:00;5.1;94;0;0;0;M\n18/11/2021;00:30;5.3;94;0;0;0;M\n18/11/2021;01:00;5.4;94;0.2;0;0;M\n"
	lines := PrepareLines(rawCsv)
	if len(lines) == expectedRecordCount && lines[0] != "18/11/2021;00:00;5.1;94;0;0;0;M" {
		t.Fatalf("Prepared lines incorrect")
	}
}

func TestToRecordsShouldNotIssueErrorAndGiveCorrectRecords(t *testing.T) {
	const expectedRecordCount = 3
	date := "18/11/2021"
	expectedFirstRecord := WeatherRecord{Date: date, Time: "00:00", Temperature: 5.1, Humidity: 94, Rain: 0}
	rawCsv := []string{"18/11/2021;00:00;5.1;94;0;0;0;M", "18/11/2021;00:30;5.3;94;0;0;0;M", "18/11/2021;01:00;5.4;94;0.2;0;0;M"}
	result, err := ToRecords(rawCsv)
	if err != nil {
		t.Fatalf("error occured during parsing !!")
	}
	if len(result) == expectedRecordCount && result[0] != expectedFirstRecord && result[1].Date == date && result[2].Date == date {
		t.Fatalf("incorrectly decoded data")
	}
}

//func TestToRecordsShouldIssueError(t *testing.T) {
//	incorrectRawCsv := []string{"18/11/2021;00:00;", "18/11/2021;00:30;5.3;94;0;0;0;M", "18/11/2021;01:00;5.4;94;0.2;0;0;M"}
//	_, err := ToRecords(incorrectRawCsv)
//	if err == nil {
//		t.Fatalf("error should occur during parsing !!")
//	}
//}

// TODO: test for incorrect first record
// TODO: test for incorrect subsequent record
