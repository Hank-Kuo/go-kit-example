package http

import (
	"net/url"

	"github.com/Hank-Kuo/go-kit-example/pkg/errors"
)

func copyURL(base *url.URL, path string) *url.URL {
	next := *base
	next.Path = path
	return &next
}

func customErrorEncoder(err errors.Error) int {

	// switch {
	// case errors.Contains(err, http.Notfound):
	// 	return http.StatusForbidden
	// }
	return 0
}
