package service

/**
业务代码实现层
*/

import (
	"context"
	"errors"
	discover2 "micro-go/register/discover"
)

type Service interface {
	// 健康检查接口
	HealthCheck() bool
	// 打招呼接口
	SaveHello() string
	// 服务发现接口
	DiscoveryService(ctx context.Context, serviceName string) ([]interface{}, error)
}

var ErrNotServiceInstances = errors.New("instances are not existed")

type DiscoveryServiceImpl struct {
	discoveryClient discover2.DiscoveryClient
}

func NewDiscoveryServiceImpl(discoveryClient discover2.DiscoveryClient) Service {
	return &DiscoveryServiceImpl{discoveryClient: discoveryClient}
}

// 健康服务的健康状态
func (service *DiscoveryServiceImpl) HealthCheck() bool {
	return true
}

// 打招呼
func (service *DiscoveryServiceImpl) SaveHello() string {
	return "Hello World"
}

// 服务发现
func (service *DiscoveryServiceImpl) DiscoveryService(ctx context.Context, serviceName string) ([]interface{}, error) {
	// 从consul中根据服务名获取服务实例列表

	//instances := service.discoveryClient.DiscoveryServices(serviceName, config.logg)

	return nil, nil
}
