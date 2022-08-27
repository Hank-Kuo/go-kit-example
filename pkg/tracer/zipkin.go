package tracer

import (
	"fmt"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
)

func NewZipkin(serviceName, httpPort, zipkinURL string, logger log.Logger) (zipkinTracer *zipkin.Tracer) {
	if zipkinURL != "" {
		var (
			err           error
			hostPort      = fmt.Sprintf("localhost:%s", httpPort)
			useNoopTracer = (zipkinURL == "")
			reporter      = zipkinhttp.NewReporter(zipkinURL)
		)
		zEP, _ := zipkin.NewEndpoint(serviceName, hostPort)
		zipkinTracer, err = zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(zEP), zipkin.WithNoopTracer(useNoopTracer))

		if err != nil {
			logger.Log("err", err)
			os.Exit(1)
		}
		if !useNoopTracer {
			logger.Log("tracer", "Zipkin", "type", "Native", "URL", zipkinURL)
		}
		return zipkinTracer
	}
	return nil
}
