package webapi_no

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

const holidaysUrl string = "https://webapi.no/api/v1/holidays/"
const webApiTimeFormat = "2006-01-02"

type Time struct {
	time.Time
}

func (t *Time) UnmarshalJSON(b []byte) error {
	newTime, err := time.Parse(webApiTimeFormat, string(b[1:11]))
	if err != nil {
		return err
	}

	*t = Time{newTime}
	return nil
}

type Holidays struct {
	Data []Holiday `json:"data"`
}

type Holiday struct {
	Date        Time   `json:"date"`
	Description string `json:"description"`
}

func NewHolidays(year int) (*Holidays, error) {
	h := new(Holidays)

	res, err := http.Get(holidaysUrl + strconv.Itoa(year))
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(res.Body)
		return h, fmt.Errorf("bad response code. got %d. body: %s", res.StatusCode, body)
	}

	err = json.NewDecoder(res.Body).Decode(h)
	if err != nil {
		return h, err
	}

	return h, nil
}

func NewHolidaysMulti(years []int) (*Holidays, error) {
	h := new(Holidays)

	for _, year := range years {
		holidays, err := NewHolidays(year)
		if err != nil {
			return h, err
		}

		h.Data = append(h.Data, holidays.Data...)
	}

	return h, nil
}
