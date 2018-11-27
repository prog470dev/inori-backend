package model

import (
	"database/sql"
	"log"
	"time"
)

type Offer struct {
	ID            int64  `db:"id" json:"id"`
	DriverID      int64  `db:"driver_id" json:"driver_id"`
	Start         string `db:"start" json:"start"`
	Goal          string `db:"goal" json:"goal"`
	DepartureTime string `db:"departure_time" json:"departure_time"`
	RiderCapacity int64  `db:"rider_capacity" json:"rider_capacity"`
}

type TimedOffer []Offer // 時系列ソートのための型

func (to TimedOffer) Len() int {
	return len(to)
}

func (to TimedOffer) Less(i, j int) bool {
	ti, err := time.Parse(time.RFC3339, to[i].DepartureTime)
	tj, err := time.Parse(time.RFC3339, to[j].DepartureTime)

	//TODO: パースエラーの処理
	log.Println(err)

	return ti.Before(tj)
}

func (to TimedOffer) Swap(i, j int) {
	to[i], to[j] = to[j], to[i]
}

func OfferOne(db *sql.DB, id int64) (*Offer, error) {
	offer := &Offer{}

	//TODO: 時刻の制約を追加するか検討
	//currentTime := time.Now()

	//if err := db.QueryRow("SELECT * FROM offers WHERE id = ? AND departure_time > ? LIMIT 1", id, currentTime).Scan(
	if err := db.QueryRow("SELECT * FROM offers WHERE id = ? LIMIT 1", id).Scan(
		&offer.ID,
		&offer.DriverID,
		&offer.Start,
		&offer.Goal,
		&offer.DepartureTime,
		&offer.RiderCapacity,
	); err != nil {
		return nil, err
	}

	return offer, nil
}

func OfferOneWithoutTime(db *sql.DB, id int64) (*Offer, error) {
	offer := &Offer{}

	if err := db.QueryRow("SELECT * FROM offers WHERE id = ? LIMIT 1", id).Scan(
		&offer.ID,
		&offer.DriverID,
		&offer.Start,
		&offer.Goal,
		&offer.DepartureTime,
		&offer.RiderCapacity,
	); err != nil {
		return nil, err
	}

	return offer, nil
}

func OffersAll(db *sql.DB) ([]Offer, error) {
	currentTime := time.Now()

	// 締切を１時間前まで
	currentTime = currentTime.Add(time.Duration(1) * time.Hour)

	rows, err := db.Query("SELECT * FROM offers WHERE departure_time > ?", currentTime)
	if err != nil {
		return nil, err
	}

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

	// 確認できるのは１時間後まで
	currentTime = currentTime.Add(time.Duration(-1) * time.Hour)

	rows, err := db.Query("SELECT * FROM offers WHERE driver_id = ? AND departure_time > ?", driverID, currentTime)
	if err != nil {
		return nil, err
	}

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
	result, err := db.Exec("INSERT INTO offers (driver_id, start, goal, departure_time, rider_capacity) values "+
		"(?, ?, ?, ?, ?) ",
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
