package database

import (
	"context"
	"log"
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

func FeedGPS(c context.Context, userid int64, lat float64, lon float64, timestamp time.Time) error {
	_, err := pool.Exec(context.Background(), "UPDATE \"GPS\" SET \"lat\"=$2, \"lon\"=$3, \"timestamp\"=$4 WHERE \"user\" = $1", userid, lat, lon, timestamp)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func GetLocation(c context.Context, id int64) (lat float64, lon float64, timestamp time.Time) {
	var row pgx.Row = pool.QueryRow(c, "SELECT \"lat\", \"lon\", \"timestamp\" FROM \"GPS\" WHERE \"user\" = $1", id)

	if row.Scan(&lat, &lon, &timestamp) != nil {
		return 0, 0, time.Time{}
	}
	return
}

func ChangeInfo(c context.Context, user int64, status string, battery byte, isCharging bool, event byte) error {
	_, err := pool.Exec(context.Background(), "UPDATE \"UserInfo\" SET \"status\"=$2, \"battery\"=$3, \"isCharging\"=$4, \"event\": $5 WHERE \"user\" = $1", user, status, battery, isCharging, event)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
