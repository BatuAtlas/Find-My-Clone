package model

import (
	"encoding/json"
	"log"
	"reflect"

	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Status   bool        `json:"status"`
	Code     int         `json:"code"`
	Response interface{} `json:"response"`
}

type Request struct {
	Endpoint string          `json:"endpoint"`
	Method   string          `json:"method"`
	Data     json.RawMessage `json:"data"`
}

func ParseResponse(status bool, code int, response interface{}) []byte {

	if reflect.TypeOf(response) == reflect.TypeOf([]byte{}) {
		json.Unmarshal(response.([]byte), &response)
	}

	data, err := json.Marshal(Response{Status: status, Code: code, Response: response})
	if err != nil {
		data2 := ParseResponse(false, 231, fiber.Map{"message": "Json Handle Error", "errorcode": 48})
		log.Println("Errcode 231 = " + err.Error())
		return data2
	}
	return data
}

func ParseJson(data []byte) (res map[string]interface{}) {
	json.Unmarshal(data, &res)
	return
}
