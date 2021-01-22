package discover

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type InstanceInfo struct {
	ID                string            `json:"id"`                  // 服务实例ID
	Name              string            `json:"name"`                // 服务名
	Service           string            `json:"service,omitempty"`   // 服务发现时返回的服务名
	Tags              []string          `json:"tags,omitempty"`      // 标签,可用于进行服务过滤
	Address           string            `json:"address"`             // 服务实例HOST
	Port              int               `json:"port"`                // 服务实例端口
	Meta              map[string]string `json:"meta,omitempty"`      // 元数据
	EnableTagOverride bool              `json:"enable_tag_override"` // 是否允许标签覆盖
	Check             Check             `json:"check,omitempty"`     // 健康检查相关配置
	Weights           Weights           `json:"weights,omitempty"`   // 权重相关
}

type Check struct {
	DeregisterCriticalServiceAfter string   `json:"deregister_critical_service_after"` // 多久之前注销服务
	Args                           []string `json:"args"`                              // 请求参数
	HTTP                           string   `json:"http"`                              // 健康检查地址
	Interval                       string   `json:"interval"`                          // Consul主动进行健康检查
	TTL                            string   `json:"ttl"`                               // 服务实例主动提交健康检查
}

type Weights struct {
	Passing int `json:"passing"`
	Warning int `json:"warning"`
}

type MyDiscoverClient struct {
	Host string // Consul的端口
	Port int    // consul的端口
}

func (consulClient *MyDiscoverClient) Register(serviceName, instanceId, healthCheckUrl, instanceHost string, instancePort int, meta map[string]string, logger *log.Logger) bool {
	// 封装服务实例的元数据
	instanceInfo := &InstanceInfo{
		ID:                instanceId,
		Name:              serviceName,
		Address:           instanceHost,
		Port:              instancePort,
		Meta:              meta,
		EnableTagOverride: false,
		Check: Check{
			DeregisterCriticalServiceAfter: "30s",
			HTTP:                           "http://" + instanceHost + ":" + strconv.Itoa(instancePort) + healthCheckUrl,
			Interval:                       "15s",
		},
		Weights: Weights{
			Passing: 10,
			Warning: 1,
		},
	}
	byteData, _ := json.Marshal(instanceInfo)
	// 向Consul发送服务注册的请求
	req, err := http.NewRequest("PUT", "http://"+consulClient.Host+":"+strconv.Itoa(consulClient.Port)+"/v1/agent/service/register", bytes.NewReader(byteData))
	if err != nil {
		return false
	}
	// 检查注册结果

	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Println("Register service error")
		return false
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Println("Register service error")
		return false
	}

	log.Println("Register service success!")

	return true
}

func (consulClient *MyDiscoverClient) DeRegister(instanceId string, logger *log.Logger) bool {
	panic("implement me")
}

func (consulClient *MyDiscoverClient) DiscoveryServices(serviceName string, logger *log.Logger) []interface{} {
	panic("implement me")
}

func NewMyDiscoverClient(consulHost string, consulPort int) DiscoveryClient {
	return &MyDiscoverClient{
		Host: consulHost,
		Port: consulPort,
	}
}
