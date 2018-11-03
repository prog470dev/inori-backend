package controller

import (
	"database/sql"
	"net/http"
)

type Reservation struct {
	DB *sql.DB
}

func (d *Reservation) GetOwnReservation(w http.ResponseWriter, r *http.Request) {
}

func (d *Reservation) CreateReservation(w http.ResponseWriter, r *http.Request) {
}

func (d *Reservation) CancelReservation(w http.ResponseWriter, r *http.Request) {
}
