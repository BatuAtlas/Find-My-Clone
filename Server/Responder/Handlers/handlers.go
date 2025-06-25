package handlers

import (
	database "FindMy/Responder/Database"
	model "FindMy/Responder/Model"
	"time"

	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

////////////////

// /User/:id
func UserQueryMiddleware(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "me" {
		c.Locals("id", c.Locals("user.id"))
		return c.Next()
	}

	i64, err := strconv.ParseInt(id, 10, 64)

	if err != nil || i64 == 0 || i64 == -1 {
		return c.Send(model.ParseResponse(false, 49, fiber.Map{"message": "invalid id"}))
	}

	bool := database.CheckUser(context.Background(), i64)

	if !bool {
		return c.Send(model.ParseResponse(false, 55, fiber.Map{"message": "user not found"}))
	}

	c.Locals("id", i64)
	return c.Next()
}

// User/:id
func GetUserFromID(c *fiber.Ctx) error {
	id := c.Locals("id").(int64)
	nickname, profilephoto := database.GetUser(context.Background(), id)

	return c.Send(model.ParseResponse(true, 4391, fiber.Map{"nickname": nickname, "profilephoto": profilephoto}))
}

// User/
func GetUser(c *fiber.Ctx) error {
	id := c.Locals("user.id").(int64)
	nickname, profilephoto := database.GetUser(context.Background(), id)

	return c.Send(model.ParseResponse(true, 4392, fiber.Map{"nickname": nickname, "profilephoto": profilephoto}))
}

////////////////
// User/
/*func UserMiddleware(c *fiber.Ctx) error {
	//it doesnt needed, because auth middleware already checked the user and writed id c.Locals("user.id")
}*/

func PostAuth(c *fiber.Ctx) error {
	jsn := model.ParseJson(c.Body())
	token, ok := jsn["token"].(string)

	if !ok || len(token) != 64 {
		return c.Send(model.ParseResponse(false, 422, fiber.Map{"message": "invalid token"}))
	}

	user, expires := database.CheckToken(context.Background(), token) //c.context(): if the client has been disconnected, the database query will be cancelling. (it cant be because c.context() is fake http request, not websocket's context)

	if user == -1 {
		return c.Send(model.ParseResponse(false, 481, fiber.Map{"message": "wrong token"}))
	}

	if expires != nil && expires.Before(time.Now()) {
		database.RemoveToken(context.Background(), token) // context.Background(): even if the client disconnects, it will removed

		c.ClearCookie("token")
		return c.Send(model.ParseResponse(false, 562, fiber.Map{"message": "expired token"}))
	}

	if expires == nil {
		c.Cookie(&fiber.Cookie{Name: "token", Value: token, HTTPOnly: true, Secure: true})
	} else {
		c.Cookie(&fiber.Cookie{Name: "token", Value: token, Expires: *expires, HTTPOnly: true, Secure: true})
	}

	return c.Send(model.ParseResponse(true, 630, fiber.Map{"message": "authorization succesfully completed"}))
}

func AuthCheckerMiddleware(c *fiber.Ctx) error {
	if c.Path() == "/Auth" {
		return c.Next()
	}

	token := c.Cookies("token")

	if len(token) != 64 {
		c.ClearCookie("token")
		return c.Send(model.ParseResponse(false, 701, fiber.Map{"message": "invalid token"}))
	}

	user, expires := database.CheckToken(context.Background(), token) //c.context(): if the client has been disconnected, the database query will be cancelling. (it cant be because c.context() is fake http request, not websocket's context)

	if user == -1 {
		c.ClearCookie("token")
		return c.Send(model.ParseResponse(false, 727, fiber.Map{"message": "wrong token"}))
	}

	if expires != nil && expires.Before(time.Now()) {
		database.RemoveToken(context.Background(), token) // context.Background(): even if the client disconnects, it will removed
		c.ClearCookie("token")
		return c.Send(model.ParseResponse(false, 832, fiber.Map{"message": "expired token"}))
	}

	c.Locals("user.id", user)
	return c.Next()
}

////////////////

func GetLocation(c *fiber.Ctx) error {
	id := c.Locals("id").(int64)
	lat, lon, timestamp := database.GetLocation(context.Background(), id)

	if timestamp.IsZero() {
		return c.Send(model.ParseResponse(true, 1311, fiber.Map{
			"message": "there isn't location data",
		}))
	}

	return c.Send(model.ParseResponse(true, 1312, fiber.Map{
		"lat":       lat,
		"lon":       lon,
		"timestamp": timestamp,
	}))
}

func GetInfo(c *fiber.Ctx) error {
	id := c.Locals("id").(int64)
	status, isCharging, battery, event, lastUpdate := database.GetInfo(context.Background(), id)

	if lastUpdate == nil {
		return c.Send(model.ParseResponse(true, 1465, fiber.Map{"message": "this user doesn't have info"}))
	}

	return c.Send(model.ParseResponse(true, 1469, fiber.Map{
		"status":     status,
		"isCharging": isCharging,
		"battery":    battery,
		"event":      event,
		"lastUpdate": lastUpdate,
	}))
}
