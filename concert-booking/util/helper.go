package util

import (
	"encoding/json"
	"net/http"
)

var Format = "2006-01-02 15:04:05"

func ResponseWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func ResponseWithError(w http.ResponseWriter, code int, message string) {
	ResponseWithJson(w, code, map[string]string{"error": message})
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	ResponseWithJson(w, http.StatusOK, "OKOKOK")
}
