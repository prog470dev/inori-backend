package controller

import (
	"database/sql"
	"fmt"
	"github.com/prog470dev/inori-backend/model"
)

type RecommendInterface interface {
	// offers テーブルへ
	GetOffersWithTime() ([]int, error) //すべてのoffer数を取得: () -> 書く時間の累積offer数
	// demand テーブルへ
	GetDemandWithTime() ([]demand, error) // 現在から24時間以内の需要データのリスト:() -> start, endのリスト
}

type demand struct {
	riderID int64
	start   int
	end     int
}

type Recommend struct {
	DB *sql.DB
}

func (rec *Recommend) PushRecommend() error {
	// TODO: recommendとdemandsをmodelから取得

	offerCounts := []int{} // タイムゾーンごとのoffer数（すべて）
	// offerCountsを累積和に変換
	for i := 0; i < len(offerCounts)-1; i++ {
		offerCounts[i+1] = offerCounts[i]
	}
	demands := []demand{} // 需要のリスト（ライダーはユニークであってほしい..., 今から近い方）

	//demands[i]について、[start, end] の範囲にあるオファーの数をカウント
	for _, d := range demands {
		var l int
		if l == 0 {
			l = 0
		} else {
			l = offerCounts[d.start-1]
		}
		r := offerCounts[d.end]
		sum := r - l // 希望時間帯に含まれるオファーの総数

		//カウントが1以上の場合は対象ライダー(d.ID == role_id, role == rider @tokens テーブル)向けにプッシュ通知を送信
		token, err := model.TokenOneRider(rec.DB, d.riderID)
		if err != nil {
			return err
		}
		pushData := &PushRecommendData{
			To:          token.PushToken,
			Type:        "recommend_offer",
			MessageBody: fmt.Sprint("予約希望時間帯に%d件のオファーがあります。", sum),
			Title:       fmt.Sprint("予約希望時間帯に%d件のオファーがあります。", sum),
		}
		err = SendPushRecommendMessage(pushData)
		if err != nil {
			return err
		}
	}

	return nil
}
