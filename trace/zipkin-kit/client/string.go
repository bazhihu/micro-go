package client

import (
	"context"
	"errors"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
	"micro-go/trace/zipkin-kit/pb"
	endpts "micro-go/trace/zipkin-kit/string-service/endpoint"
	"micro-go/trace/zipkin-kit/string-service/service"
)

func StringDiff(conn *grpc.ClientConn, clientTracer kitgrpc.ClientOption) service.Service {
	var ep = kitgrpc.NewClient(conn,
		"pb.StringService",
		"Diff",
		EncodeGRPCStringRequest,
		DecodeGRPCStringResponse,
		pd.StringResponse{},
		clientTracer).Endpoint()

	StringEp := endpts.StringEndpoints{
		StringEndpoint: ep,
	}
	return StringEp
}

// 加密GRPC请求
func DecodeGRPCStringResponse(ctx context.Context, i interface{}) (response interface{}, err error) {
	resp := i.(*pd.StringResponse)
	return endpts.StringResponse{
		Result: resp.Result,
		Error:  errors.New(resp.Err),
	}, nil
}

// 解密密GRPC请求
func EncodeGRPCStringRequest(ctx context.Context, i interface{}) (request interface{}, err error) {
	req := i.(*endpts.StringRequest)
	return &pd.StringRequest{
		RequestType: string(req.RequestType),
		A:           string(req.A),
		B:           string(req.B),
	}, nil
}
