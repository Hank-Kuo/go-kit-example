package utils

import (
	"errors"
	"net/http"

	customErr "github.com/Hank-Kuo/go-kit-example/pkg/errors"
)

var (
	ErrLessZeroe   = errors.New("less then zero")
	ErrIntOverflow = errors.New("integer overflow")
)

func CustomErrorEncoder(err customErr.Error) int {

	// switch {
	// case customErr.Contains(err, ErrTwoZeroes):
	// 	return http.StatusInternalServerError
	// }
	return http.StatusInternalServerError
}
