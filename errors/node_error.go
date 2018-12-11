package errors

// New returns an error that formats as the given text.
func NewNodeError(text string, code int32) error {
	return &NodeError{text, code}
}

type NodeError struct {
	log  string
	code int32
}

func (e *NodeError) Error() string {
	return e.log
}

func (e *NodeError) Code() int32 {
	return e.code
}
