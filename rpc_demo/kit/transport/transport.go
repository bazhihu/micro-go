package transport

import (
	"context"
	"errors"
	"github.com/go-kit/kit/transport/grpc"
	endpoint2 "micro-go/rpc_demo/kit/endpoint"
	"micro-go/rpc_demo/pd"
)

var (
	ErrBadRequest = errors.New("invalid request parameter")
)

type grpcServer struct {
	concat grpc.Handler
	diff   grpc.Handler
}

func (s *grpcServer) Concat(ctx context.Context, r *pd.StringRequest) (*pd.StringResponse, error) {
	_, resp, err := s.concat.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pd.StringResponse), nil
}

func (s *grpcServer) Diff(ctx context.Context, r *pd.StringRequest) (*pd.StringResponse, error) {
	_, resp, err := s.diff.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pd.StringResponse), nil
}

func NewStringServer(ctx context.Context, endpoints endpoint2.StringEndpoints) pd.StringServiceServer {
	return &grpcServer{
		concat: grpc.NewServer(endpoints.StringEndpoint, DecodeStringRequest, EncodeStringResponse),
		diff:   grpc.NewServer(endpoints.StringEndpoint, DecodeDiffStringRequest, EncodeStringResponse),
	}
}

func DecodeStringRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*pd.StringRequest)

	return endpoint2.StringRequest{
		RequestType: "Concat",
		A:           req.A,
		B:           req.B,
	}, nil
}

func DecodeDiffStringRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*pd.StringRequest)

	return endpoint2.StringRequest{
		RequestType: "Diff",
		A:           req.A,
		B:           req.B,
	}, nil
}

func EncodeStringResponse(ctx context.Context, r interface{}) (interface{}, error) {
	resp := r.(endpoint2.StringResponse)
	if resp.Error != nil {
		return &pd.StringResponse{
			Ret: resp.Result,
			Err: resp.Error.Error(),
		}, nil
	}
	return &pd.StringResponse{
		Ret: resp.Result,
		Err: "",
	}, nil
}
