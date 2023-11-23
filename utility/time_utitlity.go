package utility

import (
	"errors"
	"strings"
	"time"
)

func ReturnCurrentTime() time.Time {
	return time.Now()
}

func FormatTime(dateParams string) (outputTime time.Time, err error) {

	var formattedDate time.Time

	if strings.TrimSpace(dateParams) == "" {
		return formattedDate, errors.New("empty string date time")
	}

	if len(dateParams) > 0 {
		dateTime, err := time.Parse("2006-01-02", dateParams)
		if err != nil {
			return formattedDate, errors.New("date format should be YYYY-MM-DD")
		}
		formattedDate = dateTime
	}
	return formattedDate, nil
}
