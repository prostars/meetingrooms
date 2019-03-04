package util

import (
	"testing"
	"github.com/stretchr/testify/assert"
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

