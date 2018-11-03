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
	driverID, err := strconv.ParseInt(mux.Vars(r)["driver_id"], 10, 64)
	if err != nil {
		return
	}

	driver, err := model.DriverOne(d.DB, driverID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // 不適切なuser_idと判断(本当はDBのエラーかも)
		return
	}

	JSON(w, http.StatusOK, struct {
		Driver model.Driver `json:"driver"`
	}{
		Driver: *driver,
	})
}

func (d *Driver) UpdateDriver(w http.ResponseWriter, r *http.Request) {
	// model.Driverの型でレスポンスのボディを受ける
}

func (d *Driver) CreteDriver(w http.ResponseWriter, r *http.Request) {

}

func (d *Driver) SignInDriver(w http.ResponseWriter, r *http.Request) {

}

func (d *Driver) SignUpDriver(w http.ResponseWriter, r *http.Request) {

}
