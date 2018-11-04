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

type Offer struct {
	DB *sql.DB
}

type OfferResp struct {
	Offer          model.Offer `json:"offer"`
	ReservedRiders []int64     `json:"reserved_riders"`

	/** 予備のデータ構造 **/
	//ID             int64         `json:"id"`
	//Driver         model.Driver  `json:"driver"`
	//Start          string        `json:"start"`
	//Goal           string        `json:"goal"`
	//DepartureTime  string        `json:"departure_time"`
	//RiderCapacity  int64         `json:"rider_capacity"`
	//ReservedRiders []model.Rider `json:"reserved_riders"`
}

func (o *Offer) GetDriverOffers(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.URL.RequestURI())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	query := u.Query()

	var offers []model.Offer

	if len(query["driver_id"]) == 0 { // すべて
		offers, err = model.OffersAll(o.DB)
		if NotFoundOrErr(w, err) != nil {
			return
		}
	} else { // 選択ドライバのみ
		driverID, err := strconv.ParseInt(query["driver_id"][0], 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		offers, err = model.OffersWithDriver(o.DB, driverID)
		if NotFoundOrErr(w, err) != nil {
			return
		}
	}

	resps := []OfferResp{}
	for _, off := range offers {
		reservations, err := model.ReservationsWithOffer(o.DB, off.ID)
		if NotFoundOrErr(w, err) != nil {
			return
		}

		// 予約中のライダーのリスト作成
		riders := []int64{}
		for _, rider := range reservations {
			riders = append(riders, rider.ID)
		}

		resp := OfferResp{
			Offer:          off,
			ReservedRiders: riders,
		}

		resps = append(resps, resp)
	}

	JSON(w, http.StatusOK, struct {
		Offers []OfferResp `json:"offers"`
	}{
		Offers: resps,
	})
}

func (o *Offer) GetOfferDetail(w http.ResponseWriter, r *http.Request) {
	offerID, err := strconv.ParseInt(mux.Vars(r)["offer_id"], 10, 64)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	offer, err := model.OfferOne(o.DB, offerID)
	if err != nil { //該当なしはerr=nil
		log.Println(err, "A")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	reservations, err := model.ReservationsWithOffer(o.DB, offerID)
	if err != nil { //該当なしはerr=nil
		log.Println(err, "B")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 予約中のライダーのリスト作成
	riders := []int64{}
	for _, rider := range reservations {
		riders = append(riders, rider.ID)
	}

	JSON(w, http.StatusOK, struct {
		Offer          model.Offer `json:"offer"`
		ReservedRiders []int64     `json:"reserved_riders"`
	}{
		Offer:          *offer,
		ReservedRiders: riders,
	})
}

func (d *Offer) DeleteOffer(w http.ResponseWriter, r *http.Request) {
	offerID, err := strconv.ParseInt(mux.Vars(r)["offer_id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	offer := model.Offer{}
	offer.ID = offerID
	_, err = offer.Delete(d.DB)
	if NotFoundOrErr(w, err) != nil {
		return
	}

	JSON(w, http.StatusOK, struct {
		ID int64 `json:"id"`
	}{
		ID: offerID,
	})
}

func (d *Offer) CreateOffer(w http.ResponseWriter, r *http.Request) {
	var offer model.Offer
	if err := json.NewDecoder(r.Body).Decode(&offer); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := offer.Insert(d.DB)
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
