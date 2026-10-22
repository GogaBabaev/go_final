package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

func afterNow(date, now time.Time) bool {

	dateStr := date.Format(dateFormat)
	nowStr := now.Format(dateFormat)
	return dateStr > nowStr

}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	now := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var nowTime time.Time
	var err error

	if now == "" {
		nowTime = time.Now()
	} else {
		nowTime, err = time.Parse(dateFormat, now)
	}

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	result, err := NextDate(nowTime, date, repeat)

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte(result))

}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	date, err := time.Parse(dateFormat, dstart)

	if err != nil {
		return "", err
	}

	if repeat == "" {

		return "", errors.New("rule of repeat not found")

	}

	parts := strings.Split(repeat, " ")

	switch parts[0] {

	case "d":
		if len(parts) != 2 {
			return "", errors.New("invalid repeat format")
		}

		interval, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", err
		}

		if interval < 1 || interval > 400 {
			return "", errors.New("invalid interval")
		}

		for {
			date = date.AddDate(0, 0, interval)

			if date.After(now) {
				break
			}

		}

		return date.Format(dateFormat), nil

	case "y":

		for {
			date = date.AddDate(1, 0, 0)

			if date.After(now) {
				break
			}

		}

		return date.Format(dateFormat), nil

	case "w":

		if len(parts) != 2 {
			return "", errors.New("invalid repeat format")
		}

		weekDays := strings.Split(parts[1], ",")

		var allowedDays [8]bool

		for _, dayStr := range weekDays {
			day, err := strconv.Atoi(dayStr)

			if err != nil {
				return "", err
			}

			if day < 1 || day > 7 {
				return "", errors.New("invalid weekday")
			}

			allowedDays[day] = true
		}

		for {
			date = date.AddDate(0, 0, 1)

			weekday := int(date.Weekday())

			if weekday == 0 {
				weekday = 7
			}

			if allowedDays[weekday] && date.After(now) {
				break
			}
		}

		return date.Format(dateFormat), nil

	case "m":

		if len(parts) < 2 || len(parts) > 3 {

			return "", errors.New("invalid repeat format")

		}

		monthDays := strings.Split(parts[1], ",")

		var allowedDays [32]bool

		var allowedMonths [13]bool

		var lastDay bool

		var beforeLastDay bool

		for _, dayStr := range monthDays {
			day, err := strconv.Atoi(dayStr)

			if err != nil {
				return "", err
			}

			switch day {
			case -1:
				lastDay = true

			case -2:
				beforeLastDay = true

			default:
				if day < 1 || day > 31 {
					return "", errors.New("invalid month day")
				}

				allowedDays[day] = true
			}

		}

		if len(parts) == 3 {
			months := strings.Split(parts[2], ",")

			for _, monthsStr := range months {
				month, err := strconv.Atoi(monthsStr)

				if err != nil {
					return "", err
				}

				if month < 1 || month > 12 {
					return "", errors.New("invalid month")
				}
				allowedMonths[month] = true
			}

		} else {
			for i := 1; i <= 12; i++ {
				allowedMonths[i] = true
			}
		}

		for {
			date = date.AddDate(0, 0, 1)

			day := date.Day()

			month := int(date.Month())

			nextDay := date.AddDate(0, 0, 1)

			isLastDay := date.Month() != nextDay.Month()

			dayAfterNext := date.AddDate(0, 0, 2)

			isBeforeLastDay := date.Month() == nextDay.Month() &&

				date.Month() != dayAfterNext.Month()

			if allowedMonths[month] &&
				(allowedDays[day] ||

					(lastDay && isLastDay) ||

					(beforeLastDay && isBeforeLastDay)) &&

				date.After(now) {

				break

			}
		}

		return date.Format(dateFormat), nil

	default:

		return "", errors.New("unsupported repeat rule")

	}

}
