package otcalc

import (
	"log"
	"time"
	"toggl-overtime-calculator/pkg/webapi_no"
)

func HoursOvertime(hoursWorked float32, hoursShouldWork float32) float32 {
	return hoursWorked - hoursShouldWork
}

func HoursWorked(msWorked int) float32 {
	return float32(msWorked) / 3600000.0
}

func HoursShouldWork(businessDays int) float32 {
	return float32(businessDays) * 7.5
}

func BusinessDays(since time.Time, until time.Time) int {
	holidays := holidays(since, until)

	totalDays := float32(until.Sub(since) / (24 * time.Hour))
	weekDays := float32(since.Weekday()) - float32(until.Weekday())
	businessDays := int(1 + (totalDays*5-weekDays*2)/7)
	if until.Weekday() == time.Saturday {
		businessDays--
	}
	if since.Weekday() == time.Sunday {
		businessDays--
	}

	return businessDays - holidays
}

func holidays(since time.Time, until time.Time) int {
	if since.After(until) {
		log.Fatal("Since after until")
	}

	var years []int
	theYear := since.Year()
	for theYear <= until.Year() {
		years = append(years, theYear)
		theYear++
	}

	holidays, err := webapi_no.NewHolidaysMulti(years)
	if err != nil {
		log.Fatal(err)
	}

	count := 0
	for _, h := range holidays.Data {
		if ((h.Date.After(since) && h.Date.Before(until)) || h.Date.Equal(since) || h.Date.Equal(until)) && h.Date.Weekday() != time.Saturday && h.Date.Weekday() != time.Sunday {
			count++
		}
	}

	return count
}
