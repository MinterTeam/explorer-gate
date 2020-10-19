package errors

type GateError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

func (ge *GateError) Error() string {
	return ge.Message
}
