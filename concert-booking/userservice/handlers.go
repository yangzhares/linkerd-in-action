package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/yangzhares/linkerd-in-action/concert-booking/db"
	"github.com/yangzhares/linkerd-in-action/concert-booking/util"
)

type UserHandler struct {
	DB                 *db.DB
	BookingServiceAddr string
	ConcertServiceAddr string
}

type booking struct {
	UserID    string `json:"user_id"`
	Date      string `json:"date"`
	ConcertID string `json:"concert_id"`
}

type bookings []booking

type concert struct {
	ConcertName string `json:"concert_name"`
	Singer      string `json:"singer"`
	Location    string `json:"location"`
	//	Street      string `json:"street"` //v2
}

type result struct {
	Date string `json:"date"`
	*concert
}

func NewUserHandler(db *db.DB, bookingServiceAddr, concertServiceAddr string) *UserHandler {
	return &UserHandler{
		DB:                 db,
		BookingServiceAddr: bookingServiceAddr,
		ConcertServiceAddr: concertServiceAddr,
	}
}

func (handler *UserHandler) AddUser(w http.ResponseWriter, r *http.Request) {
	user := User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		util.ResponseWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request payload: %v", err))
		return
	}

	defer r.Body.Close()

	rows, _ := handler.DB.Raw("SELECT id FROM users WHERE id = ?", user.ID).Rows()
	defer rows.Close()

	if rows.Next() {
		util.ResponseWithError(w, http.StatusFound, "User already existed")
		return
	}

	dbconn := handler.DB.Exec("INSERT INTO users(id, name, age) VALUES(?, ?, ?)", user.ID, user.Name, user.Age)
	if dbconn.Error != nil {
		util.ResponseWithError(w, http.StatusInternalServerError, fmt.Sprintf("Create user failure: %v", dbconn.Error))
		return
	}

	util.ResponseWithJson(w, http.StatusOK, user)
}

func (handler *UserHandler) FindUsers(w http.ResponseWriter, r *http.Request) {
	var users []User
	rows, err := handler.DB.Raw("SELECT id, name, age FROM users").Rows()
	if err != nil {
		util.ResponseWithError(w, http.StatusNotFound, fmt.Sprintf("User not found: %v", err))
		return
	}

	defer rows.Close()

	for rows.Next() {
		var id string
		var name string
		var age int

		rows.Scan(&id, &name, &age)
		user := User{
			ID:   id,
			Name: name,
			Age:  age,
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		util.ResponseWithError(w, http.StatusNotFound, fmt.Sprintf("User not found: %v", err))
		return
	}
	util.ResponseWithJson(w, http.StatusOK, users)
}

func (handler *UserHandler) FindUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var user User
	rows, err := handler.DB.Raw("SELECT name, age FROM users WHERE id = ?", id).Rows()
	if err != nil {
		util.ResponseWithError(w, http.StatusNotFound, fmt.Sprintf("User not found: %v", err))
		return
	}

	defer rows.Close()

	if !rows.Next() {
		util.ResponseWithError(w, http.StatusNotFound, "User not found")
		return
	}
	{
		var name string
		var age int

		rows.Scan(&name, &age)
		user = User{
			ID:   id,
			Name: name,
			Age:  age,
		}
	}

	util.ResponseWithJson(w, http.StatusOK, user)
}

func (handler *UserHandler) FindUserBookingsByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	rows, err := handler.DB.Raw("SELECT id FROM users WHERE id = ?", userID).Rows()
	if err != nil {
		util.ResponseWithError(w, http.StatusNotFound, fmt.Sprintf("User not found: %v", err))
		return
	}

	defer rows.Close()

	if !rows.Next() {
		util.ResponseWithError(w, http.StatusOK, "User not found")
		return
	}

	bookings, err := quertBookings(handler.BookingServiceAddr, userID)
	if err != nil {
		util.ResponseWithError(w, http.StatusNotFound, fmt.Sprintf("Can't query user's bookings from booking service: %v", err))
		return
	}

	query := make(map[string][]result)
	query[userID] = []result{}
	for _, booking := range bookings {
		concertID := booking.ConcertID
		concert, err := queryConcert(handler.ConcertServiceAddr, concertID)
		if err != nil {
			util.ResponseWithError(w, http.StatusNotFound, fmt.Sprintf("Can't query concert from concert service: %v", err))
			return
		}

		query[userID] = append(query[userID], result{
			Date:    booking.Date,
			concert: concert,
		})
	}

	util.ResponseWithJson(w, http.StatusOK, query)
}

func quertBookings(endpoint, userID string) (bookings, error) {
	endpoint = standarizeEndpoint(endpoint)

	resp, err := http.Get(fmt.Sprintf("%s/bookings/%s", endpoint, userID))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("Bookings not found")
	}

	var bks bookings
	if err := json.NewDecoder(resp.Body).Decode(&bks); err != nil {
		return nil, err
	}

	return bks, nil
}

func queryConcert(endpoint, concertID string) (*concert, error) {
	endpoint = standarizeEndpoint(endpoint)

	resp, err := http.Get(fmt.Sprintf("%s/concerts/%s", endpoint, concertID))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New("Concert not found")
	}

	var concert concert
	if err := json.NewDecoder(resp.Body).Decode(&concert); err != nil {
		return nil, err
	}

	return &concert, nil
}

func standarizeEndpoint(endpoint string) string {
	if !strings.HasPrefix(endpoint, "http://") {
		endpoint = "http://" + endpoint
	}
	return endpoint
}
