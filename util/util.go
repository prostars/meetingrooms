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
