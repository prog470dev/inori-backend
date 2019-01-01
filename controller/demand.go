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

	schedule := [7]model.Dem{}
	for i := 0; i < 7; i++ {
		schedule[i].Day = int64(i)
	}

	dems, err := model.DemandOne(d.DB, riderID)

	for _, dem := range dems {
		day := dem.Day
		dir := dem.Dir

		atom := model.DemAtom{
			Start: dem.Start,
			End:   dem.End,
		}

		if dir == 0 {
			schedule[day].School = atom
		} else { //dir == 1
			schedule[day].Home = atom
		}
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

	//TODO: 存在しない rider_id の場合は 400 を返す.
	//TODO: トランザクション
	//TODO: 不適切な時間指定のフィルタリング

	// 削除
	err := model.DeleteWithRider(d.DB, demRider.RiderID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 挿入
	for _, e := range demRider.Schedule {
		var dem model.Demand

		dem.RiderID = demRider.RiderID
		dem.Day = e.Day

		//school
		dem.Dir = 0
		dem.Start = e.School.Start
		dem.End = e.School.End
		if dem.Start != 0 && dem.End != 0 { //データなし
			_, err := dem.Insert(d.DB)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		//home
		dem.Dir = 1
		dem.Start = e.Home.Start
		dem.End = e.Home.End
		if dem.Start != 0 && dem.End != 0 { //データなし
			_, err = dem.Insert(d.DB)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}

	err = JSON(w, http.StatusOK, demRider)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (d *Demand) GetDemandAggregate(w http.ResponseWriter, r *http.Request) {
	days := model.DemAggResp{}
	days.Days = [7][24 * 4]int{}

	// 登下校の判定
	dir := mux.Vars(r)["dir"]
	if dir != "school" && dir != "home" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	demandAgg, err := model.DemandAggregate(d.DB, dir)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cnt := 0
	day := [24 * 4]int{}
	for _, agg := range demandAgg {
		day[cnt%(24*4)] = int(agg.Value)
		cnt++
		if cnt%(24*4) == 0 {
			days.Days[cnt/(24*4)-1] = day
			day = [24 * 4]int{}
		}
	}

	err = JSON(w, http.StatusOK, days)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
