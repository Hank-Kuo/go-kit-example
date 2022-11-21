package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	kitgrpc "github.com/go-kit/kit/transport/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/kelseyhightower/envconfig"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/Hank-Kuo/go-kit-example/internal/app/price/endpoints"
	priceEndpoint "github.com/Hank-Kuo/go-kit-example/internal/app/price/endpoints"
	priceService "github.com/Hank-Kuo/go-kit-example/internal/app/price/service"
	priceGrpcTransport "github.com/Hank-Kuo/go-kit-example/internal/app/price/transports/grpc"
	priceHttpTransport "github.com/Hank-Kuo/go-kit-example/internal/app/price/transports/http"
	pb "github.com/Hank-Kuo/go-kit-example/pb/price"
	"github.com/Hank-Kuo/go-kit-example/pkg/tracer"
)

type Config struct {
	ServiceName string `envconfig:"QS_SERVICE_NAME" default:"price"`
	ServiceHost string `envconfig:"QS_SERVICE_HOST" default:"localhost"`
	LogLevel    string `envconfig:"QS_LOG_LEVEL" default:"error"`
	HttpPort    string `envconfig:"QS_HTTP_PORT" default:"8180"`
	GrpcPort    string `envconfig:"QS_GRPC_PORT" default:"8181"`
	ZipkinURL   string `envconfig:"QS_ZIPKIN_URL"`
	JaegerURL   string `envconfig:"QS_JAEGER_URL"`
}

func main() {
	fmt.Println("==> Starting price service")
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	}

	var cfg Config
	err := envconfig.Process("qs", &cfg)
	if err != nil {
		level.Error(logger).Log("err", err)
		os.Exit(1)
	}

	logger = level.NewFilter(logger, level.AllowInfo())
	logger = log.With(logger, "service", cfg.ServiceName)
	logger = log.With(logger, "caller", log.DefaultCaller)
	level.Info(logger).Log("version", priceService.Version, "commitHash", priceService.CommitHash, "buildTimeStamp", priceService.BuildTimeStamp)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// endpoint-level metrics.
	var duration metrics.Histogram
	{
		duration = prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "example",
			Subsystem: "pricesvc",
			Name:      "request_duration_seconds",
			Help:      "Request duration in seconds.",
		}, []string{"method", "success"})
	}

	// service-level metrics.
	var cout1, cout2 metrics.Counter
	{
		cout1 = prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "example",
			Subsystem: "pricesvc",
			Name:      "integers_summed",
			Help:      "Total count of integers summed via the Sum method.",
		}, []string{})
		cout2 = prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "example",
			Subsystem: "addsvc",
			Name:      "characters_concatenated",
			Help:      "Total count of characters concatenated via the Concat method.",
		}, []string{})
	}

	http.DefaultServeMux.Handle("/metrics", promhttp.Handler())

	// jaeger
	otTracer := tracer.NewJaeger(cfg.ServiceName, cfg.JaegerURL, logger)

	// zipkin
	zipkin := tracer.NewZipkin(cfg.ServiceName, cfg.HttpPort, cfg.ZipkinURL, logger)

	// service
	srv := priceService.NewService(logger, cout1, cout2)
	endpoints := priceEndpoint.New(srv, logger, duration, otTracer, zipkin)
	hs := health.NewServer()
	hs.SetServingStatus(cfg.ServiceName, healthgrpc.HealthCheckResponse_SERVING)

	wg := &sync.WaitGroup{}
	go startHTTPServer(ctx, wg, endpoints, otTracer, zipkin, cfg.ServiceHost, cfg.HttpPort, logger)
	go startGRPCServer(ctx, wg, endpoints, otTracer, zipkin, cfg.ServiceHost, cfg.GrpcPort, hs, logger)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	cancel()
	wg.Wait()
}

func startHTTPServer(ctx context.Context, wg *sync.WaitGroup, endpoints priceEndpoint.Endpoints, otTracer stdopentracing.Tracer, zipkinTracer *zipkin.Tracer, host string, port string, logger log.Logger) {
	wg.Add(1)
	defer wg.Done()

	if port == "" {
		level.Error(logger).Log("protocol", "HTTP", "exposed", port, "err", "port is not assigned exist")
		return
	}

	p := fmt.Sprintf("%s:%s", host, port)
	// create a server
	srv := &http.Server{
		Addr:    p,
		Handler: priceHttpTransport.NewHTTPHandler(endpoints, otTracer, zipkinTracer, logger),
	}
	level.Info(logger).Log("protocol", "HTTP", "exposed", port)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			level.Info(logger).Log("Listen", err)
		}
	}()

	<-ctx.Done()

	// shut down gracefully, but wait no longer than 5 seconds before halting
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	srv.Shutdown(shutdownCtx)

	level.Info(logger).Log("protocol", "HTTP", "Shutdown", "http server gracefully stopped")
}

func startGRPCServer(ctx context.Context, wg *sync.WaitGroup, endpoints endpoints.Endpoints, tracer stdopentracing.Tracer, zipkinTracer *zipkin.Tracer, host string, port string, hs *health.Server, logger log.Logger) {
	wg.Add(1)
	defer wg.Done()

	p := fmt.Sprintf("%s:%s", host, port)
	listener, err := net.Listen("tcp", p)
	if err != nil {
		level.Error(logger).Log("protocol", "GRPC", "listen", port, "err", err)
		os.Exit(1)
	}

	var server *grpc.Server
	level.Info(logger).Log("protocol", "GRPC", "exposed", port)
	server = grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor))
	pb.RegisterPriceServer(server, priceGrpcTransport.MakeGRPCServer(endpoints, tracer, zipkinTracer, logger))
	healthgrpc.RegisterHealthServer(server, hs)
	reflection.Register(server)

	go func() {
		err = server.Serve(listener)
		if err != nil {
			fmt.Printf("grpc serve : %s\n", err)
		}
	}()

	<-ctx.Done()

	server.GracefulStop()
	fmt.Println("grpc server gracefully stopped")
}
