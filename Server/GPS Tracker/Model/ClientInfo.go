package model

type CLientInfo struct {
	Status     string `json:"status"` //status message (can be nil)
	IsCharging bool   `json:"ischarging"`
	Battery    int    `json:"battery"`
	Event      int    `json:"eventid"` //a number of event id (can be nil)
}

// clients doesn't send every time all of them. It is sends just modified ones. This struct just a full preview
