package model

import "time"

// see https://ashitani.jp/golangtips/tips_time.html#time_Format
// see https://pliutau.com/working-with-db-time-in-go/
// https://twinbird-htn.hatenablog.com/entry/2017/03/31/140335
const TimeLayout = "2006-01-02 15:04:05"

type Offer struct {
	ID            int64     `db:"id" json:"id"`
	DriverID      int64     `db:"driver_id" json:"driver_id"`
	Start         string    `db:"start" json:"start"`
	Goal          string    `db:"goal" json:"goal"`
	DepartureTime time.Time `db:"departure_time" json:"departure_time"` //use TimeLayout
	RiderCapacity int64     `db:"rider_capacity" json:"rider_capacity"`
}
