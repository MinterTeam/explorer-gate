package domain

import "time"

type ExplorerErrorResponse struct {
	Error Error `json:"error"`
}

type ExplorerStatusResponse struct {
	Data ExplorerStatusData `json:"data"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ExplorerStatusData struct {
	AverageBlockTime      float64   `json:"averageBlockTime"`
	BipPriceUsd           float64   `json:"bipPriceUsd"`
	LatestBlockHeight     int       `json:"latestBlockHeight"`
	LatestBlockTime       time.Time `json:"latestBlockTime"`
	MarketCap             float64   `json:"marketCap"`
	TotalTransactions     int       `json:"totalTransactions"`
	TransactionsPerSecond float64   `json:"transactionsPerSecond"`
}
