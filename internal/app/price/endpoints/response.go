package endpoints

import (
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/Hank-Kuo/go-kit-example/pkg/response"
)

var (
	_ httptransport.Headerer = (*SumResponse)(nil)
	// _ httptransport.Failer      = (*SumResponse)(nil)
	// _ endpoint.Failer           = SumResponse{}
	_ httptransport.StatusCoder = (*SumResponse)(nil)
)

type SumResponse struct {
	Res int64 `json:"res"`
	Err error `json:"-"`
}

type ExchangeResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type PriceResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func (r SumResponse) StatusCode() int {
	return http.StatusOK
}

func (r SumResponse) Headers() http.Header {
	return http.Header{}
}

func (r SumResponse) Response() interface{} {
	return response.Response{Data: r}
}

func (r SumResponse) Failed() error {
	return r.Err
}
