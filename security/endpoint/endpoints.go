package endpoint

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"micro-go/security/model"
	"micro-go/security/service"
	"net/http"
)

const (
	OAuth2ErrorKey         = "OAuth2Error"
	OAuth2DetailsKey       = "OAuth2Details"
	OAuth2ClientDetailsKey = "OAuth2ClientDetails"
)

var (
	ErrInvalidClientRequest = errors.New("invalid client message")
	ErrInvalidUserRequest   = errors.New("invalid user message")
	ErrNotPermit            = errors.New("not permit")
)

type OAuth2Endpoints struct {
	TokenEndpoint       endpoint.Endpoint
	CheckTokenEndpoint  endpoint.Endpoint
	HealthCheckEndpoint endpoint.Endpoint
	SimpleEndpoint      endpoint.Endpoint
	AdminEndpoint       endpoint.Endpoint
}

// 验证客户端信息
func MakeClientAuthorizationMiddleware(logger log.Logger) endpoint.Middleware {
	return func(e endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if err, ok := ctx.Value(OAuth2ErrorKey).(error); ok {
				return nil, err
			}
			if _, ok := ctx.Value(OAuth2ClientDetailsKey).(*model.ClientDetails); !ok {
				return nil, ErrInvalidClientRequest
			}
			return e(ctx, request)
		}
	}
}

// 令牌信息
func MakeOAuth2AuthorizationMiddleware(logger log.Logger) endpoint.Middleware {
	return func(e endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if err, ok := ctx.Value(OAuth2ErrorKey).(error); ok {
				return nil, err
			}
			if _, ok := ctx.Value(OAuth2DetailsKey).(*model.OAuth2Details); !ok {
				return nil, ErrInvalidUserRequest
			}
			return e(ctx, request)
		}
	}
}

// 用户信息
func MakeAuthorityAuthorizationMiddleware(authority string, logger log.Logger) endpoint.Middleware {
	return func(e endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if err, ok := ctx.Value(OAuth2ErrorKey).(error); ok {
				return nil, err
			}
			if details, ok := ctx.Value(OAuth2DetailsKey).(*model.OAuth2Details); !ok {
				return nil, ErrInvalidClientRequest
			} else {
				for _, value := range details.User.Authorities {
					if value == authority {
						return e(ctx, request)
					}
				}
				return nil, ErrNotPermit
			}
		}
	}
}

type TokenRequest struct {
	GrantType string
	Reader    *http.Request
}

type TokenResponse struct {
	AccessToken *model.OAuth2Token `json:"access_token"`
	Error       string             `json:"error"`
}

func MakeTokenEndpoint(svc service.TokenGranter, clientService service.ClientDetailsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*TokenRequest)
		token, err := svc.Grant(ctx, req.GrantType, ctx.Value(OAuth2ClientDetailsKey).(*model.ClientDetails), req.Reader)

		var errString = ""
		if err != nil {
			errString = err.Error()
		}

		return TokenResponse{AccessToken: token, Error: errString}, nil
	}
}

// ----------------------------

type CheckTokenRequest struct {
	Token         string
	ClientDetails model.ClientDetails
}

type CheckTokenResponse struct {
	OAuthDetails *model.OAuth2Details `json:"o_auth_details"`
	Err          string               `json:"error"`
}

func MakeCheckTokenEndpoint(svc service.TokenService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*CheckTokenRequest)
		tokenDetails, err := svc.GetOAuth2DetailsByAccessToken(req.Token)

		var errString = ""
		if err != nil {
			errString = err.Error()
		}

		return CheckTokenResponse{OAuthDetails: tokenDetails, Err: errString}, nil
	}
}

// -----------------------------

type SimpleRequest struct {
}

type SimpleResponse struct {
	Result string `json:"result"`
	Error  string `json:"error"`
}

func MakeSimpleEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		result := svc.SimpleData(ctx.Value(OAuth2DetailsKey).(*model.OAuth2Details).User.Username)
		return &SimpleResponse{Result: result}, nil
	}
}

// -----------------------

type AdminRequest struct {
}

type AdminResponse struct {
	Result string `json:"result"`
	Error  string `json:"error"`
}

func MakeAdminEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		result := svc.AdminData(ctx.Value(OAuth2DetailsKey).(*model.OAuth2Details).User.Username)
		return &AdminResponse{Result: result}, nil
	}
}

// ------------------------------

// HealthRequest 健康检查请求结构
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
