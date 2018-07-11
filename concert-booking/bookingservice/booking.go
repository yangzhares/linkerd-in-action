package main

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/yangzhares/linkerd-in-action/concert-booking/util"
)

type Booking struct {
	UserID    string    `json:"user_id"`
	Date      time.Time `json:"date"`
	ConcertID string    `json:"concert_id"`
}

func (b *Booking) UnmarshalJSON(buf []byte) error {
	var raws map[string]string

	err := json.Unmarshal(buf, &raws)
	if err != nil {
		return err
	}

	for k, v := range raws {
		if strings.ToLower(k) == "user_id" {
			b.UserID = v
		}
		if strings.ToLower(k) == "concert_id" {
			b.ConcertID = v
		}

		if strings.ToLower(k) == "date" {
			t, err := time.ParseInLocation(util.Format, v, time.Local)
			if err != nil {
				return err
			}

			b.Date = t
		}
	}
	return nil
}

func (b *Booking) MarshalJSON() ([]byte, error) {
	stub := struct {
		UserID    string `json:"user_id"`
		Date      string `json:"date"`
		ConcertID string `json:"concert_id"`
	}{
		UserID:    b.UserID,
		Date:      b.Date.Format(util.Format),
		ConcertID: b.ConcertID,
	}

	return json.Marshal(stub)
}
