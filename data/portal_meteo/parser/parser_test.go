package parser

import "testing"

func TestPrepareLines(t *testing.T) {
	rawCsv := "DATE;TIME;outsideTemperature;outsideHumidity;rainfall;leafWetness1;leafWetness2;solarRadiation\n18/11/2021;00:00;5.1;94;0;0;0;M\n18/11/2021;00:30;5.3;94;0;0;0;M\n18/11/2021;01:00;5.4;94;0.2;0;0;M\n"
	lines := PrepareLines(rawCsv)
	if lines[0] != "18/11/2021;00:00;5.1;94;0;0;0;M" {
		t.Fatalf("Prepared lines incorrect")
	}
}

func TestToRecordsShouldNotIssueError(t *testing.T) {
	rawCsv := []string{"18/11/2021;00:00;5.1;94;0;0;0;M", "18/11/2021;00:30;5.3;94;0;0;0;M", "18/11/2021;01:00;5.4;94;0.2;0;0;M"}
	_, err := ToRecords(rawCsv)
	if err != nil {
		t.Fatalf("error occured during parsing !!")
	}
}

func TestToRecordsShouldIssueError(t *testing.T) {
	incorrectRawCsv := []string{"18/11/2021;00:00;", "18/11/2021;00:30;5.3;94;0;0;0;M", "18/11/2021;01:00;5.4;94;0.2;0;0;M"}
	_, err := ToRecords(incorrectRawCsv)
	if err == nil {
		t.Fatalf("error should occur during parsing !!")
	}
}
