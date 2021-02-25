package main

import (
	"context"
	"flag"
	"github.com/go-kit/kit/log"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"micro-go/trace/zipkin-kit/string-service/service"
	"os"
)

func main() {
	var (
		consulHost = flag.String("consul.host", "127.0.0.1", "consul ip address")
		consulPost = flag.String("consul.port", "8500", "consul port")

		serviceHost = flag.String("service.host", "127.0.0.1", "service ip address")
		servicePort = flag.String("service.port", "9009", "service port")

		zipkinURL = flag.String("zipkin.url", "http://127.0.0.1:9411/api/v2/spans", "zipkinURL")
		grpcAddr  = flag.String("grpc", ":9008", "gRPC listen address")
	)
	flag.Parse()

	ctx := context.Background()
	errChan := make(chan error)

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var zipkinTracer *zipkin.Tracer
	{
		var (
			err           error
			hostPort      = *serviceHost + ":" + *servicePort
			serviceName   = "string-service"
			useNoopTracer = (*zipkinURL == "")
			reporter      = zipkinhttp.NewReporter(*zipkinURL)
		)
		defer reporter.Close()

		zEP, _ := zipkin.NewEndpoint(serviceName, hostPort)
		zipkinTracer, err = zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(zEP), zipkin.WithNoopTracer(useNoopTracer))

		if err != nil {
			logger.Log("err", err)
			os.Exit(1)
		}
		if !useNoopTracer {
			logger.Log("tracer", "Zipkin", "type", "Native", "URL", *zipkinURL)
		}
	}

	var svc service.Service
	svc = service.StringService{}

	// add logging middleware to service
	svc = LoggingMiddleware(logger)(svc)

	endpoint := endpoint.MakeString

}
