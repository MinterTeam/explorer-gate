package errors

// New returns an error that formats as the given text.
func NewNodeError(text string, code int) error {
	return &NodeError{
		Message: text,
		Code:    code,
	}
}

type GateError struct {
	Error   string                 `json:"error"`
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

type TxResult struct {
}

type NodeErrorResponse struct {
	Error NodeError `json:"error"`
}

func (ner *NodeErrorResponse) GetNodeError() error {
	return &ner.Error
}

type NodeError struct {
	Message  string `json:"message"`
	Data     string `json:"data"`
	Code     int    `json:"code"`
	TxResult struct {
		Code int    `json:"code"`
		Log  string `json:"log"`
	} `json:"tx_result"`
}

func (e *NodeError) Error() string {
	return e.Message
}

func (e *NodeError) GetCode() int {
	return e.Code
}
func (e *NodeError) GetTxCode() int {
	return e.TxResult.Code
}
func (e *NodeError) GetMessage() string {
	return e.Message
}
func (e *NodeError) GetLog() string {
	return e.TxResult.Log
}

type OldNodeError struct {
	Log  string `json:"log"`
	Code int    `json:"code"`
}

func (e *OldNodeError) Error() string {
	return e.Log
}
