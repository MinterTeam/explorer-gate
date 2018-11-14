package errors

// New returns an error that formats as the given text.
func NewNodeError(text string, code uint) error {
	return &NodeError{text, code}
}

type NodeError struct {
	log  string
	code uint
}

func (e *NodeError) Error() string {
	return e.log
}

func (e *NodeError) Code() uint {
	return e.code
}
