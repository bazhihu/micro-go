package endpoint

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"micro-go/rpc_demo/kit/service"
	"micro-go/rpc_demo/pd"
	"strings"
)

var (
	ErrInvalidRequestType = errors.New("RequestType has only two type: Concat, Diff")
)

type StringEndpoints struct {
	StringEndpoint      endpoint.Endpoint
	HealthCheckEndpoint endpoint.Endpoint
}

type StringRequest struct {
	RequestType string `json:"request_type"`
	A           string `json:"a"`
	B           string `json:"b"`
}

type StringResponse struct {
	Result string `json:"result"`
	Error  error  `json:"error"`
}

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
			res, _ = svc.Concat(ctx, a, b)
		} else if strings.EqualFold(req.RequestType, "Diff") {
			res, _ = svc.Diff(ctx, a, b)
		} else {
			return nil, ErrInvalidRequestType
		}
		return StringResponse{
			Result: res,
			Error:  opError,
		}, nil
	}
}

func DecodeStringRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*pd.StringRequest)

	return StringRequest{
		RequestType: "Concat",
		A:           req.A,
		B:           req.B,
	}, nil
}

func DecodeDiffStringRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*pd.StringRequest)

	return StringRequest{
		RequestType: "Diff",
		A:           req.A,
		B:           req.B,
	}, nil
}

func EncodeStringResponse(ctx context.Context, r interface{}) (interface{}, error) {
	resp := r.(StringResponse)
	if resp.Error != nil {
		return &pd.StringResponse{
			Ret: resp.Result,
			Err: resp.Error.Error(),
		}, nil
	}
	return &pd.StringResponse{
		Ret: resp.Result,
		Err: "",
	}, nil
}
