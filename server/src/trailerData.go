package src

import "time"

type TrailerData struct {
	Latt               float64   `json:"latt"`
	Longt              float64   `json:"longt"`
	GpsTime            time.Time `json:"GpsTime"`
	OsTime             time.Time `json: "OsTime"`
	Weight             float64   `json:"Weight"`
	WeightStatus       int8      `json:"WeightStatus"`
	ShuntVoltage       float64   `json:"ShuntVoltage"`
	PowerSupplyVoltage float64   `json:"PowerSupplyVoltage"`
	SerialNumber       string    `json: "SerialNumber"`
	CpuTemp            float64   `json: "CpuTemp"`
}

type TrailerDataMqtt struct {
	Latt               float64 `json:"latt"`
	Longt              float64 `json:"longt"`
	GpsTime            string  `json:"GpsTime"`
	OsTime             string  `json: "OsTime"`
	Weight             float64 `json:"Weight"`
	WeightStatus       int8    `json:"WeightStatus"`
	ShuntVoltage       float64 `json:"ShuntVoltage"`
	PowerSupplyVoltage float64 `json:"PowerSupplyVoltage"`
	SerialNumber       string  `json: "SerialNumber"`
	CpuTemp            float64 `json: "CpuTemp"`
}
