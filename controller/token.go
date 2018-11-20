package controller

import (
	"database/sql"
	"encoding/json"
	"github.com/prog470dev/inori-backend/model"
	"net/http"
)

type Token struct {
	DB *sql.DB
}

/*
 プッシュトークン登録処理をほぼ同じ処理だが分けた理由:
 今後DriverとRicerで処理が変更されることを考慮
*/

func (t *Token) RegisterPushTokenDriver(w http.ResponseWriter, r *http.Request) {
	type rbody struct {
		ID    int64  `db:"id" json:"id"`
		Token string `db:"token" json:"token"`
	}

	var b rbody
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//TODO: Riderが登録済みであることを確認 (table: drivers)

	token := &model.Token{
		ID:        0,
		Role:      "driver",
		RoleID:    b.ID,
		PushToken: b.Token,
	}

	_, err := token.InsertOrUpdateToken(t.DB)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	JSON(w, http.StatusOK, b)
}

func (t *Token) RegisterPushTokenRider(w http.ResponseWriter, r *http.Request) {
	type rbody struct {
		ID    int64  `db:"id" json:"id"`
		Token string `db:"token" json:"token"`
	}

	var b rbody
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//TODO: Riderが登録済みであることを確認 (table: riders)

	token := &model.Token{
		ID:        0,
		Role:      "rider",
		RoleID:    b.ID,
		PushToken: b.Token,
	}

	_, err := token.InsertOrUpdateToken(t.DB)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	JSON(w, http.StatusOK, b)
}
