package model

import "database/sql"

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
	Days [][]int `json:"days"`
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

func DemandOne(db *sql.DB, riderID int64) ([]Demand, error) {
	rows, err := db.Query("SELECT * FROM demand_aggregate WHERE rider_id = ?", riderID)
	if err != nil {
		return nil, err
	}

	dems := []Demand{}
	for rows.Next() {
		dem := Demand{}
		err = rows.Scan(
			dem.RiderID,
			dem.Day,
			dem.Dir,
			dem.Start,
			dem.End,
		)
		if err != nil {
			return nil, err
		}
		dems = append(dems, dem)
	}

	return dems, nil
}

// ino.demand_aggregate は別のテーブルだが、キャッシュ的な意味で使っている。
func DemandAggregate(db *sql.DB) ([]DemandAgg, error) {
	rows, err := db.Query("SELECT * FROM demand_aggregate")
	if err != nil {
		return nil, err
	}

	aggs := []DemandAgg{}
	for rows.Next() {
		demandAgg := DemandAgg{}
		err = rows.Scan(
			demandAgg.TimeZone,
			demandAgg.Value,
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
func Aggregate() error {
	//TODO: impl
	return nil
}
