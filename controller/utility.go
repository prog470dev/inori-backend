package controller

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, code int, data interface{}) error {
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(data)
}

func NotFoundOrErr(w http.ResponseWriter, err error) error {
	if err == nil {
		return nil
	}
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	w.WriteHeader(http.StatusInternalServerError)
	return err
}
