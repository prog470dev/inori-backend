package controller

import (
	"database/sql"
	"encoding/json"
	"github.com/prog470dev/inori-backend/model"
	"net/http"
)

type User struct {
	DB *sql.DB
}

func (u *User) Get(w http.ResponseWriter, r *http.Request) {
	users, err := model.Select(u.DB)
	if err != nil {
		return
	}

	JSON(w, http.StatusOK, struct {
		Users []model.User `json:"users"`
	}{
		Users: users,
	})
}

func (u *User) Post(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return
	}

	result, err := user.Insert(u.DB)
	if err != nil {
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		return
	}

	JSON(w, http.StatusOK, struct {
		UserID int64 `json:"id"`
	}{
		UserID: id,
	})
}
