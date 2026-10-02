package Counter

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

func afterNow(now, after time.Time) bool {
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	after = time.Date(after.Year(), after.Month(), after.Day(), 0, 0, 0, 0, after.Location())
	return after.After(now)
}
func isAllowedMDay(date time.Time, targetDay int) bool {
	day := date.Day()

	if targetDay > 0 {
		return day == targetDay
	}
	firstDay := time.Date(date.Year(), date.Month()+1, 1, 0, 0, 0, 0, date.Location())
	if targetDay == -1 {
		return day == firstDay.AddDate(0, 0, -1).Day()
	}
	if targetDay == -2 {
		return day == firstDay.AddDate(0, 0, -2).Day()
	}
	return false
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	date, err := time.Parse(DateFormat, dstart)
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
			if err != nil {
				return "", err
			}
			if avlday < 1 || avlday > 7 {
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
		monthError := errors.New("wrong m-repeat rule")

		if len(ruleParts) < 2 || len(ruleParts) > 3 || ruleParts[1] == "" {
			return "", monthError
		}

		dayStr := strings.Split(ruleParts[1], ",")
		targetDay := make([]int, 0, len(dayStr))

		for _, v := range dayStr {
			day, err := strconv.Atoi(v)
			if err != nil {
				return "", err
			}
			if (day < 1 || day > 31) && day != -1 && day != -2 {
				return "", monthError
			}
			targetDay = append(targetDay, day)
		}

		targetMonth := make(map[int]bool)

		if len(ruleParts) == 3 {

			if ruleParts[2] == "" {
				return "", monthError
			}

			monthStr := strings.Split(ruleParts[2], ",")

			for _, val := range monthStr {
				mth, err := strconv.Atoi(val)
				if err != nil {
					return "", err
				}
				if mth < 1 || mth > 12 {
					return "", monthError
				}
				targetMonth[mth] = true
			}
		}

		for {
			currentMonth := int(date.Month())

			if len(targetMonth) > 0 && !targetMonth[currentMonth] {
				date = time.Date(date.Year(), date.Month()+1, 1, 0, 0, 0, 0, date.Location())
				continue
			}

			dayOk := false

			for _, tDay := range targetDay {
				if isAllowedMDay(date, tDay) {
					dayOk = true
					break
				}
			}

			if dayOk && afterNow(now, date) {
				break
			}

			date = date.AddDate(0, 0, 1)
		}
	default:
		return "", errors.New("unknown rule")
	}
	return date.Format(DateFormat), nil
}
