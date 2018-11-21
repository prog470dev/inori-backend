package controller

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/prog470dev/inori-backend/model"
	"log"
	"net/http"
	"net/url"
	"sort"
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

func (o *Offer) GetOffers(w http.ResponseWriter, r *http.Request) {
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

	//Offerの時系列順ソート
	sort.Sort(model.TimedOffer(offers))

	resps := []OfferResp{}
	for _, off := range offers {
		reservations, err := model.ReservationsWithOffer(o.DB, off.ID)
		if NotFoundOrErr(w, err) != nil {
			return
		}

		// 予約中のライダーのリスト作成
		riders := []int64{}
		for _, reservation := range reservations {
			riders = append(riders, reservation.RiderID)
		}

		//キャパシティオーバーの場合は非表示 (ドライば固定の場合は表示)
		if len(query["driver_id"]) == 0 && int(off.RiderCapacity) == len(reservations) {
			continue
		}

		// 時間文字列の変換
		t, err := SwitchTimeStrStyle(off.DepartureTime)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		off.DepartureTime = t

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
	if NotFoundOrErr(w, err) != nil {
		log.Println(err)
		return
	}

	reservations, err := model.ReservationsWithOffer(o.DB, offerID)
	if err != nil { //該当なしはerr=nil (通す)
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 予約中のライダーのリスト作成
	riders := []int64{}
	for _, reservation := range reservations {
		riders = append(riders, reservation.RiderID)
	}

	// 時間文字列の変換
	t, err := SwitchTimeStrStyle(offer.DepartureTime)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	offer.DepartureTime = t

	JSON(w, http.StatusOK, struct {
		Offer          model.Offer `json:"offer"`
		ReservedRiders []int64     `json:"reserved_riders"`
	}{
		Offer:          *offer,
		ReservedRiders: riders,
	})
}

// 注意：取得系のAPIは時間外のものはヒットしないが、削除系はヒットする。
// （理由は、ユーザ側の情報が古く、時間外のOfferが見えているときに削除できなくなるため。）
func (o *Offer) DeleteOffer(w http.ResponseWriter, r *http.Request) {
	offerID, err := strconv.ParseInt(mux.Vars(r)["offer_id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	offer := &model.Offer{}
	// delete文は対象レコードの有無に関する情報は返さないため、事前にselectを実行
	offer, err = model.OfferOneWithoutTime(o.DB, offerID)
	log.Println(err)
	if NotFoundOrErr(w, err) != nil {
		log.Println(err)
		return
	}

	// 関連reservationの情報を取得 (プッシュ通知のため)
	// 実際の通知を削除後にやる理由：Expoサーバのトラブルでreservationが削除されないのを防ぐため
	reservations, err := model.ReservationsWithOffer(o.DB, offer.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 関連reservationの削除
	_, err = model.DeleteOfferReservation(o.DB, int(offer.ID))
	if err != nil && err != sql.ErrNoRows {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// offer の削除 (関連reservationは削除済み)
	_, err = offer.Delete(o.DB)
	if NotFoundOrErr(w, err) != nil {
		return
	}

	//プッシュ通知 (複数のライダ向け)
	for _, reserve := range reservations {
		token, err := model.TokenOneRider(o.DB, reserve.ID)
		if err != nil {
			JSON(w, http.StatusOK, struct {
				ID      int64  `json:"id"`
				Message string `json:"message"`
			}{
				ID:      offer.ID,
				Message: "Failed to send push notification.",
			})
			return
		}
		pushData := &PushData{
			To:          token.PushToken,
			Type:        "canceled_offer",
			OfferID:     offer.ID,
			MessageBody: "予約済みのオファーがキャンセルされました。",
			Title:       "予約済みのオファーがキャンセルされました。",
		}
		err = SendPushMessage(pushData)
		if err != nil {
			JSON(w, http.StatusOK, struct {
				ID      int64  `json:"id"`
				Message string `json:"message"`
			}{
				ID:      offer.ID,
				Message: "Failed to send push notification.",
			})
			return
		}
	}

	JSON(w, http.StatusOK, struct {
		ID int64 `json:"id"`
	}{
		ID: offer.ID,
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
