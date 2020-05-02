package utils

import (
	"errors"
	"time"
)

// GetTimeFromString transform a string to a duration
func GetTimeFromString(input string) (time.Time, error) {
	validTimeFormats := []string{
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
		time.Kitchen,
	}

	for _, layout := range validTimeFormats {
		deadlineTime, err := time.Parse(layout, input)
		if err == nil {
			if layout == time.Kitchen {
				now := time.Now()
				deadlineTime = time.Date(now.Year(),
					now.Month(),
					now.Day(),
					deadlineTime.Hour(),
					deadlineTime.Minute(),
					0,
					0,
					now.Location())

				// If time is before now i refer to that time but the next day
				if deadlineTime.Before(now) {
					deadlineTime = deadlineTime.Add(24 * time.Hour)
				}
			}
			return deadlineTime, nil
		}
	}

	return time.Now(), errors.New("invalid input")
}
