package model

type Reservation struct {
	ID      int64 `db:"id" json:"id"`
	OfferID int64 `db:"offer_id" json:"offer_id"`
	RiderID int64 `db:"rider_id" json:"rider_id"`
}
