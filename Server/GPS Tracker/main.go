package main

import (
	connections "FindMy/GPSTracker/Connections"
	database "FindMy/GPSTracker/Database"
	handlers "FindMy/GPSTracker/Handlers"
	model "FindMy/GPSTracker/Model"
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/valyala/fasthttp"
)

func main() {
	database.Initialize()
	app := fiber.New()
	wsRouter := fiber.New()

	wsRouter.Post("/Auth", handlers.PostAuth)
	wsRouter.Use("*", handlers.AuthCheckerMiddleware)
	//wsRouter.Use("/User/:id", handlers.UserQueryMiddleware)
	wsRouter.Post("/User/Location/Feed", handlers.PostFeed)
	//wsRouter.Get("/User/Location", handlers.GetLocation)
	wsRouter.Post("/User/Info", handlers.PostInfo)

	wsRouter.Post("/Subscribe/", handlers.Subscribe) // when client entered the android app ui, subscribe the gps traffic.
	wsRouter.Post("/UnSubscribe/", handlers.UnSubscribe)
	// its automatic calling on 20 minutes later when subscribed -i don't do it yet- or when client closed the android app ui

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}

		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		connections.Connection[c.Headers("Sec-Websocket-Key")] = map[string]interface{}{"isOpen": true}
		session := connections.Connection[c.Headers("Sec-Websocket-Key")]

		session["websocket.connection"] = c

		c.SetCloseHandler(func(code int, text string) error {
			if session["subscriptions"] != nil {
				for _, user := range session["subscriptions"].([]int64) {
					if connections.Userconnection[user] != nil && connections.Userconnection[user]["subscribers"] != nil {
						connections.Userconnection[user]["subscribers"] = database.RemoveValue(
							connections.Userconnection[user]["subscribers"].([]int64),
							session["user.id"].(int64))
					}
				}
			}

			if session["user.id"] != nil {
				delete(connections.Userconnection, session["user.id"].(int64))
			}

			delete(session, c.Headers("Sec-Websocket-Key"))
			return nil
		})

		for session["isOpen"].(*bool) == nil {
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
			ctx.Request.Header.Add("Sec-Websocket-Key", c.Headers("Sec-Websocket-Key"))
			ctx.Request.SetBody(jsn.Data)
			ctx.SetUserValue("session", session)

			wsRouter.Handler()(ctx)

			c.WriteMessage(websocket.TextMessage, ctx.Response.Body())
		}

	}))

	log.Fatal(app.Listen(":4215"))
}
