package model

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
	RiderID  int64 `json:"rider_id"`
	Schedule []Dem `json:"schedule"` // golangで集合は扱いにくいので配列を使用
}

type DemAggResp struct {
	Days [][]int `json:"days"`
}

// TABLE: demand_school, demand_home
type Demand struct {
	RiderID int64 `db:"rider_id" json:"rider_id"`
	Day     int64 `db:"day" json:"day"`
	Start   int64 `db:"start" json:"start"`
	End     int64 `db:"end" json:"end"`
}

// TABLE: demand_aggregate_school, demand_aggregate_home
type DemAgg struct {
	TimeZone int64 `db:"time_zone" json:"time_zone"`
	Value    int64 `db:"value" json:"value"`
}
