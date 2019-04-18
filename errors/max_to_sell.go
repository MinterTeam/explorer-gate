package errors

type MaximumValueToSellReachedError struct {
	log  string
	code int32
	want string
	need string
}

func NewMaximumValueToSellReachedError(text string, code int32, want string, need string) error {
	return &MaximumValueToSellReachedError{text, code, want, need}
}

func (e *MaximumValueToSellReachedError) Error() string {
	return e.log
}

func (e *MaximumValueToSellReachedError) Code() int32 {
	return e.code
}

func (e *MaximumValueToSellReachedError) Want() string {
	return e.want
}

func (e *MaximumValueToSellReachedError) Need() string {
	return e.need
}
