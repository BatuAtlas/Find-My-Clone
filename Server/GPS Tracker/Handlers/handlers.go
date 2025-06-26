package handlers

import (
	connections "FindMy/GPSTracker/Connections"
	database "FindMy/GPSTracker/Database"
	model "FindMy/GPSTracker/Model"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

func PostAuth(c *fiber.Ctx) error {
	jsn := model.ParseJson(c.Body())
	token, ok := jsn["token"].(string)

	if !ok || len(token) != 64 {
		return c.Send(model.ParseResponse(false, 19, fiber.Map{"message": "invalid token"}))
	}

	user, expires := database.CheckToken(context.Background(), token) //c.context(): if the client has been disconnected, the database query will be cancelling. (it cant be because c.context() is fake http request, not websocket's context)

	if user == -1 {
		return c.Send(model.ParseResponse(false, 27, fiber.Map{"message": "wrong token"}))
	}

	if expires != nil && expires.Before(time.Now()) {
		database.RemoveToken(context.Background(), token) // context.Background(): even if the client disconnects, it will removed
		return c.Send(model.ParseResponse(false, 32, fiber.Map{"message": "expired token"}))
	}

	session := c.Context().UserValue("session").(map[string]interface{})

	session["auth.token"] = token
	session["user.id"] = user
	connections.Userconnection[user] = session

	return c.Send(model.ParseResponse(true, 40, fiber.Map{"message": "authorization succesfully completed"}))
}

func PostFeed(c *fiber.Ctx) error {
	session := c.Context().UserValue("session").(map[string]interface{})
	jsn := model.ParseJson(c.Body())

	lat, ok := jsn["lat"].(float64)
	lon, ok2 := jsn["lat"].(float64)
	//timestamp, ok3 := jsn["timestamp"].(string) /*now.toISOString();*/
	t := time.Now()
	//t, err := time.Parse(time.RFC3339, timestamp)

	if !ok || !ok2 /*||  !ok3 || err != nil*/ {
		return c.Send(model.ParseResponse(false, 72, fiber.Map{"message": "invalid form"}))
	}

	id := session["user.id"].(int64)
	err2 := database.FeedGPS(context.Background(), id, lat, lon, t)

	if err2 != nil {
		return c.Send(model.ParseResponse(true, 78, fiber.Map{"message": "failed :cc"}))
	}

	if session["subscribers"] != nil && len(session["subscribers"].([]int64)) != 0 {
		SubscribesPublish(session["subscribers"].([]int64), lat, lon, t, id)
	}

	return c.Send(model.ParseResponse(true, 81, fiber.Map{"message": "OK!"}))
}

func AuthCheckerMiddleware(c *fiber.Ctx) error {
	session := c.Context().UserValue("session").(map[string]interface{})
	if session["auth.token"] != nil {
		return c.Next()
	}

	if c.Path() == "/Auth" {
		return c.Next()
	}

	return c.Send(model.ParseResponse(false, 91, fiber.Map{"message": "please login!"}))
}

func GetLocation(c *fiber.Ctx) error {
	session := c.Context().UserValue("session").(map[string]interface{})
	lat, lon, timestamp := database.GetLocation(context.Background(), session["user.id"].(int64))

	if lat == 0 && lon == 0 {
		return c.Send(model.ParseResponse(false, 99, fiber.Map{"message": "failed :("}))
	}

	return c.Send(model.ParseResponse(true, 102, fiber.Map{"lat": lat, "lon": lon, "timestamp": timestamp}))
}

func PostInfo(c *fiber.Ctx) error {
	session := c.Context().UserValue("session").(map[string]interface{})
	jsn := model.ParseJson(c.Body())

	status, ok1 := jsn["status"].(string)
	battery, ok2 := jsn["battery"].(byte)
	isCharging, ok3 := jsn["isCharging"].(bool)
	event, ok4 := jsn["event"].(byte)

	if !ok1 || !ok2 || !ok3 || !ok4 {
		return c.Send(model.ParseResponse(false, 115, fiber.Map{"message": "invalid form"}))
	}

	err := database.ChangeInfo(context.Background(), session["user.id"].(int64), status, battery, isCharging, event)

	if err != nil {
		return c.Send(model.ParseResponse(true, 121, fiber.Map{"message": "failed :(cc"}))
	}
	return c.Send(model.ParseResponse(true, 123, fiber.Map{"message": "OK!"}))

}

func Subscribe(c *fiber.Ctx) error {
	session := c.Context().UserValue("session").(map[string]interface{})
	jsn := model.ParseJson(c.Body())

	users, ok1 := jsn["users"].([]int64) // users are usually all friends.

	if !ok1 || len(users) <= 0 {
		return c.Send(model.ParseResponse(false, 1215, fiber.Map{"message": "invalid form"}))
	}

	friends := database.GetFriends(context.Background(), session["user.id"].(int64))

	if !database.IsSubset(users, friends) { //is a all users friend?
		return c.Send(model.ParseResponse(false, 4128, fiber.Map{"message": "some users isn't your friend"}))
	}

	for _, user := range users {
		if connections.Userconnection[user] == nil {
			continue
		}

		subs := connections.Userconnection[user]["subscribers"]
		if subs == nil {
			subs = new([]int64)
		}

		subs = append(subs.([]int64), user)

	}

	return c.Send(model.ParseResponse(true, 14523, fiber.Map{"message": "OK!"}))
}
