package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	kitzipkin "github.com/go-kit/kit/tracing/zipkin"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"google.golang.org/grpc"
	pd "micro-go/trace/zipkin-kit/pb"
	endpoints "micro-go/trace/zipkin-kit/string-service/endpoint"
	"micro-go/trace/zipkin-kit/string-service/service"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
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

	// string endpoint
	endpoint := endpoints.MakeStringEndpoint(ctx, svc)
	endpoint = kitzipkin.TraceEndpoint(zipkinTracer, "string-endpoint")(endpoint)

	// 健康检查的Endpoint
	healthEndpoint := endpoints.MakeHealthEndpoint(ctx, svc)
	healthEndpoint = kitzipkin.TraceEndpoint(zipkinTracer, "health-endpoint")(healthEndpoint)

	endpts := endpoints.StringEndpoints{
		StringEndpoint:      endpoint,
		HealthCheckEndpoint: healthEndpoint,
	}

	// 创建http.Handler
	r := MakeHttpHandler(ctx, endpts, zipkinTracer, logger)

	// 创建注册对象
	register := Register(*consulHost, *consulPost, *serviceHost, *servicePort, logger)

	go func() {
		fmt.Println("Http Server start at port:" + *servicePort)
		// 启动注册
		register.Register()
		handler := r
		errChan <- http.ListenAndServe(":"+*servicePort, handler)
	}()

	// grpc server
	go func() {
		fmt.Println("grpc Server start at port" + *grpcAddr)
		listener, err := net.Listen("tcp", *grpcAddr)
		if err != nil {
			errChan <- err
			return
		}

		serverTracer := kitzipkin.GRPCServerTrace(zipkinTracer, kitzipkin.Name("string-grpc-transport"))

		handler := NewGrpcServer(ctx, endpts, serverTracer)
		gRPCServer := grpc.NewServer()
		pd.RegisterStringServiceServer(gRPCServer, handler)
		errChan <- gRPCServer.Serve(listener)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	error := <-errChan
	// 服务退出取消注册
	register.Deregister()
	fmt.Println(error)
}

func Register(consulHost, consulPort, svcHost, svcPort string, logger log.Logger) (register sd.Registrar) {

	// 创建Consul 客户端连接
	var client consul.Client
	{
		consulCfg := api.DefaultConfig()
		consulCfg.Address = consulHost + ":" + consulPort
		consulClient, err := api.NewClient(consulCfg)
		if err != nil {
			logger.Log("create consul client error:", err)
			os.Exit(1)
		}

		client = consul.NewClient(consulClient)
	}

	// 设置Consul 对服务健康检查的参数
	check := api.AgentServiceCheck{
		HTTP:     "http://" + svcHost + ":" + svcPort + "/health",
		Interval: "10s",
		Timeout:  "1s",
		Notes:    "Consul check service health status.",
	}

	port, _ := strconv.Atoi(svcPort)

	// 设置微服务想Consul的注册信息
	reg := api.AgentServiceRegistration{
		ID:      "string-service" + uuid.New().String(),
		Name:    "string-service",
		Tags:    []string{"string-service"},
		Port:    port,
		Address: svcHost,
		Check:   &check,
	}

	// 创建注册对象
	register = consul.NewRegistrar(client, &reg, logger)
	return
}
