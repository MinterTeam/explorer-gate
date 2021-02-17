package domain

type CoinEstimate struct {
	Value      string `json:"value"`
	Commission string `json:"commission"`
	SwapFrom   string `json:"swap_from"`
}
