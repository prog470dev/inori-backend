package model

import (
	"database/sql"
)

type Rider struct {
	ID        int64  `db:"id" json:"id"`
	FirstName string `db:"first_name" json:"first_name"`
	LastName  string `db:"last_name" json:"last_name"`
	Grade     string `db:"grade" json:"grade"`
	Major     string `db:"major" json:"major"`
	Mail      string `db:"mail" json:"mail"`
	Phone     string `db:"phone" json:"phone"`
	ImageUrl  string `db:"image_url" json:"image_url"`
}

func RiderOne(db *sql.DB, id int64) (*Rider, error) {
	rider := &Rider{}

	if err := db.QueryRow("SELECT * FROM riders WHERE id = ? LIMIT 1", id).Scan(
		&rider.ID,
		&rider.FirstName,
		&rider.LastName,
		&rider.Grade,
		&rider.Major,
		&rider.Mail,
		&rider.Phone,
		&rider.ImageUrl,
	); err != nil {
		return nil, err
	}

	return rider, nil
}

func RiderOneWithMail(db *sql.DB, mail string) (*Rider, error) {
	rider := &Rider{}

	if err := db.QueryRow("SELECT * FROM riders WHERE mail = ? LIMIT 1", mail).Scan(
		&rider.ID,
		&rider.FirstName,
		&rider.LastName,
		&rider.Grade,
		&rider.Major,
		&rider.Mail,
		&rider.Phone,
		&rider.ImageUrl,
	); err != nil {
		return nil, err
	}

	return rider, nil
}

func (d *Rider) Update(db *sql.DB) (sql.Result, error) {
	result, err := db.Exec("UPDATE riders SET first_name=?, last_name=?, grade=?, major=?, mail=?, phone=?, image_url=? WHERE id = ?",
		d.FirstName, d.LastName, d.Grade, d.Major, d.Mail, d.Phone, d.ImageUrl, d.ID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (d *Rider) Insert(db *sql.DB) (sql.Result, error) {
	result, err := db.Exec("INSERT INTO riders (first_name, last_name, grade, major, mail, phone, image_url) values"+
		" (?, ?, ?, ?, ?, ?, ?) ",
		d.FirstName, d.LastName, d.Grade, d.Major, d.Mail, d.Phone, d.ImageUrl)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (d *Rider) Delete(db *sql.DB) (sql.Result, error) {
	result, err := db.Exec("DELETE FROM riders WHERE id = ?", d.ID)
	if err != nil {
		return nil, err
	}

	return result, nil
}
