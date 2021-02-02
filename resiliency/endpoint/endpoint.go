package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"micro-go/resiliency/service"
)

// define endpoint
type UseStringEndpoints struct {
	UseStringEndpoint   endpoint.Endpoint
	HealthCheckEndpoint endpoint.Endpoint
}

// define request struct
type UseStringRequest struct {
	RequestType string `json:"request_type"`
	A           string `json:"a"`
	B           string `json:"b"`
}

// define response struct
type UseStringResponse struct {
	Result string `json:"result"`
	Error  string `json:"error"`
}

func MakeUseStringEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UseStringRequest)
		var (
			res, a, b string
			opError   error
		)
		a = req.A
		b = req.B

		res, opError = svc.UseStringService(req.RequestType, a, b)
		return UseStringResponse{Result: res}, opError
	}
}

// 健康检查请求结构
type HealthRequest struct{}

// 健康检查响应结构
type HealthResponse struct {
	Status bool `json:"status"`
}

// 创建健康检查Endpoint
func MakeHealthCheckEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		status := svc.HealthCheck()
		return HealthResponse{status}, nil
	}
}
