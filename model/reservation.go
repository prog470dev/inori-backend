package model

import (
	"database/sql"
	"log"
	"time"
)

type Reservation struct {
	ID            int64  `db:"id" json:"id"`
	OfferID       int64  `db:"offer_id" json:"offer_id"`
	RiderID       int64  `db:"rider_id" json:"rider_id"`
	DepartureTime string `db:"departure_time" json:"departure_time"`
}

func ReservationOne(db *sql.DB, id int64) (*Reservation, error) {
	reservation := &Reservation{}

	currentTime := time.Now()

	if err := db.QueryRow("SELECT * FROM reservations WHERE id = ? AND departure_time > ? LIMIT 1", id, currentTime).Scan(
		&reservation.ID,
		&reservation.OfferID,
		&reservation.RiderID,
		&reservation.DepartureTime,
	); err != nil {
		return nil, err
	}

	return reservation, nil
}

func ReservationsWithRider(db *sql.DB, riderID int64) ([]Reservation, error) {
	currentTime := time.Now()

	// 支払い受付は12時間あとまで
	currentTime = currentTime.Add(time.Duration(-12) * time.Hour)

	log.Println("currentTime: ", currentTime)

	rows, err := db.Query("SELECT * FROM reservations WHERE rider_id = ? AND departure_time > ? ", riderID, currentTime)
	if err != nil {
		return nil, err
	}

	reservations := []Reservation{}
	for rows.Next() {
		reservation := Reservation{}
		err = rows.Scan(
			&reservation.ID,
			&reservation.OfferID,
			&reservation.RiderID,
			&reservation.DepartureTime,
		)
		if err != nil {
			return nil, err
		}
		reservations = append(reservations, reservation)
		log.Println("reservation.DepartureTime> ", reservation.DepartureTime)
	}

	log.Println("len(reservations): ", len(reservations))

	return reservations, nil
}

func (r *Reservation) Insert(db *sql.DB) (sql.Result, error) {
	result, err := db.Exec("INSERT INTO reservations (offer_id, rider_id, departure_time) values "+
		"(?, ?, ?) ",
		r.OfferID, r.RiderID, r.DepartureTime)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return result, nil
}

func (r *Reservation) Delete(db *sql.DB) (sql.Result, error) {
	result, err := db.Exec("DELETE FROM reservations WHERE id = ?", r.ID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func DeleteOfferReservation(db *sql.DB, offer_id int) (sql.Result, error) {
	result, err := db.Exec("DELETE FROM reservations WHERE offer_id = ?", offer_id)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// オファーの予約
func ReservationsWithOffer(db *sql.DB, offerID int64) ([]*Reservation, error) {
	currentTime := time.Now()

	rows, err := db.Query("SELECT * FROM reservations WHERE offer_id = ? AND departure_time > ? ", offerID, currentTime)
	if err == sql.ErrNoRows {
		return []*Reservation{}, nil
	}
	if err != nil {
		return nil, err
	}

	reservations := []*Reservation{}
	for rows.Next() {
		reservation := &Reservation{}
		err = rows.Scan(
			&reservation.ID,
			&reservation.OfferID,
			&reservation.RiderID,
			&reservation.DepartureTime,
		)
		if err != nil {
			return nil, err
		}
		reservations = append(reservations, reservation)
	}

	return reservations, nil
}
