package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func RemoveToken(c context.Context, token string) (pgconn.CommandTag, error) {
	return pool.Exec(c, "DELETE FROM \"Authorization\" WHERE \"token\" = $1", token)
}

func CheckToken(c context.Context, token string) (user int64, expires *time.Time) {
	var row pgx.Row = pool.QueryRow(c, "SELECT \"user\", \"expires\" FROM \"Authorization\" WHERE \"token\" = $1", token)

	if row.Scan(&user, &expires) != nil {
		return -1, nil
	}
	return
}

func CheckUser(c context.Context, id int64) (exist bool) {
	var row pgx.Row = pool.QueryRow(c, "SELECT EXISTS (SELECT 1 FROM \"User\" WHERE \"id\" = $1)", id)

	if row.Scan(&exist) != nil {
		return false
	}
	return
}
func GetUser(c context.Context, id int64) (string, string) {
	var row pgx.Row = pool.QueryRow(c, "SELECT \"nickname\", \"profilephoto\" FROM \"User\" WHERE \"id\" = $1", id)

	var nickname string
	var profilephoto *string

	err := row.Scan(&nickname, &profilephoto)
	if err != nil {
		return "", ""
	}
	if profilephoto == nil {
		return nickname, ""
	}
	return nickname, *profilephoto
}

func GetLocation(c context.Context, id int64) (lat float64, lon float64, timestamp time.Time) {
	var row pgx.Row = pool.QueryRow(c, "SELECT \"lat\", \"lon\", \"timestamp\" FROM \"GPS\" WHERE \"user\" = $1 ORDER BY \"timestamp\" DESC LIMIT 1", id)

	if row.Scan(&lat, &lon, &timestamp) != nil {
		return 0, 0, time.Time{}
	}
	return
}

func GetInfo(c context.Context, id int64) (*string, *bool, *int16, *int16, *time.Time) {
	var row pgx.Row = pool.QueryRow(c, "SELECT  \"status\", \"isCharging\", \"battery\", \"event\", \"lastUpdate\" FROM \"Userinfo\" WHERE \"user\" = $1", id)

	var (
		status     *string
		charging   *bool
		battery    *int16
		event      *int16
		lastupdate *time.Time
	)

	row.Scan(&status, &charging, &battery, &event, &lastupdate)
	return status, charging, battery, event, lastupdate
}

func ChangeInfo(c context.Context, user int64, status string, battery byte, isCharging bool, event byte) error {
	_, err := pool.Exec(context.Background(), "UPDATE \"UserInfo\" SET \"status\"=$2, \"battery\"=$3, \"isCharging\"=$4, \"event\"=$5, \"timestamp\"=$6 WHERE \"user\" = $1", user, status, battery, isCharging, event, time.Now())
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func GetFriends(c context.Context, user int64) []int64 {
	var row pgx.Row = pool.QueryRow(c, "SELECT \"friends\" FROM \"User\" WHERE \"id\" = $1", user)
	var arr *[]int64

	err := row.Scan(&arr)
	if err != nil {
		fmt.Println(err)
		return []int64{}
	}
	if arr == nil {
		return []int64{}

	}
	return *arr
}
