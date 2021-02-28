package main

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/transport/grpc"
	"micro-go/trace/zipkin-kit/client"
	pd "micro-go/trace/zipkin-kit/pb"
	"micro-go/trace/zipkin-kit/string-service/endpoint"
)

type grpcServer struct {
	diff grpc.Handler
}

func (s *grpcServer) Diff(ctx context.Context, r *pd.StringRequest) (*pd.StringResponse, error) {
	fmt.Println(r)
	_, resp, err := s.diff.ServeGRPC(ctx, r)

	if err != nil {
		return nil, err
	}
	return resp.(*pd.StringResponse), nil
}

func NewGrpcServer(ctx context.Context, endpoints endpoint.StringEndpoints, option grpc.ServerOption) pd.StringServiceServer {
	return &grpcServer{diff: grpc.NewServer(endpoints.StringEndpoint, client.DecodeGRPCStringResponse, client.EncodeGRPCStringRequest, option)}
}
