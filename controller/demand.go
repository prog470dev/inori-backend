package controller

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/prog470dev/inori-backend/model"
	"net/http"
	"strconv"
)

type Demand struct {
	DB *sql.DB
}

func (d *Demand) GetDemandRider(w http.ResponseWriter, r *http.Request) {
	riderID, err := strconv.ParseInt(mux.Vars(r)["rider_id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	schedule := []model.Dem{}

	for i := 0; i < 7; i++ {
		dema1 := model.DemAtom{
			Start: 10,
			End:   15,
		}

		startUp := 0
		startDown := 0
		if i%2 == 0 {
			startUp = 20
			startDown = 25
		}

		dema2 := model.DemAtom{
			Start: int64(startUp),
			End:   int64(startDown),
		}

		dem := model.Dem{
			Day:    (int64(i)),
			School: dema1,
			Home:   dema2,
		}
		schedule = append(schedule, dem)
	}

	demRider := model.DemRider{
		RiderID:  riderID,
		Schedule: schedule,
	}

	err = JSON(w, http.StatusOK, demRider)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (d *Demand) ResisterDemandRider(w http.ResponseWriter, r *http.Request) {
	var demRider model.DemRider
	if err := json.NewDecoder(r.Body).Decode(&demRider); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := JSON(w, http.StatusOK, demRider)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (d *Demand) GetDemandAggregate(w http.ResponseWriter, r *http.Request) {
	days := model.DemAggResp{}

	for i := 0; i < 7; i++ {
		day := []int{}
		for i := 0; i < 96; i++ {
			day = append(day, 1)
		}
		days.Days = append(days.Days, day)
	}

	err := JSON(w, http.StatusOK, days)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
