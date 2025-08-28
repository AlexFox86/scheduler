package tasks

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const dateFmt = "20060102"

func repeatDays(now time.Time, parsedDate time.Time, parts []string) (time.Time, error) {
	if len(parts) == 1 {
		return time.Time{}, fmt.Errorf("missing day count in 'repeat' parameter")
	}

	days, err := strconv.Atoi(parts[1])
	if err != nil || days <= 0 || days > 400 {
		return time.Time{}, fmt.Errorf("invalid day count in 'repeat' parameter")
	}

	for {
		parsedDate = parsedDate.AddDate(0, 0, days)
		if parsedDate.After(now) {
			return parsedDate, nil
		}
	}
}

func repeatYears(now time.Time, parsedDate time.Time) (time.Time, error) {
	for {
		parsedDate = parsedDate.AddDate(1, 0, 0)
		if parsedDate.After(now) {
			return parsedDate, nil
		}
	}
}

// NextDate calculates the next date for the task according to the specified rule
func NextDate(now time.Time, startDate string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("'repeat' parameter not found")
	}

	parsedDate, err := time.Parse(dateFmt, startDate)
	if err != nil {
		return "", fmt.Errorf("invalid 'startDate' format")
	}

	parts := strings.Split(repeat, " ")
	if len(parts) == 0 {
		return "", fmt.Errorf("invalid 'repeat' format")
	}

	switch parts[0] {
	case "d":
		parsedDate, err = repeatDays(now, parsedDate, parts)
		if err != nil {
			return "", err
		}

	case "y":
		parsedDate, err = repeatYears(now, parsedDate)
		if err != nil {
			return "", err
		}

	default:
		return "", fmt.Errorf("unsupported 'repeat' type")
	}

	return parsedDate.Format(dateFmt), nil
}
