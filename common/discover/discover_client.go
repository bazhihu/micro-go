package discover

import "log"

type DiscoveryClient interface {
	/**
	服务注册接口
	@param serviceName 服务名
	@param instanceId 服务实例Id
	@param instancePort 服务实例端口
	@param healthCheckUrl 健康实例地址
	@param instanceHost 服务实例地址
	@param meta 服务实例元数据
	*/
	Register(serviceName, instanceId, healthCheckUrl, instanceHost string, instancePort int, meta map[string]string, logger *log.Logger) bool

	/**
	服务注销接口
	@param instanceId 服务实例Id
	*/
	DeRegister(instanceId string, logger *log.Logger) bool

	/**
	服务发现接口
	@param serviceName 服务名
	*/
	DiscoveryServices(serviceName string, logger *log.Logger) []interface{}
}
