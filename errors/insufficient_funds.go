package errors

type InsufficientFundsError struct {
	log   string
	code  int32
	value string
	coin  string
}

func NewInsufficientFundsError(text string, code int32, val string, coin string) error {
	return &InsufficientFundsError{text, code, val, coin}
}

func (e *InsufficientFundsError) Error() string {
	return e.log
}

func (e *InsufficientFundsError) Code() int32 {
	return e.code
}

func (e *InsufficientFundsError) Value() string {
	return e.value
}

func (e *InsufficientFundsError) Coin() string {
	return e.coin
}
