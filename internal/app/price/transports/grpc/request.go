package grpc

import (
	"context"

	"github.com/Hank-Kuo/go-kit-example/internal/app/price/endpoints"
	pb "github.com/Hank-Kuo/go-kit-example/pb/price"
)

func encodeGRPCSumRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(endpoints.SumRequest)
	return &pb.SumRequest{Price: req.Price, Fee: req.Fee}, nil
}

func decodeGRPCSumRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.SumRequest)
	return endpoints.SumRequest{Price: req.Price, Fee: req.Fee}, nil
}
