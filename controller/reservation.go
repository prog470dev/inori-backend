package controller

import (
	"database/sql"
)

type Reservation struct {
	DB *sql.DB
}
