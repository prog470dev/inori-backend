package model

import "database/sql"

type Reservation struct {
	ID      int64 `db:"id" json:"id"`
	OfferID int64 `db:"offer_id" json:"offer_id"`
	RiderID int64 `db:"rider_id" json:"rider_id"`
}

func ReservationsWithRider(db *sql.DB, riderID int64) ([]Reservation, error) {
	rows, err := db.Query("SELECT * FROM reservations WHERE rider_id = ?", riderID)
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
		)
		if err != nil {
			return nil, err
		}
		reservations = append(reservations, reservation)
	}

	return reservations, nil
}

func (r *Reservation) Insert(db *sql.DB) (sql.Result, error) {
	result, err := db.Exec("INSERT INTO reservations (offer_id, rider_id) values (?, ?) ", r.OfferID, r.RiderID)
	if err != nil {
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
	rows, err := db.Query("SELECT * FROM reservations WHERE offer_id = ?", offerID)
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
		)
		if err != nil {
			return nil, err
		}
		reservations = append(reservations, reservation)
	}

	return reservations, nil
}
