package main

import (
	model "FindMy/GPSTracker/Model"
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/valyala/fasthttp"
)

var connections = make(map[string]bool)

func main() {
	app := fiber.New()

	/*app.Get("/", func(c *fiber.Ctx) error {
		return c.Send(model.ParseResponse(true, c.Body()))
	})*/

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}

		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		connections[c.Headers("Sec-Websocket-Key")] = true

		c.SetCloseHandler(func(code int, text string) error {
			connections[c.Headers("Sec-Websocket-Key")] = false
			return nil
		})

		for connections[c.Headers("Sec-Websocket-Key")] {
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

			ctx := fasthttp.RequestCtx{}
			ctx.Request.Header.SetMethod(jsn.Method)
			ctx.Request.SetRequestURI(jsn.Endpoint)
			ctx.Request.SetBody(jsn.Data)

			app.Handler()(&ctx)

			c.WriteMessage(websocket.TextMessage, ctx.Response.Body())
		}

	}))

	//TODO: create a class, its works both websocket and http. if is websocket, gets a endpoint url and process etc.
	// * there will be multiple channel support (websocket and http)
	log.Fatal(app.Listen(":4215"))
}
