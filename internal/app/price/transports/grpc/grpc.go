package grpc

import (
	"context"

	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/tracing/zipkin"
	grpctransport "github.com/go-kit/kit/transport/grpc"

	stdopentracing "github.com/opentracing/opentracing-go"
	stdzipkin "github.com/openzipkin/zipkin-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Hank-Kuo/go-kit-example/internal/app/price"
	"github.com/Hank-Kuo/go-kit-example/internal/app/price/endpoints"
	pb "github.com/Hank-Kuo/go-kit-example/pb/price"
	"github.com/Hank-Kuo/go-kit-example/pkg/errors"
)

type grpcServer struct {
	sum grpctransport.Handler
}

func (s *grpcServer) Sum(ctx context.Context, req *pb.SumRequest) (rep *pb.SumResponse, err error) {
	_, rp, err := s.sum.ServeGRPC(ctx, req)
	if err != nil {
		return nil, grpcEncodeError(errors.Cast(err))
	}
	rep = rp.(*pb.SumResponse)
	return rep, nil
}

func MakeGRPCServer(endpoints endpoints.Endpoints, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer, logger log.Logger) (req pb.PriceServer) {
	options := []grpctransport.ServerOption{
		// grpctransport.ServerBefore(telepresence.GRPCToContext()),
		grpctransport.ServerErrorLogger(logger),
	}

	if zipkinTracer != nil {
		options = append(options, zipkin.GRPCServerTrace(zipkinTracer))
	}

	return &grpcServer{
		sum: grpctransport.NewServer(
			endpoints.SumEndpoint,
			decodeGRPCSumRequest,
			encodeGRPCSumResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(otTracer, "Sum", logger)))...,
		),
	}
}

func NewGRPCClient(conn *grpc.ClientConn, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer, logger log.Logger) price.Service {
	// global client middlewares
	options := []grpctransport.ClientOption{
		// grpctransport.ClientBefore(telepresence.ContextToGRPC()),
	}

	if zipkinTracer != nil {
		options = append(options, zipkin.GRPCClientTrace(zipkinTracer))
	}

	var sumEndpoint endpoint.Endpoint
	{
		sumEndpoint = grpctransport.NewClient(
			conn,
			"pb.Add",
			"Sum",
			encodeGRPCSumRequest,
			decodeGRPCSumResponse,
			pb.SumResponse{},
			append(options, grpctransport.ClientBefore(opentracing.ContextToGRPC(otTracer, logger), jwt.ContextToGRPC()))...,
		).Endpoint()
		sumEndpoint = opentracing.TraceClient(otTracer, "Sum")(sumEndpoint)
	}

	return endpoints.Endpoints{
		SumEndpoint: sumEndpoint,
	}
}

func grpcEncodeError(err errors.Error) error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if ok {
		return status.Error(st.Code(), st.Message())
	}

	// switch {

	// case errors.Contains(err, jwt.ErrXJWTContextMissing):
	// 	return status.Error(codes.Unauthenticated, err.Error())
	// default:

	// }
	return status.Error(codes.Internal, "internal server error")
}
