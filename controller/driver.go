package controller

import (
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/prog470dev/inori-backend/model"
	"net/http"
	"strconv"
)

type Driver struct {
	DB *sql.DB
}

func (d *Driver) GetDriverDetail(w http.ResponseWriter, r *http.Request) {

	// パスパラメータの取得
	driverID, err := strconv.ParseInt(mux.Vars(r)["driver_id"], 10, 64)
	if err != nil {
		return
	}

	//ダミーデータ
	driver := model.Driver{
		ID:        driverID,
		FirstName: "d",
		LastName:  "d",
		Grade:     "d",
		Major:     "d",
		Mail:      "d",
		Phone:     "d",
		CarColor:  "d",
		CarNumber: "d",
	}

	JSON(w, http.StatusOK, struct {
		Driver model.Driver `json:"driver"`
	}{
		Driver: driver,
	})
}
