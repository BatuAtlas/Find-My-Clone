package handlers

import (
	database "FindMy/GPSTracker/Database"
	model "FindMy/GPSTracker/Model"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func PostAuth(c *fiber.Ctx) error {
	jsn := model.ParseJson(c.Body())
	token, ok := jsn["token"].(string)

	if !ok || len(token) != 64 {
		return c.Send(model.ParseResponse(false, 19, fiber.Map{"message": "invalid token"}))
	}

	pool := database.GetPool()
	fmt.Println(pool)
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

	return c.Send(model.ParseResponse(true, 40, fiber.Map{"message": "authorization succesfully completed"}))
}

// /User/:id
func UserQueryMiddleware(c *fiber.Ctx) error {
	id := c.Params("id")
	i64, err := strconv.ParseInt(id, 10, 64)

	if err != nil || i64 == 0 || i64 == -1 {
		return c.Send(model.ParseResponse(false, 49, fiber.Map{"message": "invalid id"}))
	}

	bool := database.CheckUser(context.Background(), i64)

	if !bool {
		return c.Send(model.ParseResponse(false, 55, fiber.Map{"message": "user not found"}))
	}

	return c.Next()
}

func PostFeed(c *fiber.Ctx) error {
	session := c.Context().UserValue("session").(map[string]interface{})
	jsn := model.ParseJson(c.Body())

	lat, ok := jsn["lat"].(float64)
	lon, ok2 := jsn["lat"].(float64)
	timestamp, ok3 := jsn["timestamp"].(string) /*now.toISOString();*/

	t, err := time.Parse(time.RFC3339, timestamp)

	if !ok || !ok2 || !ok3 || err != nil {
		return c.Send(model.ParseResponse(false, 72, fiber.Map{"message": "invalid form"}))
	}

	err2 := database.FeedGPS(context.Background(), session["user.id"].(int64), lat, lon, t)

	if err2 != nil {
		return c.Send(model.ParseResponse(true, 78, fiber.Map{"message": "failed :cc"}))
	}

	return c.Send(model.ParseResponse(true, 81, fiber.Map{"message": "OK!"}))

}

func AuthCheckerMiddleware(c *fiber.Ctx) error {
	session := c.Context().UserValue("session").(map[string]interface{})
	if session["auth.token"] != nil {
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
