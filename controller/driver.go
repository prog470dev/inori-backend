package controller

import (
	"bytes"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/olahol/go-imageupload"
	"github.com/prog470dev/inori-backend/model"
	"golang.org/x/crypto/bcrypt"
	"io"
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
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	driver, err := model.DriverOne(d.DB, driverID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // 不適切なuser_idと判断(本当はDBのエラーかも)
		return
	}

	_ = JSON(w, http.StatusOK, struct {
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

	_ = JSON(w, http.StatusOK, struct {
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

	_ = JSON(w, http.StatusOK, struct {
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

	_ = JSON(w, http.StatusOK, struct {
		ID int64 `json:"id"`
	}{
		ID: id,
	})
}

//TODO: driverとriderの重複を取り除く
func (d *Driver) PostImage(w http.ResponseWriter, r *http.Request) {
	driverID, err := strconv.ParseInt(mux.Vars(r)["driver_id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//存在しないdriverを弾く
	driver, err := model.DriverOne(d.DB, driverID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, fileHeader, err := r.FormFile("face_image")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	file, err := imageupload.Process(r, "face_image")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fileName := fileHeader.Filename

	// 名前のハッシュ化
	converted, err := bcrypt.GenerateFromPassword([]byte(fileName), 10)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fileName = hex.EncodeToString(converted[:]) + ".jpg"

	thumb, err := imageupload.ThumbnailJPEG(file, 512, 512, 85)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	re := io.NewSectionReader(bytes.NewReader(thumb.Data), 0, int64(len(thumb.Data)))
	f := sectionReadCloser{re}

	//アップロード
	url, err := AddFileToS3(fileName, f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//driver情報の更新（更新までサーバ側でやる）
	driver.ImageUrl = url
	_, err = driver.UpdateImage(d.DB)
	if err != nil {
		//
		fmt.Println(err)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_ = JSON(w, http.StatusOK, struct {
		ImageUrl string `json:"image_url"`
	}{
		ImageUrl: url,
	})
}
