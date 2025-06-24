package handlers

import (
	database "FindMy/GPSTracker/Database"
	model "FindMy/GPSTracker/Model"
	"context"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func PostAuth(c *fiber.Ctx) error {
	jsn := model.ParseJson(c.Body())
	token, ok := jsn["token"].(string)

	if !ok || len(token) != 64 {
		return c.Send(model.ParseResponse(false, fiber.Map{"message": "invalid token"}))
	}

	user, expires := database.CheckToken(c.Context(), token) //c.context(): if the client has been disconnected, the database query will be cancelling.

	if user == -1 {
		return c.Send(model.ParseResponse(false, fiber.Map{"message": "wrong token"}))
	}

	if expires != nil && expires.Before(time.Now()) {
		database.RemoveToken(context.Background(), token) // context.Background(): even if the client disconnects, it will removed
		return c.Send(model.ParseResponse(false, fiber.Map{"message": "expired token"}))
	}

	session := c.Context().UserValue("session").(map[string]interface{})

	session["auth.token"] = token
	session["user.id"] = user

	return c.Send(model.ParseResponse(true, fiber.Map{"message": "authorization succesfully completed"}))
}

// /User/:id
func UserQueryMiddleware(c *fiber.Ctx) error {
	id := c.Params("id")
	i64, err := strconv.ParseInt(id, 10, 64)

	if err != nil || i64 == 0 || i64 == -1 {
		return c.Send(model.ParseResponse(false, fiber.Map{"message": "invalid id"}))
	}

	bool := database.CheckUser(c.Context(), i64)

	if !bool {
		return c.Send(model.ParseResponse(false, fiber.Map{"message": "user not found"}))
	}

	return c.Next()
}
