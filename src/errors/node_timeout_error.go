package errors

// New returns an error that formats as the given text.
func NewNodeTimeOutError(text string, code int) error {
	return &NodeError{Message: text, Code: code}
}

type NodeTimeOutError struct {
	log  string
	code int32
}

func (e *NodeTimeOutError) Error() string {
	return e.log
}

func (e *NodeTimeOutError) Code() int32 {
	return e.code
}
