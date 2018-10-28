package model

import (
	"database/sql"
)

type User struct {
	ID    int64  `db:"id" json:"id"`
	Name  string `db:"name" json:"name"`
	Mail  string `db:"mail" json:"mail"`
	Phone string `db:"phone" json:"phone"`
}

func Select(db *sql.DB) ([]User, error) {
	rows, err := db.Query("select * from users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		user := User{}
		err = rows.Scan(
			&user.ID,
			&user.Name,
			&user.Mail,
			&user.Phone,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (u *User) Insert(db *sql.DB) (sql.Result, error) {
	result, err := db.Exec("insert into users (name, mail, phone) values (?, ?, ?) ", u.Name, u.Mail, u.Phone)
	if err != nil {
		return nil, err
	}

	return result, nil
}
