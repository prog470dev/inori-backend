package controller

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
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

func SwitchTimeStrStyle(timeStr string) (string, error) {
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return "", err
	}
	return t.Format("2006-01-02 15:04:05"), nil
}

type PushData struct {
	To          string
	Type        string
	OfferID     int64
	MessageBody string
	Title       string
}

func SendPushMessage(pushData *PushData) error {
	type data struct {
		Type        string `json:"type"`
		OfferID     int64  `json:"offer_id"`
		MessageBody string `json:"message_body"`
	}

	body := struct {
		To        string `json:"to"`
		Data      data   `json:"data"`
		Title     string `json:"title"`
		Body      string `json:"body"`
		Sound     string `json:"sound"`
		Badge     int64  `json:"badge"`
		ChannelId string `json:"channelId"` //TODO: string と nil を両立
	}{
		To: pushData.To,
		Data: data{
			Type:        pushData.Type,
			OfferID:     pushData.OfferID,
			MessageBody: pushData.MessageBody,
		},
		Title:     pushData.Title, //TODO: 空だとエラー
		Body:      "",
		Sound:     "default",
		Badge:     1,
		ChannelId: "null",
	}

	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	//テスト
	var buf bytes.Buffer
	buf.Write(b)
	log.Println(buf.String())

	req, err := http.NewRequest(
		"POST",
		"https://exp.host/--/api/v2/push/send",
		bytes.NewBuffer(b),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	log.Println(resp.Status)

	return nil
}
