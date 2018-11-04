package model

import (
	"database/sql"
	"time"
)

type Offer struct {
	ID            int64     `db:"id" json:"id"`
	DriverID      int64     `db:"driver_id" json:"driver_id"`
	Start         string    `db:"start" json:"start"`
	Goal          string    `db:"goal" json:"goal"`
	DepartureTime time.Time `db:"departure_time" json:"departure_time"` //use TimeLayout
	RiderCapacity int64     `db:"rider_capacity" json:"rider_capacity"`
}

func OfferOne(db *sql.DB, id int64) (*Offer, error) {
	offer := &Offer{}

	currentTime := time.Now()

	if err := db.QueryRow("SELECT * FROM offers WHERE id = ? AND departure_time > ? LIMIT 1", id, currentTime).Scan(
		&offer.ID,
		&offer.DriverID,
		&offer.Start,
		&offer.Goal,
		&offer.DepartureTime,
		&offer.RiderCapacity,
	); err != nil {
		return nil, err
	}
	defer db.Close()

	return offer, nil
}

func OffersAll(db *sql.DB) ([]Offer, error) {
	currentTime := time.Now()

	rows, err := db.Query("SELECT * FROM offers WHERE departure_time > ?", currentTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	offers := []Offer{}
	for rows.Next() {
		offer := Offer{}
		err = rows.Scan(
			&offer.ID,
			&offer.DriverID,
			&offer.Start,
			&offer.Goal,
			&offer.DepartureTime,
			&offer.RiderCapacity,
		)
		if err != nil {
			return nil, err
		}
		offers = append(offers, offer)
	}

	return offers, nil
}

func OffersWithDriver(db *sql.DB, driverID int64) ([]Offer, error) {
	currentTime := time.Now()

	rows, err := db.Query("SELECT * FROM offers WHERE driver_id = ? AND departure_time > ?", driverID, currentTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	offers := []Offer{}
	for rows.Next() {
		offer := Offer{}
		err = rows.Scan(
			&offer.ID,
			&offer.DriverID,
			&offer.Start,
			&offer.Goal,
			&offer.DepartureTime,
			&offer.RiderCapacity,
		)
		if err != nil {
			return nil, err
		}
		offers = append(offers, offer)
	}

	return offers, nil
}

func (o *Offer) Insert(db *sql.DB) (sql.Result, error) {
	result, err := db.Exec("INSERT INTO offers (driver_id, start, goal, departure_time, rider_capacity) values"+
		" (?, ?, ?, ?, ?) ",
		o.DriverID, o.Start, o.Goal, o.DepartureTime, o.RiderCapacity)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (o *Offer) Delete(db *sql.DB) (sql.Result, error) {
	result, err := db.Exec("DELETE FROM offers WHERE id = ?", o.ID)

	if err != nil {
		return nil, err
	}

	return result, nil
}
