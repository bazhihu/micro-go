package endpoint

import (
	"github.com/go-kit/kit/endpoint"
	"micro-go/trace/zipkin-kit/string-service/service"
)

// StringEndpoint define endpoint
type StringEndpoints struct {
	StringEndpoint      endpoint.Endpoint
	HealthCheckEndpoint endpoint.Endpoint
}

func (s StringEndpoints) Concat(req service.StringRequest, ret *string) error {
	panic("implement me")
}

func (s StringEndpoints) Diff(req service.StringRequest, ret *string) error {
	panic("implement me")
}

func (s StringEndpoints) HealthCheck() bool {
	panic("implement me")
}

// StringRequest define request struct
type StringRequest struct {
	RequestType string `json:"request_type"`
	A           string `json:"a"`
	B           string `json:"b"`
}

// StringResponse define response struct
type StringResponse struct {
	Result string `json:"result"`
	Error  error  `json:"error"`
}
