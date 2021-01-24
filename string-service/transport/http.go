package transport

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"micro-go/string-service/endpoint"
	"net/http"
)

var (
	ErrorBadRequest = errors.New("invalid request parameter")
)

// make handler use mux
func MakeHttpHandler(ctx context.Context, endpoints endpoint.StringEndpoints, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	options := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
	}

	// 字符串操作接口
	r.Methods("POST").Path("/op/{type}/{a}/{b}").Handler(kithttp.NewServer(
		endpoints.StringEndpoints, decodeStringRequest, encodeStringResponse, options...))

	r.Path("/metrics").Handler(promhttp.Handler())

	// create health check handler
	r.Methods("GET").Path("/health").Handler(kithttp.NewServer(
		endpoints.HealthCheckEndpoints, decodeHealthCheckRequest, encodeStringResponse, options...))

	return r
}

func decodeHealthCheckRequest(i context.Context, request2 *http.Request) (request interface{}, err error) {
	return endpoint.HealthRequest{}, nil
}

func encodeStringResponse(i context.Context, writer http.ResponseWriter, response interface{}) error {
	writer.Header().Set("Content-Type", "application/json;charset=utf-8")
	return json.NewEncoder(writer).Encode(response)
}

func decodeStringRequest(i context.Context, request2 *http.Request) (request interface{}, err error) {
	vars := mux.Vars(request2)
	requestType, ok := vars["type"]
	if !ok {
		return nil, ErrorBadRequest
	}
	pa, ok := vars["a"]
	if !ok {
		return nil, ErrorBadRequest
	}
	pb, ok := vars["b"]
	if !ok {
		return nil, ErrorBadRequest
	}
	return endpoint.StringRequest{
		RequestType: requestType,
		A:           pa,
		B:           pb,
	}, nil
}

// 自定义错误响应
func encodeError(ctx context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
}
