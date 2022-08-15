package helper

import "time"

const (
	datetimeFormat = "2006-01-02 15:04:05"
	dateFormat     = "2006-01-02"
)

func StringToDatetime(date string) (time.Time, error) {
	parseDatetime, err := time.Parse(datetimeFormat, date)
	return parseDatetime, err
}

func StringToDate(date string) (time.Time, error) {
	parseDate, err := time.Parse(dateFormat, date)
	return parseDate, err
}
