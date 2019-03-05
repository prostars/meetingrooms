package util

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
)

func TestBuildBookingTimes(t *testing.T) {
	times := BuildBookingTimes("1000", "1230")
	assert.Equal(t, times, []string{"1000", "1030", "1100", "1130", "1200", "1230"})
}

func TestSplitStartTimeAndUserName(t *testing.T) {
	time, name := SplitStartTimeAndUserName("1500:prostars")
	assert.Equal(t, time, "1500")
	assert.Equal(t, name, "prostars")
}

func TestGetNextWeekDate(t *testing.T) {
	nextWeek := GetNextWeekDate("20190304")
	assert.Equal(t, nextWeek, "20190311")
}

func TestIsValidRepeatCount(t *testing.T) {
	assert.Equal(t, IsValidRepeatCount(10), true)
	assert.Equal(t, IsValidRepeatCount(-10), false)
	assert.Equal(t, IsValidRepeatCount(-0), false)
}

func TestIsValidTime(t *testing.T) {
	assert.Equal(t, IsValidTime("0800"), true)
	assert.Equal(t, IsValidTime("2330"), true)
	assert.Equal(t, IsValidTime("0730"), false)
	assert.Equal(t, IsValidTime("1210"), false)
	assert.Equal(t, IsValidTime("2340"), false)
	assert.Equal(t, IsValidTime("2390"), false)
	assert.Equal(t, IsValidTime("2400"), false)
}

func TestIsValidDate(t *testing.T) {
	const layOut = "20060102"
	now := time.Now().Format(layOut)
	oldDate, _ := time.Parse(layOut, "20101205")
	assert.Equal(t, IsValidDate(now), true)
	assert.Equal(t, IsValidDate(GetNextWeekDate(now)), true)
	assert.Equal(t, IsValidDate(oldDate.Format(layOut)), false)
}