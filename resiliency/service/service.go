package service

import (
	"micro-go/common/discover"
	"micro-go/common/loadbalance"
)

// Service constants
const (
	StringServiceCommandName = "String.string"
	StringService            = "string"
)

type UseStringService struct {
	// 服务发现客户端
	discoveryClient discover.DiscoveryClient
	// 负载均衡器
	loadbalance loadbalance.LoadBalance
}
