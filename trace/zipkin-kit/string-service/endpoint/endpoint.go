package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"micro-go/trace/zipkin-kit/string-service/service"
)

// StringEndpoint define endpoint
type StringEndpoints struct {
	StringEndpoint      endpoint.Endpoint
	HealthCheckEndpoint endpoint.Endpoint
}

func (s StringEndpoints) Concat(a, b string) (string, error) {
	ctx := context.Background()
	resp, err := s.StringEndpoint(ctx, StringRequest{
		RequestType: "Concat",
		A:           a,
		B:           b,
	})
	response := resp.(StringResponse)
	return response.Result, err
}

func (s StringEndpoints) Diff(ctx context.Context, a, b string) (string, error) {
	resp, err := s.StringEndpoint(ctx, StringRequest{
		RequestType: "Diff",
		A:           a,
		B:           b,
	})
	response := resp.(StringResponse)
	return response.Result, err
}

func (s StringEndpoints) HealthCheck() bool {
	return true
}

// StringRequest define request struct
type StringRequest struct {
	RequestType string `json:"request_type"`
	A           string `json:"a"`
	B           string `json:"b"`
}

// make string endpoint
func MakeStringEndpoint(ctx context.Context, svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

	}
}

// StringResponse define response struct
type StringResponse struct {
	Result string `json:"result"`
	Error  error  `json:"error"`
}

// HealthRequest 健康检查请求结构
type HealthRequest struct{}

// HealthResponse 健康检查响应结构
type HealthResponse struct {
	Status bool `json:"status"`
}

func MakeHealthEndpoint(ctx context.Context, svc service.Service) endpoint.Endpoint {

}
