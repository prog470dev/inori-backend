package model

import (
	"database/sql"
)

type Driver struct {
	ID        int64  `db:"id" json:"id"`
	FirstName string `db:"first_name" json:"first_name"`
	LastName  string `db:"last_name" json:"last_name"`
	Grade     string `db:"grade" json:"grade"`
	Major     string `db:"major" json:"major"`
	Mail      string `db:"mail" json:"mail"`
	Phone     string `db:"phone" json:"phone"`
	CarColor  string `db:"car_color" json:"car_color"`
	CarNumber string `db:"car_number" json:"car_number"`
}

func DriverOne(db *sql.DB, id int64) (*Driver, error) {
	driver := &Driver{}

	if err := db.QueryRow("SELECT * FROM drivers WHERE id = ? LIMIT 1", id).Scan(
		&driver.ID,
		&driver.FirstName,
		&driver.LastName,
		&driver.Grade,
		&driver.Major,
		&driver.Mail,
		&driver.Phone,
		&driver.CarColor,
		&driver.CarNumber,
	); err != nil {
		return nil, err
	}

	return driver, nil
}

func DriverOneWithMail(db *sql.DB, mail string) (*Driver, error) {
	driver := &Driver{}

	if err := db.QueryRow("SELECT * FROM drivers WHERE mail = ? LIMIT 1", mail).Scan(
		&driver.ID,
		&driver.FirstName,
		&driver.LastName,
		&driver.Grade,
		&driver.Major,
		&driver.Mail,
		&driver.Phone,
		&driver.CarColor,
		&driver.CarNumber,
	); err != nil {
		return nil, err
	}

	return driver, nil
}

func (d *Driver) Update(db *sql.DB) (sql.Result, error) {
	result, err := db.Exec("UPDATE drivers SET first_name=?, last_name=?, grade=?, major=?, mail=?, phone=?, car_color=?, car_number=? WHERE id = ?",
		d.FirstName, d.LastName, d.Grade, d.Major, d.Mail, d.Phone, d.CarColor, d.CarNumber, d.ID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (d *Driver) Insert(db *sql.DB) (sql.Result, error) {
	result, err := db.Exec("INSERT INTO drivers (first_name, last_name, grade, major, mail, phone, car_color, car_number) values"+
		" (?, ?, ?, ?, ?, ?, ?, ?) ",
		d.FirstName, d.LastName, d.Grade, d.Major, d.Mail, d.Phone, d.CarColor, d.CarNumber)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (d *Driver) Delete(db *sql.DB) (sql.Result, error) {
	result, err := db.Exec("DELETE FROM drivers WHERE id = ?", d.ID)
	if err != nil {
		return nil, err
	}

	return result, nil
}
