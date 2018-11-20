package model

import (
	"database/sql"
)

type Token struct {
	ID        int64  `db:"id" json:"id"`
	Role      string `db:"role" json:"role"`
	RoleID    int64  `db:"role_id" json:"role_id"`
	PushToken string `db:"push_token" json:"push_token"`
}

// 新規レコードは追加、既存レコードは更新
func (t *Token) InsertOrUpdateToken(db *sql.DB) (sql.Result, error) {
	stmt, err := db.Prepare(`INSERT INTO tokens (role, role_id, push_token) values (?, ?, ?) ON DUPLICATE KEY UPDATE push_token = ?;`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.Exec(t.Role, t.RoleID, t.PushToken, t.PushToken)
}

func TokenOneDriver(db *sql.DB, driverID int64) (*Token, error) {
	token := &Token{}

	if err := db.QueryRow("SELECT * FROM tokens WHERE role = ? AND role_id = ? LIMIT 1", "driver", driverID).Scan(
		&token.ID,
		&token.Role,
		&token.RoleID,
		&token.PushToken,
	); err != nil {
		return nil, err
	}

	return token, nil
}

func TokenOneRider(db *sql.DB, driverID int64) (*Token, error) {
	token := &Token{}

	if err := db.QueryRow("SELECT * FROM tokens WHERE role = ? AND role_id = ? LIMIT 1", "rider", driverID).Scan(
		&token.ID,
		&token.Role,
		&token.RoleID,
		&token.PushToken,
	); err != nil {
		return nil, err
	}

	return token, nil
}
