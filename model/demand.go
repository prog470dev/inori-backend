package model

import (
	"database/sql"
)

const resolution = 4

type DemAtom struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
}

type Dem struct {
	Day    int64   `json:"day"`
	School DemAtom `json:"school"`
	Home   DemAtom `json:"home"`
}

type DemRider struct {
	RiderID  int64  `json:"rider_id"`
	Schedule [7]Dem `json:"schedule"` // golangで集合は扱いにくいので配列を使用
}

type DemAggResp struct {
	Days [7][24 * 4]int `json:"days"`
}

// TABLE: demand_school, demand_home
type Demand struct {
	RiderID int64 `db:"rider_id" json:"rider_id"`
	Day     int64 `db:"day" json:"day"`
	Dir     int64 `db:"dir" json:"dir"`
	Start   int64 `db:"start" json:"start"`
	End     int64 `db:"end" json:"end"`
}

// TABLE: demand_aggregate_school, demand_aggregate_home
type DemandAgg struct {
	TimeZone int64 `db:"time_zone" json:"time_zone"`
	Value    int64 `db:"value" json:"value"`
}

func DemandOne(db *sql.DB, riderID int64) ([]*Demand, error) {
	rows, err := db.Query("SELECT * FROM demand WHERE rider_id = ?", riderID)
	if err != nil {
		return nil, err
	}

	dems := []*Demand{}
	for rows.Next() {
		dem := &Demand{}
		err = rows.Scan(
			&dem.RiderID,
			&dem.Day,
			&dem.Dir,
			&dem.Start,
			&dem.End,
		)
		if err != nil {
			return nil, err
		}
		dems = append(dems, dem)
	}

	return dems, nil
}

// ino.demand_aggregate は別のテーブルだが、キャッシュ的な意味で使っている。
func DemandAggregate(db *sql.DB, dir string) ([]*DemandAgg, error) {
	table := "demand_aggregate_" + dir

	//TODO: bad smell!! (文字列連結)
	rows, err := db.Query("SELECT * FROM " + table)
	if err != nil {
		return nil, err
	}

	aggs := []*DemandAgg{}
	for rows.Next() {
		demandAgg := &DemandAgg{}
		err = rows.Scan(
			&demandAgg.TimeZone,
			&demandAgg.Value,
		)
		if err != nil {
			return nil, err
		}
		aggs = append(aggs, demandAgg)
	}

	return aggs, nil
}

// 対象ライダーのレコードを全削除
func DeleteWithRider(db *sql.DB, riderID int64) error {
	_, err := db.Exec("DELETE FROM demand WHERE rider_id = ?", riderID)
	if err != nil {
		return err
	}

	return nil
}

// レコードを新規挿入（DeleteWithRider の後）
func (d *Demand) Insert(db *sql.DB) (sql.Result, error) {
	result, err := db.Exec("INSERT INTO demand (rider_id, day, dir, start, end) values (?, ?, ?, ?, ?) ",
		d.RiderID, d.Day, d.Dir, d.Start, d.End)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// 定期的な（初期は高頻度）集計処理
func Aggregate(db *sql.DB) error {
	rows, err := db.Query("SELECT * FROM demand")
	if err != nil {
		return err
	}

	aggSchool := [24*resolution*7 + 1]int{} //最後に番兵が必要
	aggHome := [24*resolution*7 + 1]int{}

	for rows.Next() {
		dem := &Demand{}
		err = rows.Scan(
			&dem.RiderID,
			&dem.Day,
			&dem.Dir,
			&dem.Start,
			&dem.End,
		)
		if err != nil {
			return err
		}

		start := dem.Day*(24*resolution) + dem.Start
		end := dem.Day*(24*resolution) + dem.End + 1

		if dem.Dir == 0 {
			aggSchool[start]++
			aggSchool[end]--
		} else { // dem.Dir == 1
			aggHome[start]++
			aggHome[end]--
		}
	}

	// imos
	for i := 0; i < len(aggSchool)-1; i++ {
		if i > 0 {
			aggSchool[i] += aggSchool[i-1]
			aggHome[i] += aggHome[i-1]
		}
	}

	// ループの分離理由: こちらは前のインデックスの要素に影響を受けないので、並列化の可能性があるため。
	for i := 0; i < len(aggSchool)-1; i++ {
		_, err := db.Exec("UPDATE demand_aggregate_school SET value=? WHERE time_zone=?", aggSchool[i], i)
		if err != nil {
			return err
		}
		_, err = db.Exec("UPDATE demand_aggregate_home SET value=? WHERE time_zone=?", aggHome[i], i)
		if err != nil {
			return err
		}
	}

	return nil
}
