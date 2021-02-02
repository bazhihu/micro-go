package service

import (
	"errors"
	"github.com/afex/hystrix-go/hystrix"
	"micro-go/common/discover"
	"micro-go/common/loadbalance"
)

// Service constants
const (
	StringServiceCommandName = "String.string"
	StringService            = "string"
)

var (
	ErrHystrixFallbackExecute = errors.New("hystrix fall back execute")
)

type Service interface {
	// 远程调用 string-service 服务
}

type UseStringService struct {
	// 服务发现客户端
	discoveryClient discover.DiscoveryClient
	// 负载均衡器
	loadbalance loadbalance.LoadBalance
}

func NewUseStringService(client discover.DiscoveryClient, lb loadbalance.LoadBalance) UseStringService {
	hystrix.ConfigureCommand(StringServiceCommandName, hystrix.CommandConfig{
		RequestVolumeThreshold: 5,
	})
	return &UseStringService{
		discoveryClient: client,
		loadbalance:     lb,
	}
}
