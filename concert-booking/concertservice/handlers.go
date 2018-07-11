package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	//uuid "github.com/satori/go.uuid"

	"github.com/yangzhares/linkerd-in-action/concert-booking/db"
	"github.com/yangzhares/linkerd-in-action/concert-booking/util"
)

type ConcertHandler struct {
	DB *db.DB
}

//v2
var uuidGenerator = "http://httpbin.org/uuid"

func NewConcertHandler(db *db.DB) *ConcertHandler {
	return &ConcertHandler{
		DB: db,
	}
}

func (h *ConcertHandler) AddConcert(w http.ResponseWriter, r *http.Request) {
	var concert Concert

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		util.ResponseWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request payload: %v", err))
		return
	}

	if err := concert.UnmarshalJSON(body); err != nil {
		util.ResponseWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request payload: %v", err))
		return
	}
	defer r.Body.Close()

	if concert.ID != "" {
		rows, _ := h.DB.Raw("SELECT id FROM concerts WHERE id = ?", concert.ID).Rows()
		defer rows.Close()

		if rows.Next() {
			util.ResponseWithError(w, http.StatusFound, "Concert already existed")
			return
		}
	}

	if concert.ID == "" {
		rows, _ := h.DB.Raw("SELECT id FROM concerts WHERE concert_name = ? and singer = ?", concert.ConcertName, concert.Singer).Rows()
		defer rows.Close()

		if rows.Next() {
			util.ResponseWithError(w, http.StatusFound, "Concert already existed")
			return
		}

		// v1
		//concert.ID = uuid.NewV4().String()

		//v2

		resp, err := http.Get(uuidGenerator)
		if err != nil {
			util.ResponseWithError(w, http.StatusInternalServerError, fmt.Sprintf("UUID Generator error: %v", err))
		}

		defer resp.Body.Close()
		type stubID struct {
			UUID string `json:"uuid"`
		}
		var id stubID

		if err := json.NewDecoder(resp.Body).Decode(&id); err == nil {
			concert.ID = id.UUID
		}

	}

	dbconn := h.DB.Exec("INSERT INTO concerts(id, concert_name, singer, start_date, end_date, location, street) VALUES(?, ?, ?, ?, ?, ?, ?)", concert.ID, concert.ConcertName, concert.Singer, concert.StartDate, concert.EndDate, concert.Location, concert.Street)
	if dbconn.Error != nil {
		util.ResponseWithError(w, http.StatusInternalServerError, fmt.Sprintf("Add concert failure: %v", err))
		return
	}
	util.ResponseWithJson(w, http.StatusOK, concert)
}

func (h *ConcertHandler) FindConcertByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var concert Concert
	rows, err := h.DB.Raw("SELECT * FROM concerts WHERE id = ?", id).Rows()
	if err != nil {
		util.ResponseWithError(w, http.StatusNotFound, fmt.Sprintf("Concert not found: %v", err))
		return
	}

	defer rows.Close()

	if !rows.Next() {
		util.ResponseWithError(w, http.StatusNotFound, "Concert not found")
		return
	}
	{
		var id string
		var concertName string
		var singer string
		var startDate time.Time
		var endDate time.Time
		var location string
		var street string

		rows.Scan(&id, &concertName, &singer, &startDate, &endDate, &location, &street)
		concert = Concert{
			ID:          id,
			ConcertName: concertName,
			Singer:      singer,
			StartDate:   startDate,
			EndDate:     endDate,
			Location:    location,
			Street:      street,
		}
	}

	util.ResponseWithJson(w, http.StatusOK, concert)
}

func (h *ConcertHandler) FindConcerts(w http.ResponseWriter, r *http.Request) {
	var concerts []Concert

	rows, err := h.DB.Raw("SELECT * FROM concerts").Rows()
	if err != nil {
		util.ResponseWithError(w, http.StatusNotFound, fmt.Sprintf("Concert not found: %v", err))
		return
	}

	defer rows.Close()

	for rows.Next() {
		var id string
		var concertName string
		var singer string
		var startDate time.Time
		var endDate time.Time
		var location string
		var street string

		rows.Scan(&id, &concertName, &singer, &startDate, &endDate, &location, &street)
		concert := Concert{
			ID:          id,
			ConcertName: concertName,
			Singer:      singer,
			StartDate:   startDate,
			EndDate:     endDate,
			Location:    location,
			Street:      street,
		}

		concerts = append(concerts, concert)
	}

	if err := rows.Err(); err != nil {
		util.ResponseWithError(w, http.StatusNotFound, fmt.Sprintf("Concert not found: %v", err))
		return
	}
	util.ResponseWithJson(w, http.StatusOK, concerts)
}
