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
	//TODO: 意味的にURLにIDがほしいが、実装的にはボディにIDがいるので、URLにはいらない？
	_, err := strconv.ParseInt(mux.Vars(r)["driver_id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var driver model.Driver
	if err := json.NewDecoder(r.Body).Decode(&driver); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = driver.Update(d.DB)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	JSON(w, http.StatusOK, struct {
		Driver model.Driver `json:"driver"`
	}{
		Driver: driver,
	})
}

func (d *Driver) SignInDriver(w http.ResponseWriter, r *http.Request) {
	type Rb struct {
		Mail string `json:"mail"`
	}

	var rb Rb
	if err := json.NewDecoder(r.Body).Decode(&rb); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	driver, err := model.DriverOneWithMail(d.DB, rb.Mail)
	if NotFoundOrErr(w, err) != nil {
		return
	}

	// サインイン成功
	JSON(w, http.StatusOK, struct {
		Driver model.Driver `json:"driver"`
	}{
		Driver: *driver,
	})
}

func (d *Driver) SignUpDriver(w http.ResponseWriter, r *http.Request) {
	var driver model.Driver
	if err := json.NewDecoder(r.Body).Decode(&driver); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := driver.Insert(d.DB)
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

	JSON(w, http.StatusOK, struct {
		ID int64 `json:"id"`
	}{
		ID: id,
	})
}
