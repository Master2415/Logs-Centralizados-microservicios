package models

import "time"

type Health struct {
	Checks []GeneralCheck `json:"checks"`
}

type GeneralCheck struct {
	Status  string         `json:"status"`
	Checks  []ServiceCheck `json:"chechks"`
	Version string         `json:"version"`
	Uptime  string         `json:"uptime"`
}

type ServiceCheck struct {
	Data   CheckData `json:"data"`
	Name   string    `json:"name"`
	Status string    `json:"status"`
}

type CheckData struct {
	From   time.Time `json:"from"`
	Status string    `json:"status"`
}
