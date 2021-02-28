package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/zipkin"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	gozipkin "github.com/openzipkin/zipkin-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"micro-go/trace/zipkin-kit/string-service/endpoint"
	"net/http"
)

var (
	ErrorBadRequest = errors.New("invalid request parameter")
)

func MakeHttpHandler(_ context.Context, endpoints endpoint.StringEndpoints, zipkinTracer *gozipkin.Tracer, logger log.Logger) http.Handler {
	r := mux.NewRouter()

	zipkinServer := zipkin.HTTPServerTrace(zipkinTracer, zipkin.Name("http-transport"))

	options := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(kithttp.DefaultErrorEncoder),
		zipkinServer,
	}

	r.Methods("POST").Path("/op/{type}/{a}/{b}").Handler(kithttp.NewServer(endpoints.StringEndpoint, decodeStringRequest, encodeStringResponse, options...))

	// prometheus metrics
	r.Path("/metrics").Handler(promhttp.Handler())

	r.Methods("GET").Path("/health").Handler(kithttp.NewServer(endpoints.HealthCheckEndpoint, decodeHealthCheckEndpoint, encodeStringResponse, options...))

	return r
}

func decodeHealthCheckEndpoint(_ context.Context, _ *http.Request) (request interface{}, err error) {
	return endpoint.HealthRequest{}, nil
}

func encodeStringResponse(_ context.Context, writer http.ResponseWriter, i interface{}) error {
	writer.Header().Set("Content-Type", "application/json;charset=utf-8")
	return json.NewEncoder(writer).Encode(i)
}

func decodeStringRequest(_ context.Context, request2 *http.Request) (request interface{}, err error) {
	vars := mux.Vars(request2)
	requestType, ok := vars["type"]
	if !ok {
		return nil, ErrorBadRequest
	}

	a, ok := vars["a"]
	if !ok {
		return nil, ErrorBadRequest
	}
	b, ok := vars["b"]
	if !ok {
		return nil, ErrorBadRequest
	}

	return endpoint.StringRequest{
		RequestType: requestType,
		A:           a,
		B:           b,
	}, nil
}
