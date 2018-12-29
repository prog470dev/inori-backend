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

type DemAgg struct {
	TimeZone int64 `db:"time_zone" json:"time_zone"`
	Value    int64 `db:"value" json:"value"`
}

type DemAggResp struct {
	Days [][]int `json:"days"`
}
