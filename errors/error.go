package errors

type GateErrorV1 struct {
	Code string                 `json:"code"`
	Log  string                 `json:"log"`
	Data map[string]interface{} `json:"data,omitempty"`
}

type GateError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

func (ge *GateError) Error() string {
	return ge.Message
}
