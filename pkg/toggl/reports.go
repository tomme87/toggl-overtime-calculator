package toggl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const TimeFormat = "2006-01-02"

const detailedReportUrl = "https://toggl.com/reports/api/v2/details"

type ReportRequest struct {
	UserAgent   string
	WorkspaceId int
	Since       time.Time
	Until       time.Time
	Page        int
}

func (rr *ReportRequest) Url() string {
	u, err := url.Parse(detailedReportUrl)
	if err != nil {
		log.Fatal("Unable to create URL")
	}

	q := u.Query()
	q.Add("user_agent", rr.UserAgent)
	q.Add("workspace_id", strconv.Itoa(rr.WorkspaceId))
	q.Add("since", rr.Since.Format(TimeFormat))
	q.Add("until", rr.Until.Format(TimeFormat))
	q.Add("page", strconv.Itoa(rr.Page))

	u.RawQuery = q.Encode()

	return u.String()
}

type Currency struct {
	Currency string  `json:"currency"`
	Amount   float32 `json:"amount"`
}

type Report struct {
	TotalGrand      int            `json:"total_grand"`
	TotalBillable   int            `json:"total_billable"`
	TotalCurrencies []Currency     `json:"total_currencies"`
	Data            []DetailedData `json:"data"`
}

func (r *Report) TotalDuration() int {
	totalDur := 0
	for _, entry := range r.Data {
		totalDur += entry.Dur
	}

	return totalDur
}

type DetailReport struct {
	Report
	TotalCount int `json:"total_count"`
	PerPage    int `json:"per_page"`
}

type DetailedData struct {
	Id          int       `json:"id"`
	Pid         int       `json:"pid"`
	Project     string    `json:"project"`
	Client      string    `json:"client"`
	Tid         int       `json:"tid"`
	Task        string    `json:"task"`
	User        string    `json:"user"`
	Description string    `json:"description"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
	Dur         int       `json:"dur"`
	Updated     time.Time `json:"updated"`
	UseStop     bool      `json:"use_stop"`
	IsBillable  bool      `json:"is_billable"`
	Billable    float32   `json:"billable"`
	Cur         string    `json:"cur"`
	Tags        []string  `json:"tags"`
}

func NewDetailReport(reportRequest ReportRequest) (*DetailReport, error) {
	rr := new(DetailReport)

	request, err := http.NewRequest(http.MethodGet, reportRequest.Url(), nil)
	if err != nil {
		log.Fatal(err)
	}

	request.SetBasicAuth("some-token", "api_token")

	client := &http.Client{}

	res, err := client.Do(request)
	if err != nil {
		return rr, err
	}

	if res.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(res.Body)
		return rr, fmt.Errorf("bad response code. got %d. body: %s", res.StatusCode, body)
	}

	err = json.NewDecoder(res.Body).Decode(rr)
	if err != nil {
		return rr, err
	}

	return rr, nil
}
