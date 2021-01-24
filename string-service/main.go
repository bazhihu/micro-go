package main

import (
	"context"
	"flag"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"micro-go/common/discover"
	"micro-go/string-service/config"
	"micro-go/string-service/endpoint"
	"micro-go/string-service/plugins"
	"micro-go/string-service/service"
	"micro-go/string-service/transport"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	// 获取命令行参数
	var (
		servicePort = flag.Int("service.port", 10085, "service port")
		serviceHost = flag.String("service.host", "127.0.0.1", "service host")
		consulPort  = flag.Int("consul.port", 8500, "consul port")
		consulHost  = flag.String("consul.host", "127.0.0.1", "consul host")
		serviceName = flag.String("service.name", "string", "service name")
	)
	flag.Parse()

	ctx := context.Background()
	errChan := make(chan error)

	var discoveryClient discover.DiscoveryClient
	discoveryClient, err := discover.NewKitDiscoverClient(*consulHost, *consulPort)
	if err != nil {
		config.Logger.Println("Get Consul Client Failed")
		os.Exit(-1)
	}

	var svc service.Service
	svc = service.StringService{}
	svc = plugins.LoggingMiddleware(config.KitLogger)(svc)

	// 创建 endpoint
	stringEndpoint := endpoint.MakeStringEndpoint(svc)
	healthEndpoint := endpoint.MakeHealthCheckEndpoint(svc)

	// 把StringEndpoint 和healthCheckEndpoint 封装至StringEndpoints
	endpts := endpoint.StringEndpoints{
		StringEndpoints:      stringEndpoint,
		HealthCheckEndpoints: healthEndpoint,
	}

	// 创建http.Handler
	r := transport.MakeHttpHandler(ctx, endpts, config.KitLogger)

	instanceId := *serviceName + "-" + uuid.NewV4().String()

	// http server
	go func() {
		config.Logger.Println("HTTP server start at port:" + strconv.Itoa(*servicePort))
		// 启动前 执行注册服务发现
		if !discoveryClient.Register(*serviceName, instanceId, "/health", *serviceHost, *servicePort, nil, config.Logger) {
			config.Logger.Printf("string-service for service %s failed.", serviceName)
			os.Exit(-1)
		}
		handler := r
		errChan <- http.ListenAndServe(":"+strconv.Itoa(*servicePort), handler)
	}()

	// 监控系统关闭信号 ctr+c
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()
	error := <-errChan

	// 服务退出 取消注册服务发现
	discoveryClient.DeRegister(instanceId, config.Logger)
	config.Logger.Println(error)
}
