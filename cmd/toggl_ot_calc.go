package main

import (
	"fmt"
	"log"
	"time"
	"toggl-overtime-calculator/internal/pkg/config"
	"toggl-overtime-calculator/pkg/otcalc"
	"toggl-overtime-calculator/pkg/toggl"
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

	fmt.Printf("%+v\n", report)
	fmt.Println(report.TotalDuration() / 3600000.0)
	fmt.Println(report.TotalGrand / 3600000.0)
}
