package main

import (
	database "FindMy/Responder/Database"
	"log"

	"github.com/gofiber/fiber/v2"
)

var connections = make(map[string]map[string]interface{})

func main() {
	database.Initialize()
	app := fiber.New()

	log.Fatal(app.Listen(":4216"))
}
