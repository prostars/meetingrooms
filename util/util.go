package util

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func BuildBookingTimes(startTime string, endTime string) []string {
	start := 0
	end := 0

	if v, err := strconv.Atoi(startTime); err == nil {
		start = v
	} else {
		return nil
	}
	if v, err := strconv.Atoi(endTime); err == nil {
		end = v
	} else {
		return nil
	}

	if start >= end {
		return nil
	}

	p := start
	times := []string{fmt.Sprintf("%04d", p)}

	for p < end {
		if p % 100 == 0 {
			p += 30
		} else {
			p += 70
		}
		times = append(times, fmt.Sprintf("%04d", p))
	}
	return times
}

func SplitStartTimeAndUserName(value string) (string, string) {
	s := strings.Split(value, ":")
	if len(s) != 2 {
		return "", ""
	}
	return s[0], s[1]
}

func GetNextWeekDate(date string) string {
	const layOut = "20060102"
	const sevenDays = time.Hour * 24 * 7

	dateStamp, err := time.Parse(layOut, date)
	if err != nil {
		return ""
	}
	afterWeek := dateStamp.Add(sevenDays)
	return afterWeek.Format("20060102")
}

func IsValidDate(date string) bool {
	if len(date) != 8 {
		return false
	}
	const layOut = "20060102"
	if d, err := time.Parse(layOut, date); err == nil {
		if int(time.Now().Sub(d).Hours() / 24) <= 0 {
			return true
		}
	}
	return false
}

func IsValidTime(time string) bool {
	if len(time) != 4 {
		return false
	}
	value, err := strconv.Atoi(time)
	if err != nil {
		return false
	}
	if 2330 < value || value < 800 {
		return false
	}
	v := value % 100
	if v == 0 || v == 30 {
		return true
	}
	return false
}

func IsValidRepeatCount(count int) bool {
	if 10 < count || count < 1 {
		return false
	}
	return true
}