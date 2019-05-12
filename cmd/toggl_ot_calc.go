package main

import (
	"fmt"
	"log"
	"time"
	"toggl-overtime-calculator/pkg/toggl"
)

func main() {
	timeUntil, err := time.Parse(toggl.TimeFormat, "2019-04-30")
	if err != nil {
		log.Fatal("Unable to create timeUntil")
	}

	timeSince, err := time.Parse(toggl.TimeFormat, "2019-04-01")
	if err != nil {
		log.Fatal("Unable to create timeUntil")
	}

	request := toggl.ReportRequest{
		UserAgent:   "toggle-overtime-calculator",
		WorkspaceId: 1,
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
