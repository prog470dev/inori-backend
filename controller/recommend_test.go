package controller

import (
	"fmt"
	"github.com/prog470dev/inori-backend/model"
	"testing"
	"time"
)

//func TestControllerRecommend(t *testing.T) {
//	conf := &db.Config{}
//	dbx, err := conf.Open("../dbconfig.yml") // testの場合は対象ファイルからの相対パス
//	if err != nil {
//		t.Fatalf("failed %s", err)
//	}
//
//	recommend := &Recommend{dbx}
//
//	err = recommend.PushRecommend()
//	if err != nil {
//		t.Fatalf("failed %s", err)
//	}
//
//}

/*
テストの想定状況
- 通知タイミングは金曜の夜
- 調べるオファーは翌日土曜のもの
*/

func TestCalcAccumulation(t *testing.T) {
	testDay := time.Date(2019, 1, 11, 22, 0, 0, 0, time.Local)

	nextDayStr := fmt.Sprintf("%04d-%02d-%02d",
		testDay.Add(24*time.Hour).Year(),
		testDay.Add(24*time.Hour).Month(),
		testDay.Add(24*time.Hour).Day())

	// 翌日のオファー
	offers := []model.Offer{
		{
			ID:            1,
			DriverID:      10,
			Start:         "foo",
			Goal:          "hoge",
			DepartureTime: nextDayStr + " 08:00:10",
			RiderCapacity: 2,
		},
		{
			ID:            2,
			DriverID:      20,
			Start:         "foo",
			Goal:          "hoge",
			DepartureTime: nextDayStr + " 10:00:10",
			RiderCapacity: 2,
		},
		{
			ID:            3,
			DriverID:      30,
			Start:         "foo",
			Goal:          "hoge",
			DepartureTime: nextDayStr + " 10:00:10",
			RiderCapacity: 2,
		},
	}

	// テストケース作成簡易化のための修正（DBではtime.RFC3339形式になっている）
	for i := 0; i < len(offers); i++ {
		t, err := time.Parse("2006-01-02 15:04:05", offers[i].DepartureTime)
		if err != nil {
			continue
		}
		offers[i].DepartureTime = t.Format(time.RFC3339)
		fmt.Println(offers[i].DepartureTime)
	}

	nextWeekDay := weekDays["Saturday"] // 2019-01-11(金曜日)の翌日
	tZone0 := model.Resolution * 8      // 1つのオファーが重なっている時間
	tZone1 := model.Resolution * 10     // 2つのオファーが重なっている時間

	offerCounts := calcAccumulation(offers)

	for i := 0; i < len(offerCounts); i++ {
		if i != nextWeekDay {
			continue
		}
		for j := 0; j < len(offerCounts[i]); j++ {
			if j == tZone0 {
				if offerCounts[i][j] != 1 {
					t.Fatalf("%d != %d, ", offerCounts[i][j], 1)
				}
			}
			if j == tZone1 {
				if offerCounts[i][j] != 3 {
					t.Fatalf("%d != %d, ", offerCounts[i][j], 3)
				}
			}
		}
	}
}

func TestConvertDemand(t *testing.T) {
	demands := []*model.Demand{
		{
			RiderID: 1,
			Day:     6,
			Dir:     0,
			Start:   0,
			End:     20,
		},
		{
			RiderID: 2,
			Day:     6,
			Dir:     0,
			Start:   35,
			End:     50,
		},
	}

	testDay := time.Date(2019, 1, 11, 22, 0, 0, 0, time.Local)
	demandRanges := convertDemand(demands, testDay, 0)

	for _, d := range demandRanges {
		if d.riderID == 1 {
			if 6 != d.weekDay {
				t.Fatalf("%d != %d, ", 6, d.weekDay)
			}
			if 0 != d.start {
				t.Fatalf("%d != %d, ", 0, d.start)
			}
			if 20 != d.end {
				t.Fatalf("%d != %d, ", 20, d.end)
			}
		}
		if d.riderID == 2 {
			if 6 != d.weekDay {
				t.Fatalf("%d != %d, ", 6, d.weekDay)
			}
			if 35 != d.start {
				t.Fatalf("%d != %d, ", 35, d.start)
			}
			if 50 != d.end {
				t.Fatalf("%d != %d, ", 50, d.end)
			}
		}
	}
}

func TestCalcRecommendOffers(t *testing.T) {
	testDay := time.Date(2019, 1, 11, 22, 0, 0, 0, time.Local)

	nextDayStr := fmt.Sprintf("%04d-%02d-%02d",
		testDay.Add(24*time.Hour).Year(),
		testDay.Add(24*time.Hour).Month(),
		testDay.Add(24*time.Hour).Day())

	offers := []model.Offer{
		{
			ID:            1,
			DriverID:      10,
			Start:         "foo",
			Goal:          "hoge",
			DepartureTime: nextDayStr + " 08:00:10",
			RiderCapacity: 2,
		},
		{
			ID:            2,
			DriverID:      20,
			Start:         "foo",
			Goal:          "hoge",
			DepartureTime: nextDayStr + " 10:00:10",
			RiderCapacity: 2,
		},
		{
			ID:            3,
			DriverID:      30,
			Start:         "foo",
			Goal:          "hoge",
			DepartureTime: nextDayStr + " 10:00:10",
			RiderCapacity: 2,
		},
	}

	demands := []*model.Demand{
		{
			RiderID: 1,
			Day:     6,
			Dir:     0,
			Start:   0,
			End:     20,
		},
		{
			RiderID: 2,
			Day:     6,
			Dir:     0,
			Start:   35,
			End:     50,
		},
	}

	offerCounts := calcAccumulation(offers)
	demandRanges := convertDemand(demands, testDay, 0)

	for _, d := range demandRanges {
		sum := calcRecommendOffers(offerCounts, d)
		if sum == 0 {
			continue
		}

		if d.riderID == 1 {
			if 0 != sum {
				t.Fatalf("%d != %d, ", 0, sum)
			}
		}
		if d.riderID == 2 {
			if 2 != sum {
				t.Fatalf("%d != %d, ", 2, sum)
			}
		}
	}
}
