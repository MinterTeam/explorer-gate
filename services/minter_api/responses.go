package minter_api

type Response struct {
	JSONrpc string     `json:"jsonrpc"`
	Id      string     `json:"id"`
	Error   *ErrorData `json:"error"`
}

type ErrorData struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

type SendTransactionResponse struct {
	Response
	Result *SendTransactionResult `json:"result"`
}

type SendTransactionResult struct {
	Code int32  `json:"code"`
	Data string `json:"data"`
	Log  string `json:"log"`
	Hash string `json:"hash"`
}

type TransactionResponse struct {
	Response
}
