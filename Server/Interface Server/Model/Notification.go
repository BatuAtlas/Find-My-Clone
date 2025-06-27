package model

type Notification struct {
	Type   int         `json:"type"`
	Params interface{} `json:"params"`
}
