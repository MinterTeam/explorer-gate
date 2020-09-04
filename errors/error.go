package errors

type GateError struct {
	ErrorString string                 `json:"error"`
	Code        string                 `json:"code"`
	Message     string                 `json:"message"`
	Details     map[string]interface{} `json:"details,omitempty"`
}

func (ge *GateError) Error() string {
	return ge.Message
}
