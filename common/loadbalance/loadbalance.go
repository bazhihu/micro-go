package loadbalance

import (
	"errors"
	"github.com/hashicorp/consul/api"
	"math/rand"
)

// 负载均衡器
type LoadBalance interface {
	SelectService(service []*api.AgentService) (*api.AgentService, error)
}

// 随机
type RandomLoadBalance struct {
}

func (LoadBalance *RandomLoadBalance) SelectService(services []*api.AgentService) (*api.AgentService, error) {
	if services == nil || len(services) < 1 {
		return nil, errors.New("service instances are not exist")
	}
	return services[rand.Intn(len(services))], nil
}

// 轮询

// 轮询加权法

// 一致性hash 法

// 最小连接数法
