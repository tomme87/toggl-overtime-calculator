package main

import (
	"fmt"
	"github.com/tomme87/toggl-overtime-calculator/internal/pkg/config"
	"github.com/tomme87/toggl-overtime-calculator/pkg/otcalc"
	"github.com/tomme87/toggl-overtime-calculator/pkg/toggl"
	"log"
	"time"
)

func main() {
	err := config.C.Init()
	if err != nil {
		log.Fatal(err)
	}

	until := config.C.Options.Until
	if until == "" {
		log.Fatal("Missing date until (use --until)")
	}

	since := config.C.Options.Since
	if since == "" {
		log.Fatal("Missing date since (use --since)")
	}

	timeUntil, err := time.Parse(toggl.TimeFormat, until)
	if err != nil {
		log.Fatal("Unable to create timeUntil use date format yyyy-mm-dd")
	}

	timeSince, err := time.Parse(toggl.TimeFormat, since)
	if err != nil {
		log.Fatal("Unable to create timeSince use date format yyyy-mm-dd")
	}

	request := toggl.ReportRequest{
		UserAgent:   "toggle-overtime-calculator",
		WorkspaceId: config.C.Toggl.WorkspaceId,
		Until:       timeUntil,
		Since:       timeSince,
		Page:        1,
	}

	report, err := toggl.NewDetailReport(request)
	if err != nil {
		log.Fatal(err)
	}

	hoursWorked := otcalc.HoursWorked(report.TotalGrand)

	businessDays := otcalc.BusinessDays(timeSince, timeUntil)
	hoursShouldWork := otcalc.HoursShouldWork(businessDays, timeSince, timeUntil)
	hoursOvertime := otcalc.HoursOvertime(hoursWorked, hoursShouldWork)

	fmt.Printf("Businessdays between %s and %s: %d\n", timeSince.Format(toggl.TimeFormat), timeUntil.Format(toggl.TimeFormat), businessDays)

	fmt.Printf("Hours worked: %.2f\n", hoursWorked)
	fmt.Printf("Hours a normal person would work: %.2f\n", hoursShouldWork)
	fmt.Printf("Hours Overtime: %.2f\n", hoursOvertime)

}
