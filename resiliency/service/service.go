package service

import (
	"encoding/json"
	"errors"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/hashicorp/consul/api"
	"micro-go/common/discover"
	"micro-go/common/loadbalance"
	"micro-go/resiliency/config"
	"net/http"
	"net/url"
	"strconv"
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
	UseStringService(operationType, a, b string) (string, error)
	// 健康检查
	HealthCheck() bool
}

type UseStringService struct {
	// 服务发现客户端
	discoveryClient discover.DiscoveryClient
	// 负载均衡器
	loadbalance loadbalance.LoadBalance
}

type StringResponse struct {
	Result string `json:"result"`
	Error  error  `json:"error"`
}

func (u UseStringService) UseStringService(operationType, a, b string) (string, error) {
	var (
		operationResult string
		err             error
	)

	instances := u.discoveryClient.DiscoveryServices(StringService, config.Logger)
	instanceList := make([]*api.AgentService, len(instances))
	for i := 0; i < len(instances); i++ {
		instanceList[i] = instances[i].(*api.AgentService)
	}

	// 使用负载均衡算法选取实例
	selectInstance, err := u.loadbalance.SelectService(instanceList)
	if err != nil {
		return operationResult, err
	}

	config.Logger.Printf("current string-service ID is %s and address:port is %s:%d\n", selectInstance.ID, selectInstance.Address, selectInstance.Port)

	requestUrl := url.URL{
		Scheme: "http",
		Host:   selectInstance.Address + ":" + strconv.Itoa(selectInstance.Port),
		Path:   "/op/" + operationType + "/" + a + "/" + b,
	}

	resp, err := http.Post(requestUrl.String(), "", nil)
	if err == nil {
		// 解析结果
		result := &StringResponse{}
		err = json.NewDecoder(resp.Body).Decode(result)
		if err == nil && result.Error == nil {
			operationResult = result.Result
		}
	}

	return operationResult, err
}

func (u UseStringService) HealthCheck() bool {
	return true
}

func NewUseStringService(client discover.DiscoveryClient, lb loadbalance.LoadBalance) Service {
	hystrix.ConfigureCommand(StringServiceCommandName, hystrix.CommandConfig{
		RequestVolumeThreshold: 5,
	})
	return &UseStringService{
		discoveryClient: client,
		loadbalance:     lb,
	}
}

type ServiceMiddleware func(Service) Service
