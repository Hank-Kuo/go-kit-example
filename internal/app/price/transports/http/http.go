package http

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/tracing/zipkin"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-zoo/bone"
	stdopentracing "github.com/opentracing/opentracing-go"
	stdzipkin "github.com/openzipkin/zipkin-go"
	"github.com/rs/cors"

	price "github.com/Hank-Kuo/go-kit-example/internal/app/price"
	priceEndpoint "github.com/Hank-Kuo/go-kit-example/internal/app/price/endpoints"
	"github.com/Hank-Kuo/go-kit-example/pkg/response"
)

func NewHTTPHandler(endpoints priceEndpoint.Endpoints, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer, logger log.Logger) http.Handler {
	options := []httptransport.ServerOption{
		// httptransport.ServerBefore(telepresence.HTTPToContext()),
		httptransport.ServerErrorEncoder(response.ErrorEncodeJSONResponse(customErrorEncoder)),
		httptransport.ServerErrorLogger(logger),
	}

	if zipkinTracer != nil {
		options = append(options, zipkin.HTTPServerTrace(zipkinTracer))
	}

	m := bone.New()
	sumHandler(m, endpoints, options, otTracer, logger)
	return cors.AllowAll().Handler(m)
}

func NewHTTPClient(instance string, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer, logger log.Logger) (price.PriceService, error) {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	u, err := url.Parse(instance)
	if err != nil {
		return nil, err
	}

	options := []httptransport.ClientOption{}

	if zipkinTracer != nil {
		options = append(options, zipkin.HTTPClientTrace(zipkinTracer))
	}

	// limiter := ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100))
	var sumEndpoint endpoint.Endpoint
	{
		sumEndpoint = httptransport.NewClient(
			"POST",
			copyURL(u, "/sum"),
			encodeHTTPSumRequest,
			decodeHTTPSumResponse,
			append(options, httptransport.ClientBefore(opentracing.ContextToHTTP(otTracer, logger)))...,
		).Endpoint()
		sumEndpoint = opentracing.TraceClient(otTracer, "Sum")(sumEndpoint)
		if zipkinTracer != nil {
			sumEndpoint = zipkin.TraceEndpoint(zipkinTracer, "Sum")(sumEndpoint)
		}

	}

	var exchangeEndpoint endpoint.Endpoint
	{
		exchangeEndpoint = httptransport.NewClient(
			"POST",
			copyURL(u, "/endpoint"),
			encodeHTTPSumRequest,
			decodeHTTPSumResponse,
			append(options, httptransport.ClientBefore(opentracing.ContextToHTTP(otTracer, logger)))...,
		).Endpoint()
		exchangeEndpoint = opentracing.TraceClient(otTracer, "Exchange")(sumEndpoint)
		if zipkinTracer != nil {
			exchangeEndpoint = zipkin.TraceEndpoint(zipkinTracer, "Exchange")(sumEndpoint)
		}

	}

	return priceEndpoint.Endpoints{
		SumEndpoint:      sumEndpoint,
		ExchangeEndpoint: exchangeEndpoint,
	}, nil
}

func sumHandler(m *bone.Mux, endpoints priceEndpoint.Endpoints, options []httptransport.ServerOption, otTracer stdopentracing.Tracer, logger log.Logger) {
	m.Post("/sum", httptransport.NewServer(
		endpoints.SumEndpoint,
		decodeHTTPSumRequest,
		response.EncodeJSONResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(otTracer, "Sum", logger)))...,
	))
}

func exchangeHandler(m *bone.Mux, endpoints priceEndpoint.Endpoints, options []httptransport.ServerOption, otTracer stdopentracing.Tracer, logger log.Logger) {
	m.Post("/sum", httptransport.NewServer(
		endpoints.SumEndpoint,
		decodeHTTPSumRequest,
		response.EncodeJSONResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(otTracer, "Sum", logger)))...,
	))
}
