package endpoints

import (
	"context"
	"time"

	"golang.org/x/time/rate"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/tracing/zipkin"

	stdopentracing "github.com/opentracing/opentracing-go"
	stdzipkin "github.com/openzipkin/zipkin-go"
	"github.com/sony/gobreaker"

	"github.com/Hank-Kuo/go-kit-example/internal/app/price"
)

type Endpoints struct {
	SumEndpoint      endpoint.Endpoint
	ExchangeEndpoint endpoint.Endpoint
}

func New(svc price.PriceService, logger log.Logger, duration metrics.Histogram, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer) Endpoints {
	var sumEndpoint endpoint.Endpoint
	{
		method := "sum"
		sumEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 1))(sumEndpoint)
		sumEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(sumEndpoint)
		sumEndpoint = makeSumEndpoint(svc)
		sumEndpoint = opentracing.TraceServer(otTracer, method)(sumEndpoint)
		if zipkinTracer != nil {
			sumEndpoint = zipkin.TraceEndpoint(zipkinTracer, method)(sumEndpoint)
		}
		sumEndpoint = LoggingMiddleware(log.With(logger, "method", method))(sumEndpoint)
		sumEndpoint = InstrumentingMiddleware(duration.With("method", method))(sumEndpoint)
	}
	var exchangeEndpoint endpoint.Endpoint
	{
		method := "exchange"
		exchangeEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 1))(exchangeEndpoint)
		exchangeEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(exchangeEndpoint)
		exchangeEndpoint = makeExchangeEndpoint(svc)

		exchangeEndpoint = opentracing.TraceServer(otTracer, method)(exchangeEndpoint)
		if zipkinTracer != nil {
			exchangeEndpoint = zipkin.TraceEndpoint(zipkinTracer, method)(exchangeEndpoint)
		}
		exchangeEndpoint = LoggingMiddleware(log.With(logger, "method", method))(exchangeEndpoint)
		exchangeEndpoint = InstrumentingMiddleware(duration.With("method", method))(exchangeEndpoint)
	}
	return Endpoints{
		SumEndpoint:      sumEndpoint,
		ExchangeEndpoint: exchangeEndpoint,
	}

}

func makeSumEndpoint(svc price.PriceService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SumRequest)
		if err := req.validate(); err != nil {
			return nil, err
		}
		res, err := svc.Sum(ctx, req.Price, req.Fee)
		if err != nil {
			return nil, err
		}
		return SumResponse{Res: res}, nil
	}
}

func makeExchangeEndpoint(svc price.PriceService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SumRequest)
		if err := req.validate(); err != nil {
			return nil, err
		}
		res, err := svc.Sum(ctx, req.Price, req.Fee)
		if err != nil {
			return nil, err
		}
		return SumResponse{Res: res}, nil
	}
}

func (e Endpoints) Sum(ctx context.Context, price, fee int64) (int64, error) {
	resp, err := e.SumEndpoint(ctx, SumRequest{Price: price, Fee: fee})
	if err != nil {
		return 0, err
	}
	response := resp.(SumResponse)
	return response.Res, nil
}

func (e Endpoints) Exchange(ctx context.Context, cost int64, currency string) (int64, error) {
	resp, err := e.ExchangeEndpoint(ctx, SumRequest{Price: 100, Fee: 100})
	if err != nil {
		return 0, err
	}
	response := resp.(SumResponse)

	return response.Res, nil
}
