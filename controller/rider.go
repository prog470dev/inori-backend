package controller

import (
	"database/sql"
	"net/http"
)

type Rider struct {
	DB *sql.DB
}

func (d *Rider) GetRiderDetail(w http.ResponseWriter, r *http.Request) {

}

func (d *Rider) UpdateRider(w http.ResponseWriter, r *http.Request) {

}

func (d *Rider) CreteRider(w http.ResponseWriter, r *http.Request) {

}

func (d *Rider) SignInRider(w http.ResponseWriter, r *http.Request) {

}

func (d *Rider) SignUpRider(w http.ResponseWriter, r *http.Request) {

}
