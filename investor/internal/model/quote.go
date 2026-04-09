package model

import "time"

type Quote struct {
	Symbol string
	Price  float64
	Time   time.Time
}

type ISSResponse struct {
	MarketData struct {
		Columns []string        `json:"columns"`
		Data    [][]interface{} `json:"data"`
	} `json:"marketdata"`
}
