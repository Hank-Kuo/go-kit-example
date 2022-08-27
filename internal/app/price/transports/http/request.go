package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Hank-Kuo/go-kit-example/internal/app/price/endpoints"
)

// sum
func encodeHTTPSumRequest(_ context.Context, r *http.Request, request interface{}) (err error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

func decodeHTTPSumRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.SumRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

// exchange
func encodeHTTPExchangeRequest(_ context.Context, r *http.Request, request interface{}) (err error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

func decodeHTTPExchangeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.ExchangeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}
