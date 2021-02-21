package transport

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"micro-go/security/endpoint"
	"micro-go/security/service"
	"net/http"
)

// 错误信息
var (
	ErrorBadRequest         = errors.New("invalid request parameter")
	ErrorGrantTypeRequest   = errors.New("invalid request grant type")
	ErrorTokenRequest       = errors.New("invalid request token")
	ErrInvalidClientRequest = errors.New("invalid client message")
)

func MakeHttpHandler(ctx context.Context, endpoints endpoint.OAuth2Endpoints, tokenService service.TokenService, clientService service.ClientDetailsService, logger log.Logger) http.Handler {
	r := mux.NewRouter()

	options := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
	}

	clientAuthorizationOptions := []kithttp.ServerOption{
		kithttp.ServerBefore(makeClientAuthorizationConText(clientService, logger)),
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
	}

	r.Methods("POST").Path("/oauth/token").Handler(kithttp.NewServer(
		endpoints.TokenEndpoint,
		decodeTokenRequest,
		encodeJsonResponse,
		clientAuthorizationOptions...,
	))
	r.Methods("POST").Path("/oauth/check_token").Handler(kithttp.NewServer(
		endpoints.CheckTokenEndpoint,
		decodeCheckTokenRequest,
		encodeJsonResponse,
		clientAuthorizationOptions...,
	))

	oauth2AuthorizationOptions := []kithttp.ServerOption{
		kithttp.ServerBefore(makeOAuth2AuthorizationContext(tokenService, logger)),
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
	}
	r.Methods("Get").Path("/simple").Handler(kithttp.NewServer(
		endpoints.SimpleEndpoint,
		decodeSimpleRequest,
		encodeJsonResponse,
		oauth2AuthorizationOptions...,
	))
	r.Methods("Get").Path("/admin").Handler(kithttp.NewServer(
		endpoints.AdminEndpoint,
		decodeAdminRequest,
		encodeJsonResponse,
		oauth2AuthorizationOptions...,
	))

	// create health check handler
	r.Methods("GET").Path("/health").Handler(kithttp.NewServer(
		endpoints.HealthCheckEndpoint,
		decodeSimpleRequest,
		encodeJsonResponse,
		options...,
	))

	// 通过滑动窗口模式查询监控数据 统计一段时间内的调用次数、失败次数、超时次数、和被拒绝次数（执行池已满时请求被拒绝）
	r.Path("/metrics").Handler(promhttp.Handler())

	// 添加hytrix 监控数据
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	r.Handle("/hystrix/stream", hystrixStreamHandler)
	return r
}

// 根据客户端ID 密钥获取客户端信息
func makeClientAuthorizationConText(clientService service.ClientDetailsService, logger log.Logger) kithttp.RequestFunc {
	return func(ctx context.Context, request *http.Request) context.Context {
		if clientId, clientSecret, ok := request.BasicAuth(); ok {
			clientDetails, err := clientService.GetClientDetailByClientId(ctx, clientId, clientSecret)
			if err == nil {
				return context.WithValue(ctx, endpoint.OAuth2ClientDetailsKey, clientDetails)
			}
		}
		return context.WithValue(ctx, endpoint.OAuth2ErrorKey, ErrInvalidClientRequest)
	}
}

// 根据令牌获取对应的用户信息和客户端信息
func makeOAuth2AuthorizationContext(tokenService service.TokenService, logger log.Logger) kithttp.RequestFunc {
	return func(ctx context.Context, request *http.Request) context.Context {
		// 获取访问令牌
		accessTokenValue := request.Header.Get("Authorization")
		var err error
		if accessTokenValue != "" {
			// 获取令牌对应的用户信息和客户端信息
			oauth2Details, err := tokenService.GetOAuth2DetailsByAccessToken(accessTokenValue)
			if err == nil {
				return context.WithValue(ctx, endpoint.OAuth2DetailsKey, oauth2Details)
			}
		} else {
			err = ErrorTokenRequest
		}
		return context.WithValue(ctx, endpoint.OAuth2ErrorKey, err)
	}
}

func encodeJsonResponse(ctx context.Context, writer http.ResponseWriter, i interface{}) error {
	writer.Header().Set("Content-Type", "application/json;charset=utf-8")
	return json.NewEncoder(writer).Encode(i)
}

func decodeSimpleRequest(ctx context.Context, request2 *http.Request) (request interface{}, err error) {
	return &endpoint.SimpleRequest{}, nil
}

func decodeAdminRequest(ctx context.Context, request2 *http.Request) (request interface{}, err error) {
	return &endpoint.AdminRequest{}, nil
}

func decodeCheckTokenRequest(ctx context.Context, request2 *http.Request) (request interface{}, err error) {
	tokenValue := request2.URL.Query().Get("token")
	if tokenValue == "" {
		return nil, ErrorTokenRequest
	}
	return &endpoint.CheckTokenRequest{
		Token: tokenValue,
	}, nil
}

func decodeTokenRequest(ctx context.Context, request2 *http.Request) (request interface{}, err error) {
	grantType := request2.URL.Query().Get("grant_type")
	if grantType == "" {
		return nil, ErrorGrantTypeRequest
	}
	return &endpoint.TokenRequest{
		GrantType: grantType,
		Reader:    request2,
	}, nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
}
