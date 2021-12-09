package portal_meteo

import "testing"

func TestCorrectDateParsing(t *testing.T) {
	rawDate := "02/12/2021"
	rawTime := "00:00"
	expectedFormattedDate := "2021-12-02T00:00:00Z"
	result := FormatDate(rawDate, rawTime)
	if result != expectedFormattedDate {
		t.Fatalf("Incorrect date parsing result: %s", result)
	}

}
