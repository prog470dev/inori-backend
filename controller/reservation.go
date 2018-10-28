package controller

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type Reserve struct {
	DB *sql.DB
}

type Reservation struct {
	ID       int64   `json:"id"`
	DriverID int64   `json:"driver"`
	Riders   []int64 `json:"riders"`
}

func MockHandler(w http.ResponseWriter, r *http.Request) {

	var reservation Reservation
	if err := json.NewDecoder(r.Body).Decode(&reservation); err != nil {
		http.Error(w, "http.StatusBadRequest", http.StatusBadRequest)
		return
	}

	JSON(w, http.StatusOK, struct {
		RiderCount int `json:"rider_count"`
	}{
		RiderCount: len(reservation.Riders),
	})
}

func JSON(w http.ResponseWriter, code int, data interface{}) error {
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(data)
}
