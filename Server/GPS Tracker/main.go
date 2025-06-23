package main

import (
	model "FindMy/GPSTracker/Model"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {

		return c.Send(model.HandleResponse(model.ParseResponse(true, "Hello World!")))
	})

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		defer c.Close()
		log.Println("Client connected")

		c.WriteMessage(websocket.TextMessage, []byte("Hello!"))
	}))

	//TODO: create a class, its works both websocket and http. if is websocket, gets a endpoint url and process etc.
	// * there will be multiple channel support (websocket and http)

	log.Fatal(app.Listen(":4215"))
}
