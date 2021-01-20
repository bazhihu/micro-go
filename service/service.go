package service

import (
	"context"
)

type Service interface {
	// 健康检查接口
	HealthCheck() bool
	// 打招呼接口
	SaveHello() string
	// 服务发现接口
	DiscoveryService(ctx context.Context, serviceName string) ([]interface{}, error)
}

type DiscoveryServiceImpl struct {
}

func (service *DiscoveryServiceImpl) DiscoveryService(ctx context.Context, serviceName string) ([]interface{}, error) {
	// 从consul中根据服务名获取服务实例列表

	return nil, nil
}
