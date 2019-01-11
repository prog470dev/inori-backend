package controller

import (
	"database/sql"
	"fmt"
	"github.com/prog470dev/inori-backend/model"
	"time"
)

type demandRange struct {
	riderID int64
	weekDay int64
	start   int64
	end     int64
}

type Recommend struct {
	DB *sql.DB
}

func (rec *Recommend) PushRecommend() error {
	offers, err := model.OffersAll(rec.DB)
	if err != nil {
		return err
	}

	weekDays := map[string]int{
		"Sunday":    0,
		"Monday":    1,
		"Tuesday":   2,
		"Wednesday": 3,
		"Thursday":  4,
		"Friday":    5,
		"Saturday":  6,
	}

	// TODO: recommend と demandsをmodelから取得
	offerCounts := [7][24 * 4]int{}

	for _, offer := range offers {
		departureTime, err := time.Parse(time.RFC3339, offer.DepartureTime)
		if err != nil {
			continue
		}

		timeSum := departureTime.Hour()*4 + (departureTime.Minute() / (60 / 4)) //境界値怪しい...

		offerCounts[weekDays[departureTime.Weekday().String()]][timeSum]++
	}

	// offerCountsを累積和に変換
	for i := 0; i < len(offerCounts)-1; i++ {
		offerCounts[i+1] = offerCounts[i]
	}

	demands, err := model.DemandAll(rec.DB)
	if err != nil {
		return err
	}

	dRanges := []demandRange{}
	for _, d := range demands {
		// 行きの便 かつ 翌日の曜日（-> この条件なら自動的にユーザーがユニークになる）
		if d.Dir == 1 || d.Day != int64((weekDays[time.Now().Weekday().String()]+1)%7) {
			continue
		}
		dRange := demandRange{
			riderID: d.RiderID,
			weekDay: d.Day,
			start:   d.Start,
			end:     d.End,
		}
		dRanges = append(dRanges, dRange)
	}

	for _, d := range dRanges {
		var l int
		if l == 0 {
			l = 0
		} else {
			l = offerCounts[d.weekDay][d.start-1]
		}
		r := offerCounts[d.weekDay][d.end]

		sum := r - l // 希望時間帯に含まれるオファーの総数
		if sum == 0 {
			return nil
		}

		// プッシュ通知の送信
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
