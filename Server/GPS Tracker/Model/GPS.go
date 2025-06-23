package model

type GPS struct {
	Latitude  string `json:"lat"`
	Longitude string `json:"lon"`
	Elevation int    `json:"elevation"`

	Timestamp int `json:"timestamp"`
	Speed     int `json:"speed"` //kmh (can be nil)
}
