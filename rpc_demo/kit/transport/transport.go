package transport

import (
	"context"
	"github.com/go-kit/kit/transport/grpc"
	"micro-go/rpc_demo/pd"
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
