package main

import (
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/hashicorp/consul/api"
	proxy2 "micro-go/gateway/proxy"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 创建环境变量
	var (
		consulHost = flag.String("consul.host", "127.0.0.1", "consul server ip address")
		consulPort = flag.String("consul.port", "8500", "consul server port")
	)
	flag.Parse()

	// 创建日志组件
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	// 创建consul api客户端
	consulConfig := api.DefaultConfig()
	consulConfig.Address = "http://" + *consulHost + ":" + *consulPort
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		logger.Log("err", err)
		os.Exit(1)
	}

	// 创建反向代理
	proxy := proxy2.NewReverseProxy(consulClient, logger)

	errc := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	// 开始监听
	go func() {
		logger.Log("transport", "HTTP", "addr", "9090")
		errc <- http.ListenAndServe(":9090", proxy)
	}()

	// 运行，等待结束
	logger.Log("exit", <-errc)
}
