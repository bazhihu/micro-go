package discover

type InstanceInfo struct {
	ID                string            `json:"id"`                  // 服务实例ID
	Name              string            `json:"name"`                // 服务名
	Service           string            `json:"service,omitempty"`   // 服务发现时返回的服务名
	Tags              []string          `json:"tags,omitempty"`      // 标签,可用于进行服务过滤
	Address           string            `json:"address"`             // 服务实例HOST
	Port              int               `json:"port"`                // 服务实例端口
	Meta              map[string]string `json:"meta,omitempty"`      // 元数据
	EnableTagOverride bool              `json:"enable_tag_override"` // 是否允许标签覆盖
}

type Check struct {
	DeregisterCriticalServiceAfter string   `json:"deregister_critical_service_after"` // 多久之前注销服务
	Args                           []string `json:"args"`                              // 请求参数
	HTTP                           string   `json:"http"`                              // 健康检查地址
	Interval                       string   `json:"interval"`                          // Consul主动进行健康检查
	TTL                            string   `json:"ttl"`                               // 服务实例主动提交健康检查
}
