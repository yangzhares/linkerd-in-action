package main

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/yangzhares/linkerd-in-action/concert-booking/util"
)

type Concert struct {
	ID          string    `json:"id"`
	ConcertName string    `json:"concert_name"`
	Singer      string    `json:"singer"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Location    string    `json:"location"`
	Street      string    `json:"street"`
}

func (c *Concert) UnmarshalJSON(buf []byte) error {
	var raws map[string]string

	err := json.Unmarshal(buf, &raws)
	if err != nil {
		return err
	}

	for k, v := range raws {
		if strings.ToLower(k) == "id" {
			c.ID = v
		}
		if strings.ToLower(k) == "concert_name" {
			c.ConcertName = v
		}

		if strings.ToLower(k) == "singer" {
			c.Singer = v
		}

		if strings.ToLower(k) == "location" {
			c.Location = v
		}

		if strings.ToLower(k) == "street" {
			c.Street = v
		}

		if strings.ToLower(k) == "start_date" {
			t, err := time.ParseInLocation(util.Format, v, time.Local)
			if err != nil {
				return err
			}

			c.StartDate = t
		}

		if strings.ToLower(k) == "end_date" {
			t, err := time.ParseInLocation(util.Format, v, time.Local)
			if err != nil {
				return err
			}

			c.EndDate = t
		}
	}
	return nil
}

func (c *Concert) MarshalJSON() ([]byte, error) {
	stub := struct {
		ID          string `json:"id"`
		ConcertName string `json:"concert_name"`
		Singer      string `json:"singer"`
		StartDate   string `json:"start_date"`
		EndDate     string `json:"end_date"`
		Location    string `json:"location"`
		Street      string `json:"street"`
	}{
		ID:          c.ID,
		ConcertName: c.ConcertName,
		Singer:      c.Singer,
		StartDate:   c.StartDate.Format(util.Format),
		EndDate:     c.EndDate.Format(util.Format),
		Location:    c.Location,
		Street:      c.Street,
	}

	return json.Marshal(stub)
}
