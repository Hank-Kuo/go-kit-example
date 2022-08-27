package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Hank-Kuo/go-kit-example/internal/app/price/endpoints"
	"github.com/Hank-Kuo/go-kit-example/pkg/response"
)

func decodeHTTPSumResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, response.JSONErrorDecoder(r)
	}
	var resp endpoints.SumResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

func decodeHTTPExchangeResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, response.JSONErrorDecoder(r)
	}
	var resp endpoints.ExchangeResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}
