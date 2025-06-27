package main

import (
	database "FindMy/InterfaceServer/Database"
	handlers "FindMy/InterfaceServer/Handlers"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {
	database.Initialize()
	app := fiber.New()

	if err := os.MkdirAll("./www/cdn", os.ModePerm); err != nil {
		log.Fatalf("Failed to create upload dir: %v", err)
	}

	app.Use("*", func(c *fiber.Ctx) error {
		c.Type("json")
		return c.Next()
	})
	app.Use("*", handlers.AuthCheckerMiddleware) // defines user.id local.

	app.Post("/Login", handlers.Login)
	app.Get("/User/me/photo", handlers.GetPhoto)
	app.Get("/User/me/nofitications", handlers.GetNofitications)
	app.Post("/User/me/nickname", handlers.SetNickname)

	app.Post("/User/me/photo", handlers.SetPhoto) // my best <3
	app.Post("/Signup", handlers.Signup)

	app.Post("/friends/add", handlers.AddFriend)
	app.Post("/friends/accept", handlers.AcceptFriend)
	app.Post("/friends/deny", handlers.DenyFriend)

	app.Get("/verify", handlers.Verify)
	log.Fatal(app.Listen(":4215"))
}
