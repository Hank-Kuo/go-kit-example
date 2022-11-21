package price

import (
	"errors"
)

var (
	ErrLessZeroe        = errors.New("less then zero")
	ErrIntOverflow      = errors.New("integer overflow")
	ErrCurrencyNotFound = errors.New("currency not found")
)
