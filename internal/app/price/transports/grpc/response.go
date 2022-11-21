package grpc

import (
	"context"

	"github.com/Hank-Kuo/go-kit-example/internal/app/price/endpoints"
	pb "github.com/Hank-Kuo/go-kit-example/pb/price"
	"github.com/Hank-Kuo/go-kit-example/pkg/errors"
)

func decodeGRPCSumResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.SumResponse)
	d := endpoints.SumResponse{Cost: reply.Cost}
	return d, nil
}

func encodeGRPCSumResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(endpoints.SumResponse)
	d := &pb.SumResponse{Cost: reply.Cost}

	// details, _ := ptypes.MarshalAny(d)
	// &stardandPB.StandardResponse{
	// 	Status:  "success",
	// 	Message: "Sum of two number",
	// 	Data:    details,
	// }
	return d, grpcEncodeError(errors.Cast(reply.Err))
}
