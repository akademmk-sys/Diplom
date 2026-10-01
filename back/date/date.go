package date

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

func afterNow(now, after time.Time) bool {
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	after = time.Date(after.Year(), after.Month(), after.Day(), 0, 0, 0, 0, after.Location())
	return after.After(now)
}
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", err
	}
	if repeat == "" {
		return "", errors.New("repeat rule is empty")
	}

	ruleParts := strings.Split(repeat, " ")

	switch ruleParts[0] {
	case "d":
		if len(ruleParts) != 2 {
			return "", errors.New("wrong d-repeat rule")
		}
		nDays, err := strconv.Atoi(ruleParts[1])
		if err != nil {
			return "", err
		}
		if nDays < 1 || nDays > 400 {
			return "", errors.New("invalid day count parametr")
		}
		for !afterNow(now, date) {
			date = date.AddDate(0, 0, nDays)
		}
	case "y":
		if len(ruleParts) != 1 {
			return "", errors.New("wrong y-repeat rule")
		}
		for !afterNow(now, date) {
			date = date.AddDate(1, 0, 0)
		}
	case "w":
		weekError := errors.New("wrong w-repeat rule")
		if len(ruleParts) != 2 {
			return "", weekError
		}
		if ruleParts[1] == "" {
			return "", weekError
		}
		daysStr := strings.Split(ruleParts[1], ",")
		avlDays := make(map[int]bool)
		for _, wd := range daysStr {
			avlday, err := strconv.Atoi(wd)
			if err != nil || avlday < 1 || avlday > 7 {
				return "", weekError
			}
			avlDays[avlday] = true
		}
		for {
			dayOfWeek := int(date.Weekday())
			if dayOfWeek == 0 {
				dayOfWeek = 7
			}
			if avlDays[dayOfWeek] && afterNow(now, date) {
				break
			}
			date = date.AddDate(0, 0, 1)
		}
	case "m":
		if len(ruleParts) < 2 || len(ruleParts) > 3 {
			return "", errors.New("wrong m-repeat rule")
		}
		switch len(ruleParts) {
		case 2:
			alwdDays := strings.Split(ruleParts[1], ",")
		case 3:
			alwdDays := strings.Split(ruleParts[1], ",")
		}
	}
	return date.Format(dateFormat), nil
}
