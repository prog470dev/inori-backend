package controller

import (
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/prog470dev/inori-backend/model"
	"net/http"
	"net/url"
	"strconv"
)

type Offer struct {
	DB *sql.DB
}

type Resp struct {
	// この辺は埋め込みが使える？
	ID             int64         `json:"id"`
	DriverID       int64         `json:"driver_id"`
	Start          string        `json:"start"`
	Goal           string        `json:"goal"`
	DepartureTime  string        `json:"departure_time"`
	RiderCapacity  int64         `json:"rider_capacity"`
	ReservedRiders []model.Rider `json:"reserved_riders"`
}

func (d *Offer) GetDriverOffers(w http.ResponseWriter, r *http.Request) {

	// クエリパラメータの取得
	u, _ := url.Parse(r.URL.RequestURI())
	query := u.Query()
	if len(query["driver_id"]) == 0 {
		// 不適切なクエリ
		return
	}
	driverID, err := strconv.ParseInt(query["driver_id"][0], 10, 64)
	if err != nil {
		return
	}

	// ダミーデータ
	resps := []Resp{
		{
			ID:            1,
			DriverID:      driverID,
			Start:         "a",
			Goal:          "a",
			DepartureTime: "a",
			RiderCapacity: 3,
			ReservedRiders: []model.Rider{
				{
					10,
					"a",
					"a",
					"a",
					"a",
					"a",
					"a",
				},
				{
					20,
					"a",
					"a",
					"a",
					"a",
					"a",
					"a",
				},
			},
		},
		{
			ID:            2,
			DriverID:      driverID,
			Start:         "a",
			Goal:          "a",
			DepartureTime: "a",
			RiderCapacity: 2,
			ReservedRiders: []model.Rider{
				{
					10,
					"a",
					"a",
					"a",
					"a",
					"a",
					"a",
				},
			},
		},
	}

	JSON(w, http.StatusOK, struct {
		Offers []Resp `json:"offers"`
	}{
		Offers: resps,
	})
}

func (d *Offer) GetOfferDetail(w http.ResponseWriter, r *http.Request) {

	// パスパラメータの取得
	offerID, err := strconv.ParseInt(mux.Vars(r)["offer_id"], 10, 64)
	if err != nil {
		return
	}

	// ダミーデータ
	resp := Resp{
		ID:            offerID,
		DriverID:      1,
		Start:         "b",
		Goal:          "b",
		DepartureTime: "b",
		RiderCapacity: 2,
		ReservedRiders: []model.Rider{
			{
				10,
				"b",
				"b",
				"b",
				"b",
				"b",
				"b",
			},
		},
	}

	JSON(w, http.StatusOK, struct {
		Offer Resp `json:"offer"`
	}{
		Offer: resp,
	})
}

func (d *Offer) DeleteOffer(w http.ResponseWriter, r *http.Request) {
	// パスパラメータの取得
	offerID, err := strconv.ParseInt(mux.Vars(r)["offer_id"], 10, 64)
	if err != nil {
		return
	}

	// ダミーデータ
	resp := Resp{
		ID:            offerID,
		DriverID:      1,
		Start:         "c",
		Goal:          "c",
		DepartureTime: "c",
		RiderCapacity: 2,
		ReservedRiders: []model.Rider{
			{
				10,
				"c",
				"c",
				"c",
				"c",
				"c",
				"c",
			},
		},
	}

	// DELETEしたofferの情報を返す.
	JSON(w, http.StatusOK, struct {
		Offer Resp `json:"offer"`
	}{
		Offer: resp,
	})
}

func (d *Offer) CreateOffer(w http.ResponseWriter, r *http.Request) {
	// ダミーデータ
	resp := Resp{
		ID:            0,
		DriverID:      1,
		Start:         "c",
		Goal:          "c",
		DepartureTime: "c",
		RiderCapacity: 2,
		ReservedRiders: []model.Rider{
			{
				10,
				"c",
				"c",
				"c",
				"c",
				"c",
				"c",
			},
		},
	}

	// DELETEしたofferの情報を返す.
	JSON(w, http.StatusOK, struct {
		Offer Resp `json:"offer"`
	}{
		Offer: resp,
	})
}
