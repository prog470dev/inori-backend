package controller

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/prog470dev/inori-backend/model"
	"log"
	"net/http"
	"strconv"
)

type Rider struct {
	DB *sql.DB
}

func (d *Rider) GetRiderDetail(w http.ResponseWriter, r *http.Request) {
	riderID, err := strconv.ParseInt(mux.Vars(r)["rider_id"], 10, 64)
	if err != nil {
		return
	}

	rider, err := model.RiderOne(d.DB, riderID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // 不適切なuser_idと判断(本当はDBのエラーかも)
		return
	}

	_ = JSON(w, http.StatusOK, struct {
		Rider model.Rider `json:"rider"`
	}{
		Rider: *rider,
	})
}

func (d *Rider) UpdateRider(w http.ResponseWriter, r *http.Request) {
	_, err := strconv.ParseInt(mux.Vars(r)["rider_id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var rider model.Rider
	if err := json.NewDecoder(r.Body).Decode(&rider); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = rider.Update(d.DB)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_ = JSON(w, http.StatusOK, struct {
		Rider model.Rider `json:"rider"`
	}{
		Rider: rider,
	})
}

func (d *Rider) SignInRider(w http.ResponseWriter, r *http.Request) {
	type Rb struct {
		Mail string `json:"mail"`
	}

	var rb Rb
	if err := json.NewDecoder(r.Body).Decode(&rb); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rider, err := model.RiderOneWithMail(d.DB, rb.Mail)
	if NotFoundOrErr(w, err) != nil {
		return
	}

	// サインイン成功
	_ = JSON(w, http.StatusOK, struct {
		Rider model.Rider `json:"rider"`
	}{
		Rider: *rider,
	})
}

func (d *Rider) SignUpRider(w http.ResponseWriter, r *http.Request) {
	var rider model.Rider
	if err := json.NewDecoder(r.Body).Decode(&rider); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := rider.Insert(d.DB)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_ = JSON(w, http.StatusOK, struct {
		ID int64 `json:"id"`
	}{
		ID: id,
	})
}
