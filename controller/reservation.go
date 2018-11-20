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

func (re *Reservation) GetRiderOffers(w http.ResponseWriter, r *http.Request) {
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

	reservations, err := model.ReservationsWithRider(re.DB, riderID)
	if NotFoundOrErr(w, err) != nil {
		return
	}

	// 時間文字列の変換
	for _, reserve := range reservations {
		t, err := SwitchTimeStrStyle(reserve.DepartureTime)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		reserve.DepartureTime = t
	}

	JSON(w, http.StatusOK, struct {
		Reservations []model.Reservation `json:"reservations"`
	}{
		Reservations: reservations,
	})
}

func (re *Reservation) CreateReservation(w http.ResponseWriter, r *http.Request) {
	var reservation model.Reservation
	if err := json.NewDecoder(r.Body).Decode(&reservation); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	reservations, err := model.ReservationsWithOffer(re.DB, reservation.OfferID)
	if err != nil && err != sql.ErrNoRows { // 予約がないのはOK
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 存在しないOfferはエラー
	offer, err := model.OfferOne(re.DB, reservation.OfferID)
	if NotFoundOrErr(w, err) != nil {
		log.Println(err)
		return
	}

	// 満員（クライアント側の同期がリアルタイムやられていれば基本発生しない）
	if len(reservations) == int(offer.RiderCapacity) {
		log.Println(len(reservations), int(offer.RiderCapacity))
		w.WriteHeader(http.StatusBadRequest) //TODO: 満員であることを伝えるボディを返す

		JSON(w, http.StatusOK, struct {
			Message string `json:"message"`
		}{
			Message: "no capacity",
		})
		return
	}

	log.Println(reservation.DepartureTime)

	result, err := reservation.Insert(re.DB)
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

	//プッシュ通知 (ドライバ向け)
	token, err := model.TokenOneDriver(re.DB, offer.DriverID)
	if err != nil {
		NotFoundOrErr(w, err)
		return
	}
	pushData := &PushData{
		To:          token.PushToken,
		Type:        "reserved_offer",
		OfferID:     reservation.OfferID,
		MessageBody: "あなたの相乗りオファーが予約されました。",
		Title:       "あなたの相乗りオファーが予約されました。",
	}
	err = SendPushMessage(pushData)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	JSON(w, http.StatusOK, struct {
		ID int64 `json:"id"`
	}{
		ID: id,
	})
}

func (re *Reservation) CancelReservation(w http.ResponseWriter, r *http.Request) {
	reservationID, err := strconv.ParseInt(mux.Vars(r)["reservation_id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	reservation := &model.Reservation{}
	// delete文は対象レコードの有無に関する情報は返さないため、事前にselectを実行
	reservation, err = model.ReservationOne(re.DB, reservationID)
	if NotFoundOrErr(w, err) != nil {
		return
	}

	_, err = reservation.Delete(re.DB)
	if NotFoundOrErr(w, err) != nil {
		return
	}

	//プッシュ通知 (ドライバ向け)
	offer, err := model.OfferOne(re.DB, reservation.OfferID)
	if err != nil {
		NotFoundOrErr(w, err)
		return
	}
	token, err := model.TokenOneDriver(re.DB, offer.DriverID)
	if err != nil {
		NotFoundOrErr(w, err)
		return
	}
	pushData := &PushData{
		To:          token.PushToken,
		Type:        "canceled_reservation",
		OfferID:     reservation.OfferID,
		MessageBody: "あなたの相乗りオファーへの予約がキャンセルされました。",
		Title:       "あなたの相乗りオファーへの予約がキャンセルされました。",
	}
	err = SendPushMessage(pushData)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	JSON(w, http.StatusOK, struct {
		ID int64 `json:"id"`
	}{
		ID: reservationID,
	})
}
