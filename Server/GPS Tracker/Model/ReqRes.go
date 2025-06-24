package model

import (
	"encoding/json"
	"log"
	"reflect"

	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Status   bool        `json:"status"`
	Response interface{} `json:"response"`
}

type Request struct {
	Endpoint string          `json:"endpoint"`
	Method   string          `json:"method"`
	Data     json.RawMessage `json:"data"`
}

func ParseResponse(status bool, response interface{}) []byte {

	if reflect.TypeOf(response) == reflect.TypeOf([]byte{}) {
		json.Unmarshal(response.([]byte), &response)
	}

	data, err := json.Marshal(Response{Status: status, Response: response})
	if err != nil {
		data2 := ParseResponse(false, fiber.Map{"message": "Json Handle Error", "errorcode": 48})
		log.Println("Errcode 48 = " + err.Error())
		return data2
	}
	return data
}

func ParseJson(data []byte) (res map[string]interface{}) {
	json.Unmarshal(data, &res)
	return
}
