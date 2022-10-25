package src

import "time"

type TrailerData struct {
	Latt               float64   `json:"latt"`
	Longt              float64   `json:"longt"`
	GpsTime            time.Time `json:"gpsTime"`
	OsTime             time.Time `json: "osTime"`
	Weight             float64   `json:"weight"`
	WeightStatus       bool      `json:"weightStatus"`
	ShuntVoltage       float64   `json:"shuntVoltage"`
	PowerSupplyVoltage float64   `json:"powerSupplyVoltage"`
	SerialNumber       string
}

type TrailerDataMqtt struct {
	Latt               float64 `json:"latt"`
	Longt              float64 `json:"longt"`
	GpsTime            string  `json:"gpsTime"`
	OsTime             string  `json: "osTime"`
	Weight             float64 `json:"weight"`
	WeightStatus       bool    `json:"weightStatus"`
	ShuntVoltage       float64 `json:"shuntVoltage"`
	PowerSupplyVoltage float64 `json:"powerSupplyVoltage"`
	SerialNumber       string
}
