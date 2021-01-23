package main

import (
	"context"
	"flag"
	"fmt"
	"micro-go/config"
	"micro-go/discover"
	"micro-go/endpoint"
	"micro-go/service"
	"micro-go/transport"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	uuid "github.com/satori/go.uuid"
)

func main() {

	// 从命令行中读取相关参数，没有时使用默认值
	var (
		// 服务地址和服务名
		servicePort = flag.Int("service.port", 10086, "service port")
		serviceHost = flag.String("service.host", "127.0.0.1", "service host")
		serviceName = flag.String("service.name", "SayHello", "service name")

		// consul 地址
		consulPort = flag.Int("consul.port", 8500, "consul port")
		consulHost = flag.String("consul.host", "127.0.0.1", "consul host")
	)

	flag.Parse()

	ctx := context.Background()
	errChan := make(chan error)

	// 声明服务发现客户端
	var discoveryClient discover.DiscoveryClient
	discoveryClient = discover.NewMyDiscoverClient(*consulHost, *consulPort)

	//if err != nil {
	//	config.Logger.Println("Get Consul Client failed")
	//	os.Exit(-1)
	//}

	// 声明并初始化 Service
	var svc = service.NewDiscoveryServiceImpl(discoveryClient)
	// 创建打招呼的Endpoint
	sayHelloEndpoint := endpoint.MakeSayHelloEndpoint(svc)
	// 创建服务发现的Endpoint
	discoverEndpoint := endpoint.MakeDiscoveryEndpoint(svc)
	// 创建健康检查的Endpoint
	healthEndpoint := endpoint.MakeHealthCheckEndpoint(svc)

	endpts := endpoint.DiscoveryEndpoints{
		SayHelloEndpoint:    sayHelloEndpoint,
		DiscoveryEndpoint:   discoverEndpoint,
		HealthCheckEndpoint: healthEndpoint,
	}

	// 创建http.Handler
	r := transport.MakeHttpHandler(ctx, endpts, config.KitLogger)

	// 定义服务实例ID
	instanceId := *serviceName + "-" + uuid.NewV4().String()
	// 启动 http server
	go func() {
		config.Logger.Println("Http Server start at port:" + strconv.Itoa(*servicePort))
		// 启动前执行注册
		if !discoveryClient.Register(*serviceName, instanceId, "/health", *serviceHost, *servicePort, nil, config.Logger) {
			config.Logger.Printf("string-service for service %s failed.", *serviceName)
			// 注册失败, 服务启动失败
			os.Exit(-1)
		}
		handler := r
		errChan <- http.ListenAndServe(":"+strconv.Itoa(*servicePort), handler)
	}()

	go func() {
		// 监控系统信号, 等待Ctrl + c 系统信号通知服务关闭
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	error := <-errChan
	// 服务退出取消注册
	discoveryClient.DeRegister(instanceId, config.Logger)
	config.Logger.Println(error)
}
