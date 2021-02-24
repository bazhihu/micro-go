package client

import (
	"context"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
	"micro-go/trace/zipkin-kit/pb"
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

}

// 解密GRPC请求
func DecodeGRPCStringResponse(ctx context.Context, i interface{}) (response interface{}, err error) {

}

// 加密GRPC请求
func EncodeGRPCStringRequest(ctx context.Context, i interface{}) (request interface{}, err error) {
	resp := r.(*pd.StringResponse)
	return endpoint.StringResponse{
		Result: "",
		Error:  nil,
	}
}
