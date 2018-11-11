package controller

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/prog470dev/inori-backend/model"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type Reservation struct {
	DB *sql.DB
}

func (d *Reservation) GetRiderOffers(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.URL.RequestURI())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	query := u.Query()
	if len(query["rider_id"]) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	riderID, err := strconv.ParseInt(query["rider_id"][0], 10, 64)

	reservations, err := model.ReservationsWithRider(d.DB, riderID)
	if NotFoundOrErr(w, err) != nil {
		return
	}

	JSON(w, http.StatusOK, struct {
		Reservations []model.Reservation `json:"reservations"`
	}{
		Reservations: reservations,
	})
}

func (d *Reservation) CreateReservation(w http.ResponseWriter, r *http.Request) {
	var reservation model.Reservation
	if err := json.NewDecoder(r.Body).Decode(&reservation); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	reservations, err := model.ReservationsWithOffer(d.DB, reservation.OfferID)
	if err != nil && err != sql.ErrNoRows { // 予約がないのはOK
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 存在しないOfferはエラー
	offer, err := model.OfferOne(d.DB, reservation.OfferID)
	if NotFoundOrErr(w, err) != nil {
		return
	}

	// 満員（クライアント側の同期がリアルタイムやられていれば基本発生しない）
	if len(reservations) == int(offer.RiderCapacity) {
		//TODO: 満員であることを伝えるエラー
		log.Println(len(reservations), int(offer.RiderCapacity))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := reservation.Insert(d.DB)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	JSON(w, http.StatusOK, struct {
		ID int64 `json:"id"`
	}{
		ID: id,
	})
}

func (d *Reservation) CancelReservation(w http.ResponseWriter, r *http.Request) {
	reservationID, err := strconv.ParseInt(mux.Vars(r)["reservation_id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	reservation := &model.Reservation{}
	// delete文は対象レコードの有無に関する情報は返さないため、事前にselectを実行
	reservation, err = model.ReservationOne(d.DB, reservationID)
	if NotFoundOrErr(w, err) != nil {
		return
	}

	_, err = reservation.Delete(d.DB)
	if NotFoundOrErr(w, err) != nil {
		return
	}

	JSON(w, http.StatusOK, struct {
		ID int64 `json:"id"`
	}{
		ID: reservationID,
	})
}
