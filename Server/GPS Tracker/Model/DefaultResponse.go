package model

import (
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Status   bool        `json:"status"`
	Response interface{} `json:"response"`
}

func ParseResponse(status bool, response interface{}) []byte {
	data, err := json.Marshal(Response{Status: status, Response: response})
	if err != nil {
		data2 := ParseResponse(false, fiber.Map{"message": "Json Handle Error", "errorcode": 48})
		log.Println("Errcode 48 = " + err.Error())
		return data2
	}

	return data
}
