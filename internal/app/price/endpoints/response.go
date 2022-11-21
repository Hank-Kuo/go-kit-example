package endpoints

import (
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/Hank-Kuo/go-kit-example/internal/app/price/service"
	"github.com/Hank-Kuo/go-kit-example/pkg/response"
)

var (
	_ httptransport.Headerer    = (*SumResponse)(nil)
	_ httptransport.StatusCoder = (*SumResponse)(nil)
)

type SumResponse struct {
	Cost int64 `json:"cost"`
	Err  error `json:"-"`
}

type ExchangeResponse struct {
	Cost     int64  `json:"cost"`
	Currency string `json:"currency"`
}

func (r SumResponse) StatusCode() int {
	return http.StatusOK
}

func (r SumResponse) Headers() http.Header {
	h := http.Header{}
	h.Add("X-Api-Version", service.Version)
	return h
}

func (r SumResponse) Response() interface{} {
	return response.Response{Status: "success", Message: "sum two number", Data: r}
}
