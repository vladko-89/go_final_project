package date_pkg

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"final-project/pkg/consts"
)

func NextDate(nowIn time.Time, date string, repeat string) (string, error) {

	var currentDate time.Time
	var err error

	now, _ := time.Parse(consts.FormatDate, nowIn.Format(consts.FormatDate))

	// Если дата не указана, подставляем текущую дату
	if date == "" {
		currentDate = now
	} else {
		currentDate, err = time.Parse(consts.FormatDate, date)
		if err != nil {
			return "", fmt.Errorf("invalid date format: %v", err)
		}
	}

	// Если правило повторения не поддерживается
	switch {
	case strings.HasPrefix(repeat, "d "):

		daysStr := strings.TrimPrefix(repeat, "d ")
		days, err := strconv.Atoi(daysStr)

		if err != nil || days < 1 || days > 400 {
			return "", errors.New("invalid day interval")
		}
		nextDate := currentDate

		if nextDate.Before(now) {
			for nextDate.Before(now) || nextDate.Equal(now) {
				nextDate = nextDate.AddDate(0, 0, days)
			}
		} else if nextDate.Equal(now) {

			nextDate = now
		} else {

			nextDate = nextDate.AddDate(0, 0, days)
		}

		if nextDate.Month() == time.February && nextDate.Day() == 29 {

			if !isLeapYear(nextDate.Year()) {
				nextDate = time.Date(nextDate.Year(), time.March, 1, 0, 0, 0, 0, nextDate.Location())
			}
		}

		return nextDate.Format(consts.FormatDate), nil

	case repeat == "y":

		nextDate := currentDate

		if nextDate.Year() < now.Year() {
			yearsDiff := now.Year() - nextDate.Year()
			fmt.Println(yearsDiff)
			nextDate = nextDate.AddDate(yearsDiff, 0, 0)
			return nextDate.Format(consts.FormatDate), nil
		}

		if nextDate.Before(now) {
			for nextDate.Before(now) {
				nextDate = nextDate.AddDate(1, 0, 0)
			}
		}
		if !nextDate.After(now) {
			nextDate = nextDate.AddDate(1, 0, 0) // Увеличиваем на 1 год
		} else {
			nextDate = nextDate.AddDate(1, 0, 0) // Увеличиваем на 1 год
		}

		if nextDate.Month() == time.February && nextDate.Day() == 29 {

			if !isLeapYear(nextDate.Year()) {
				nextDate = time.Date(nextDate.Year(), time.March, 1, 0, 0, 0, 0, nextDate.Location())
			}
		}

		return nextDate.Format(consts.FormatDate), nil

	default:
		return "", fmt.Errorf("правило повторения указано в неправильном формате - %s", repeat)
	}
}

func isLeapYear(year int) bool {
	if year%4 == 0 {
		if year%100 == 0 {
			return year%400 == 0

		}
		return true
	}
	return false
}
