package endpoint

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"micro-go/string-service/service"
	"strings"
)

// define endpoint
type StringEndpoints struct {
	StringEndpoints      endpoint.Endpoint
	HealthCheckEndpoints endpoint.Endpoint
}

var (
	ErrInvalidRequestType = errors.New("RequestType has only two type: Concat, Diff")
)

// define request struct
type StringRequest struct {
	RequestType string `json:"request_type"`
	A           string `json:"a"`
	B           string `json:"b"`
}

type StringResponse struct {
	Result string `json:"result"`
	Error  error  `json:"error"`
}

// make endpoint
func MakeStringEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(StringRequest)
		var (
			res, a, b string
			opError   error
		)
		a = req.A
		b = req.B
		if strings.EqualFold(req.RequestType, "Concat") {
			res, _ = svc.Concat(a, b)
		} else if strings.EqualFold(req.RequestType, "Diff") {
			res, _ = svc.Diff(a, b)
		} else {
			return nil, ErrInvalidRequestType
		}

		return StringResponse{
			Result: res,
			Error:  opError,
		}, nil
	}
}

// 健康检查请求结构
type HealthRequest struct {
}

type HealthResponse struct {
	Status bool `json:"status"`
}

// 创建健康检查Endpoint
func MakeHealthCheckEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		status := svc.HealthCheck()
		return HealthResponse{Status: status}, nil
	}
}
