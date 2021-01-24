package discover

import (
	"github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	"log"
	"strconv"
	"sync"
)

type KitDiscoverClient struct {
	Host   string // Consul Host
	Port   int    // Consul Port
	client consul.Client

	// 连接 consul 的配置
	config *api.Config
	mutex  sync.Mutex
	// 服务实例缓存字段
	instancesMap sync.Map
}

func (consulClient *KitDiscoverClient) Register(serviceName, instanceId, healthCheckUrl, instanceHost string, instancePort int, meta map[string]string, logger *log.Logger) bool {
	// 构建服务实例元数据
	serviceRegistration := &api.AgentServiceRegistration{
		ID:      instanceId,
		Name:    serviceName,
		Port:    instancePort,
		Address: instanceHost,
		Meta:    meta,
		Check: &api.AgentServiceCheck{
			Interval:                       "15s",
			HTTP:                           "http://" + instanceHost + ":" + strconv.Itoa(instancePort) + healthCheckUrl,
			DeregisterCriticalServiceAfter: "30s",
		},
	}

	// 发送服务注册到Consul中
	err := consulClient.client.Register(serviceRegistration)
	if err != nil {
		log.Println("Register service error!")
		return false
	}
	log.Println("Register service success!")
	return true
}

func (consulClient *KitDiscoverClient) DeRegister(instanceId string, logger *log.Logger) bool {
	// 构建包含服务实例ID的元数据结构体
	serviceRegistration := &api.AgentServiceRegistration{
		ID: instanceId,
	}
	// 发送服务注销请求
	err := consulClient.client.Deregister(serviceRegistration)

	if err != nil {
		logger.Println("deregister service Error!")
		return false
	}

	logger.Println("Deregister Service success!")
	return true
}

func (consulClient *KitDiscoverClient) DiscoveryServices(serviceName string, logger *log.Logger) []interface{} {
	// 该服务已监控并缓存
	instanceList, ok := consulClient.instancesMap.Load(serviceName)
	if ok {
		return instanceList.([]interface{})
	}

	// 单例模型枷锁
	consulClient.mutex.Lock()
	// 再次判断是否监控
	instanceList, ok = consulClient.instancesMap.Load(serviceName)
	if ok {
		return instanceList.([]interface{})
	}
	// 注册监控
	go func() {
		// 使用consul 服务实例监控来监控某个服务名的服务实例列表变化
		params := make(map[string]interface{})
		params["type"] = "service"
		params["service"] = serviceName
		plan, _ := watch.Parse(params)
		plan.Handler = func(u uint64, i interface{}) {
			if i == nil {
				return
			}
			v, ok := i.([]*api.ServiceEntry)
			if !ok {
				return // 数据异常，忽略
			}
			// 没有服务实例在线
			if len(v) == 0 {
				consulClient.instancesMap.Store(serviceName, []interface{}{})
			} else {
				var healthServices []interface{}
				for _, service := range v {
					if service.Checks.AggregatedStatus() == api.HealthPassing {
						healthServices = append(healthServices, service.Service)
					}
				}
				consulClient.instancesMap.Store(serviceName, healthServices)
			}
		}
		defer plan.Stop()
		plan.Run(consulClient.config.Address)
	}()
	defer consulClient.mutex.Unlock()

	// 根据服务名请求服务实例列表
	entries, _, err := consulClient.client.Service(serviceName, "", false, nil)
	if err != nil {
		consulClient.instancesMap.Store(serviceName, []interface{}{})
		logger.Println("Discover Service Error!")
		return nil
	}

	instances := make([]interface{}, len(entries))
	for i := 0; i < len(entries); i++ {
		instances[i] = entries[i].Service
	}
	consulClient.instancesMap.Store(serviceName, instances)
	return instances
}

// 创建kit的服务注册客户端
func NewKitDiscoverClient(consulHost string, consulPort int) (DiscoveryClient, error) {
	// 通过Consul Host和Consul Port创建一个consul.Client
	consulConfig := api.DefaultConfig()

	consulConfig.Address = consulHost + ":" + strconv.Itoa(consulPort)
	apiClient, err := api.NewClient(consulConfig)
	if err != nil {
		return nil, err
	}

	client := consul.NewClient(apiClient)

	return &KitDiscoverClient{
		Host:   consulHost,
		Port:   consulPort,
		config: consulConfig,
		client: client,
	}, err
}
