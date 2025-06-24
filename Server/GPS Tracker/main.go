package main

import (
	database "FindMy/GPSTracker/Database"
	handlers "FindMy/GPSTracker/Handlers"
	model "FindMy/GPSTracker/Model"
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/valyala/fasthttp"
)

var connections = make(map[string]map[string]interface{})

func main() {
	database.Initialize()
	app := fiber.New()
	wsRouter := fiber.New()

	wsRouter.Post("/Auth", handlers.PostAuth)
	wsRouter.Use("/User/:id", handlers.UserQueryMiddleware)

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}

		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		connections[c.Headers("Sec-Websocket-Key")] = map[string]interface{}{"isOpen": true}

		session := connections[c.Headers("Sec-Websocket-Key")]

		c.SetCloseHandler(func(code int, text string) error {
			session["isOpen"] = false
			return nil
		})

		for session["isOpen"].(bool) {
			var jsn model.Request
			err := c.ReadJSON(&jsn)

			if err != nil {
				var serialize, err2 = json.Marshal(jsn)
				var errmsg string
				if err2 != nil {
					errmsg = err.Error()
				} else {
					errmsg = string(serialize)
				}

				log.Print("Connection[" + c.Headers("Sec-Websocket-Key") + "] = " + errmsg)
				continue
			}

			ctx := new(fasthttp.RequestCtx)
			ctx.Request.SetRequestURI(jsn.Endpoint)
			ctx.Request.Header.SetMethod(jsn.Method)
			ctx.Request.SetBody(jsn.Data)
			ctx.SetUserValue("session", session)

			wsRouter.Handler()(ctx)

			c.WriteMessage(websocket.TextMessage, ctx.Response.Body())
		}

	}))

	log.Fatal(app.Listen(":4215"))
}
