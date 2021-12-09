package portal_meteo

import "fmt"

// FormatDate migrate date and time from csv to elasticsearch format
// date and time are combined to provide search abilities and have input for unique document id
func FormatDate(rawDate string, rawTime string) string {
	const elasticsearchDateFormat = "%s-%s-%sT%s:%s:00Z"
	year := rawDate[6:10]
	month := rawDate[3:5]
	day := rawDate[0:2]
	hour := rawTime[0:2]
	minute := rawTime[3:5]
	date := fmt.Sprintf(elasticsearchDateFormat, year, month, day, hour, minute)
	return date
}
