package model

import "encoding/json"

type Response struct {
	Status   bool        `json:"status"`
	Response interface{} `json:"response"`
}

type Error struct {
	Message   string `json:"message"`
	ErrorCode int    `json:"error"`
}

func ParseResponse(status bool, response interface{}) ([]byte, error) {
	data, err := json.Marshal(Response{Status: status, Response: response})
	if err != nil {
		return nil, err
	}

	return data, nil
}

//data, err := ParseResponse(status, err)
func HandleResponse(data []byte, err error) (response []byte) {
	if err != nil {
		data2, _ := ParseResponse(false, Error{Message: "Json Handle Error", ErrorCode: 48})
		response = data2
		return
	}

	response = data
	return
}
