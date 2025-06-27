package database

import (
	model "FindMy/InterfaceServer/Model"
	"context"
	"encoding/json"
	"fmt"
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

func CheckPassword(c context.Context, mail string, pass string) (match bool, id int64) {
	var row pgx.Row = pool.QueryRow(c, "SELECT \"user\" AS \"id\", CASE WHEN password = $1 THEN true ELSE false END AS \"match\" FROM \"Usersettings\" WHERE \"mail\" = $2;", pass, mail)

	if row.Scan(&id, &match) != nil {
		return false, 0
	}
	return
}

func GetToken(c context.Context, id int64) (token string, expires *time.Time) {
	var row pgx.Row = pool.QueryRow(c, "SELECT \"token\",\"expires\" FROM \"Authorization\" WHERE \"user\" = $1;", id)

	if row.Scan(&token, &expires) != nil {
		return "", nil
	}

	return
}

func CreateToken(c context.Context, id int64, expires *time.Time) (token string) {
	var row pgx.Row = pool.QueryRow(c, "INSERT INTO \"Authorization\" (\"user\", \"token\", \"expires\") VALUES ($1, encode(gen_random_bytes(32), 'hex'), $2) RETURNING \"token\"", id, expires)

	if row.Scan(&token) != nil {
		return ""
	}

	return
}

func UpdateNickname(c context.Context, id int64, newnickname string) (status bool) {
	_, err := pool.Exec(c, "UPDATE \"User\" SET \"nickname\"=$1 WHERE \"id\" = $2", newnickname, id)

	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

/*func GetNofitications(c context.Context, id int64) (texts []string) {
	var row pgx.Row = pool.QueryRow(c, "SELECT \"notifications\" FROM \"Usersettings\" where \"user\" = $1", id)

	if row.Scan(&texts) != nil {
		return []string{}
	}

	return
}*/

func GetPhoto(c context.Context, id int64) (addr string) {
	var row pgx.Row = pool.QueryRow(c, "SELECT \"notifications\" FROM \"Usersettings\" where \"user\" = $1", id)

	if row.Scan(&addr) != nil {
		return ""
	}

	return
}

func SetPhoto(c context.Context, id int64, photo string) (status bool) {
	_, err := pool.Exec(c, "UPDATE \"User\" SET \"profilephoto\"=$1 WHERE \"id\" = $2", photo, id)

	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func GetNofitications(c context.Context, id int64) []model.Notification {
	var row pgx.Row = pool.QueryRow(c, "SELECT \"notifications\" FROM \"Usersettings\" WHERE \"user\" = $1", id)

	var notifData []byte

	if row.Scan(&notifData) != nil {
		return []model.Notification{}
	}

	var notifications []model.Notification
	if len(notifData) > 0 {
		if err := json.Unmarshal(notifData, &notifications); err != nil {
			return []model.Notification{}
		}
	}

	return notifications
}

func AddNotification(c context.Context, id int64, notif model.Notification) bool {
	notifJson, err := json.Marshal([]model.Notification{notif}) //

	if err != nil {
		return false
	}

	_, err2 := pool.Exec(c, "UPDATE \"Usersettings\" SET \"notifications\"=\"notifications\" || $1::jsonb WHERE \"user\" = $2", notifJson, id)

	return err2 == nil
}

func SetNotification(c context.Context, id int64, notifications []model.Notification) bool {
	data, err := json.Marshal(notifications)
	if err != nil {
		return false
	}

	_, err = pool.Exec(c, `UPDATE "Usersettings" SET "notifications" = $1::jsonb WHERE "user" = $2`, data, id)

	return err == nil
}

//func RemoveNotification(c context.Context, id int64, notif model.Notification) bool {} -> GetNofitications and filter it, finally SetNotification

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

func AddFriend(c context.Context, user int64, friend int64) bool {
	_, err := pool.Exec(c, `
		UPDATE "User"
		SET friends = 
			CASE 
				WHEN friends IS NULL THEN ARRAY[$1]::bigint[]
				ELSE friends || $1
			END
		WHERE id = $2
	`, friend, user)

	return err == nil
}

func RemoveFriend(c context.Context, user int64, friend int64) bool {
	_, err := pool.Exec(c, `
		UPDATE "User"
		SET friends = array_remove(friends, $1)
		WHERE id = $2
	`, friend, user)

	return err == nil
}

func CreateUser(c context.Context, nickname string, mail string, password string) string {
	row := pool.QueryRow(c, `
WITH new_user AS (
  INSERT INTO "User" ("nickname")
  VALUES ($1) 
  RETURNING "id"
),
userinfo_insert AS (
  INSERT INTO "Userinfo" ("user", "lastUpdate")
  SELECT "id", NOW()
  FROM new_user
),
auth_insert AS (
  INSERT INTO "Authorization" ("user", "token")
  SELECT "id", encode(gen_random_bytes(32), 'hex')
  FROM new_user
  RETURNING "token", "user"
),
usersettings_insert AS (
  INSERT INTO "Usersettings"("user", "mail", "password") 
  SELECT "id", $2, $3
  FROM new_user
  RETURNING "user"
)
SELECT "user" AS "id", "token" FROM auth_insert`, nickname, mail, password)

	var userID int64
	var token string
	err := row.Scan(&userID, &token)

	if err != nil {
		return ""
	}

	return token
}

func CreateSignup(ctx context.Context, mail, pass, nickname string) string {
	var token string
	err := pool.QueryRow(ctx, `INSERT INTO "Signup" ("mail", "pass", "nickname", "token") VALUES ($1, $2, $3, encode(gen_random_bytes(32), 'hex')) RETURNING "token"`, mail, pass, nickname).Scan(&token)
	if err != nil {
		return ""
	}

	return token
}

func GetSignup(ctx context.Context, token string) (mail string, nickname string, pass string) {
	err := pool.QueryRow(ctx, `
		SELECT "mail", "nickname", "pass"
		FROM "Signup"
		WHERE "token" = $1
	`, token).Scan(&mail, &nickname, &pass)

	if err != nil {
		mail = ""
	}
	return
}

func RemoveSignup(ctx context.Context, token string) bool {
	_, err := pool.Exec(ctx, `DELETE FROM "Signup" WHERE "token" = $1;`, token)

	return err == nil
}
