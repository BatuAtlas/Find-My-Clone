package handlers

import (
	connections "FindMy/GPSTracker/Connections"
	model "FindMy/GPSTracker/Model"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func SubscribesPublish(subscribers []int64, lat float64, lon float64, timestamp time.Time, user int64) {
	for _, subscribe := range subscribers {
		c := connections.Userconnection[subscribe]["websocket.connection"]

		if c == nil {
			continue
		}

		c.(*websocket.Conn).WriteMessage(websocket.TextMessage, []byte(model.ParseResponse(true, 14308, fiber.Map{
			"lat":       lat,
			"lon":       lon,
			"timestamp": timestamp,
			"user":      user,
		})))
	}
}
