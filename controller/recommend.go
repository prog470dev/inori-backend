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

var weekDays = map[string]int{
	"Sunday":    0,
	"Monday":    1,
	"Tuesday":   2,
	"Wednesday": 3,
	"Thursday":  4,
	"Friday":    5,
	"Saturday":  6,
}

func (rec *Recommend) PushRecommend(dir int) error { //dir: 目的方向（school->0, home->1）
	offers, err := model.OffersAll(rec.DB)
	if err != nil {
		return err
	}

	offerCounts := calcAccumulation(offers)

	demands, err := model.DemandAll(rec.DB)
	if err != nil {
		return err
	}

	dRanges := convertDemand(demands, time.Now(), dir)

	for _, d := range dRanges {
		sum := calcRecommendOffers(offerCounts, d)
		if sum == 0 {
			continue
		}

		// プッシュ通知の送信
		token, err := model.TokenOneRider(rec.DB, d.riderID)
		if err != nil {
			return err
		}
		pushData := &PushRecommendData{
			To:          token.PushToken,
			Type:        "recommend_offer",
			MessageBody: fmt.Sprintf("予約希望時間帯に%d件のオファーがあります。", sum),
			Title:       fmt.Sprintf("予約希望時間帯に%d件のオファーがあります。", sum),
		}
		err = SendPushRecommendMessage(pushData)
		if err != nil {
			return err
		}
	}

	return nil
}

// タイムゾーンごとのオファー数を計算
func calcAccumulation(offers []model.Offer) [7][24 * model.Resolution]int {
	offerCounts := [7][24 * model.Resolution]int{}

	for _, offer := range offers {
		departureTime, err := time.Parse("2006-01-02 15:04:05", offer.DepartureTime)
		if err != nil {
			continue
		}

		timeSum := departureTime.Hour()*4 + (departureTime.Minute() / (60 / model.Resolution)) //境界値怪しい...

		offerCounts[weekDays[departureTime.Weekday().String()]][timeSum]++
	}

	// offerCountsを累積和に変換
	for i := 0; i < len(weekDays); i++ {
		for j := 0; j < len(offerCounts[i])-1; j++ {
			offerCounts[i][j+1] += offerCounts[i][j]
		}
	}

	return offerCounts
}

// demand をレコメンド数の計算が行える形に整形
func convertDemand(demands []*model.Demand, today time.Time, dir int) []demandRange {
	dRanges := []demandRange{}
	for _, d := range demands {
		// 方向によって見る曜日が異なる
		diff := 1
		if dir == 1 {
			diff = 0
		}
		// 行きの便 かつ 翌日の曜日（-> この条件なら自動的にユーザーがユニークになる）
		//TODO: オファーの方向（home/school）を時間だけでなく、明確にフィールドをテーブルに追加して判別
		if int(d.Dir) != dir || d.Day != int64((weekDays[today.Weekday().String()]+diff)%7) {
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

	return dRanges
}

// demand に対してレコメンドできるofferの数を計算
func calcRecommendOffers(offerCounts [7][24 * model.Resolution]int, d demandRange) int {
	var l int
	if d.start == 0 {
		l = 0
	} else {
		l = offerCounts[d.weekDay][d.start-1]
	}
	r := offerCounts[d.weekDay][d.end]

	return r - l // 希望時間帯に含まれるオファーの総数
}
