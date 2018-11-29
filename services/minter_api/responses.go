package minter_api

type Response struct {
	Code uint    `json:"code"`
	Log  *string `json:"log"`
}

type SendTransactionResponse struct {
	Response
	Result SendTransactionResult `json:"result"`
}

type SendTransactionResult struct {
	Hash string `json:"hash"`
}

type TransactionResponse struct {
	Response
}
