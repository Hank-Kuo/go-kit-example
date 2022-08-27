package tracer

import (
	"fmt"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jconfig "github.com/uber/jaeger-client-go/config"
)

func NewJaeger(svcName string, url string, logger log.Logger) opentracing.Tracer {
	if url == "" {
		return opentracing.NoopTracer{}
	}

	tracer, closer, err := jconfig.Configuration{
		ServiceName: svcName,
		Sampler: &jconfig.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jconfig.ReporterConfig{
			LocalAgentHostPort: url,
			LogSpans:           true,
		},
	}.NewTracer()

	if err != nil {
		level.Error(logger).Log("msg", fmt.Sprintf("Failed to init Jaeger: %s", err))
		os.Exit(1)
	}
	defer closer.Close()

	opentracing.SetGlobalTracer(tracer)

	return tracer
}
