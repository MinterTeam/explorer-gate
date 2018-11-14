package errors

type InsufficientFundsError struct {
	log   string
	code  uint
	value string
	coin  string
}

func NewInsufficientFundsError(text string, code uint, val string, coin string) error {
	return &InsufficientFundsError{text, code, val, coin}
}

func (e *InsufficientFundsError) Error() string {
	return e.log
}

func (e *InsufficientFundsError) Code() uint {
	return e.code
}

func (e *InsufficientFundsError) Value() string {
	return e.value
}

func (e *InsufficientFundsError) Coin() string {
	return e.coin
}
