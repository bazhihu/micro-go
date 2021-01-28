package main

import (
	"context"
	"flag"
	"github.com/go-kit/kit/log"
	"micro-go/rpc_demo/kit/service"
	"micro-go/string-service/plugins"
	"os"
)

func main() {
	flag.Parse()

	ctx := context.Background()
	errChan := make(chan error)

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
	svc = service.ServiceMiddleware()

}
