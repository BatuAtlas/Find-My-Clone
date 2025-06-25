package main

import (
	database "FindMy/Responder/Database"
	handlers "FindMy/Responder/Handlers"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	database.Initialize()
	app := fiber.New()

	app.Post("/Auth", handlers.PostAuth)
	app.Use("*", func(c *fiber.Ctx) error {
		c.Type("json")
		return c.Next()
	})
	app.Use("*", handlers.AuthCheckerMiddleware)
	app.Use("/User/:id/*", handlers.UserQueryMiddleware)

	app.Get("/User/:id/", handlers.GetUserFromID)
	app.Get("/User/:id/Location", handlers.GetLocation)
	app.Get("/User/:id/Info", handlers.GetInfo)
	app.Post("/User/:id/Info", handlers.PostInfo)

	log.Fatal(app.Listen(":4216"))
}
