package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"strings"

	"github.com/gorilla/mux"
	"github.com/yangzhares/linkerd-in-action/concert-booking/db"
	"github.com/yangzhares/linkerd-in-action/concert-booking/util"
)

type BookingHandler struct {
	DB                 *db.DB
	ConcertServiceAddr string
}

func NewBookingHandler(db *db.DB, concertServiceAddr string) *BookingHandler {
	return &BookingHandler{
		DB:                 db,
		ConcertServiceAddr: concertServiceAddr,
	}
}

func (h *BookingHandler) AddBooking(w http.ResponseWriter, r *http.Request) {
	var booking Booking

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		util.ResponseWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request payload: %v", err))
		return
	}

	if err := booking.UnmarshalJSON(body); err != nil {
		util.ResponseWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request payload: %v", err))
		return
	}
	defer r.Body.Close()

	rows, err := h.DB.Raw("SELECT date FROM bookings WHERE user_id = ? and concert_id = ?", booking.UserID, booking.ConcertID).Rows()
	if err != nil {
		util.ResponseWithError(w, http.StatusNotFound, fmt.Sprintf("Booking not found: %v", err))
		return
	}

	defer rows.Close()

	if rows.Next() {
		var date time.Time
		rows.Scan(&date)
		duration := date.Sub(booking.Date)
		if int64(duration) == 0 {
			util.ResponseWithError(w, http.StatusFound, "User already booked concert")
			return
		} else if int64(duration) > 0 {
			dbconn := h.DB.Exec("INSERT INTO bookings(user_id, date, concert_id) VALUES(?, ?, ?)", booking.UserID, booking.Date, booking.ConcertID)
			if dbconn.Error != nil {
				util.ResponseWithError(w, http.StatusInternalServerError, fmt.Sprintf("Book concert failure: %v", dbconn.Error))
				return
			}
			util.ResponseWithJson(w, http.StatusOK, booking)
			return
		}
	}

	if !strings.HasPrefix(h.ConcertServiceAddr, "http://") {
		h.ConcertServiceAddr = "http://" + h.ConcertServiceAddr
	}
	resp, err := http.Get(fmt.Sprintf("%s/concerts/%s", h.ConcertServiceAddr, booking.ConcertID))
	if err != nil {
		util.ResponseWithError(w, http.StatusNotFound, fmt.Sprintf("Can't query concert from concert service: %v", err))
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		util.ResponseWithError(w, http.StatusNotFound, fmt.Sprint("Can't query concert from concert service: Concert not found"))
		return
	}

	dbconn := h.DB.Exec("INSERT INTO bookings(user_id, date, concert_id) VALUES(?, ?, ?)", booking.UserID, booking.Date, booking.ConcertID)
	if dbconn.Error != nil {
		util.ResponseWithError(w, http.StatusInternalServerError, fmt.Sprintf("Book concert failure: %v", dbconn.Error))
		return
	}
	util.ResponseWithJson(w, http.StatusOK, booking)
}

func (h *BookingHandler) FindBookings(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	var bookings []Booking
	rows, err := h.DB.Raw("SELECT * FROM bookings WHERE user_id = ?", userID).Rows()
	if err != nil {
		util.ResponseWithError(w, http.StatusNotFound, fmt.Sprintf("Booking not found: %v", err))
		return
	}

	defer rows.Close()

	if !rows.Next() {
		util.ResponseWithError(w, http.StatusNotFound, "Booking not found")
		return
	}
	{
		var userID string
		var date time.Time
		var concertID string

		rows.Scan(&userID, &date, &concertID)
		booking := Booking{
			UserID:    userID,
			Date:      date,
			ConcertID: concertID,
		}
		bookings = append(bookings, booking)
	}

	for rows.Next() {
		var userID string
		var date time.Time
		var concertID string

		rows.Scan(&userID, &date, &concertID)
		booking := Booking{
			UserID:    userID,
			Date:      date,
			ConcertID: concertID,
		}
		bookings = append(bookings, booking)
	}

	util.ResponseWithJson(w, http.StatusOK, bookings)
}
