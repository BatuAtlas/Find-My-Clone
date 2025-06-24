package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func RemoveToken(c context.Context, token string) (pgconn.CommandTag, error) {
	return pool.Exec(c, "DELETE FROM \"Authorization\" WHERE token = $1", token)
}

func CheckToken(c context.Context, token string) (user int64, expires *time.Time) {
	var row pgx.Row = pool.QueryRow(c, "SELECT \"user\", \"expires\" FROM \"Authorization\" WHERE token = $1", token)

	if row.Scan(&user, &expires) != nil {
		return -1, nil
	}
	return
}

func CheckUser(c context.Context, id int64) (exist bool) {
	var row pgx.Row = pool.QueryRow(c, "SELECT EXISTS (SELECT 1 FROM \"User\" WHERE id = $1)", id)

	if row.Scan(&exist) != nil {
		return false
	}
	return
}
