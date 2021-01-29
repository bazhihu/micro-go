package main

import (
	"context"
	"flag"
	"github.com/go-kit/kit/log"
	"google.golang.org/grpc"
	endpoint2 "micro-go/rpc_demo/kit/endpoint"
	"micro-go/rpc_demo/kit/service"
	"micro-go/rpc_demo/kit/transport"
	"micro-go/rpc_demo/pd"
	"net"
	"os"
)

func main() {
	flag.Parse()

	ctx := context.Background()

	// 日志
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	// 初始化service
	var svc service.Service
	svc = service.StringService{}

	// add logging middleware
	svc = service.LoggingMiddleware(logger)(svc)

	endpoint := endpoint2.MakeStringEndpoint(svc)

	// 创建健康检查的Endpoint
	healthEndpoint := endpoint2.MakeHealthCheckEndpoint(svc)

	endpts := endpoint2.StringEndpoints{
		StringEndpoint:      endpoint,
		HealthCheckEndpoint: healthEndpoint,
	}

	handler := transport.NewStringServer(ctx, endpts)

	ls, _ := net.Listen("tcp", "127.0.0.1:8080")
	gRPCServer := grpc.NewServer()
	pd.RegisterStringServiceServer(gRPCServer, handler)
	gRPCServer.Serve(ls)
}
