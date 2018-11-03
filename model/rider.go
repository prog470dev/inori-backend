package model

type Rider struct {
	ID        int64  `db:"id" json:"id"`
	FirstName string `db:"first_name" json:"first_name"`
	LastName  string `db:"last_name" json:"last_name"`
	Grade     string `db:"grade" json:"grade"`
	Major     string `db:"major" json:"major"`
	Mail      string `db:"mail" json:"mail"`
	Phone     string `db:"phone" json:"phone"`
}
