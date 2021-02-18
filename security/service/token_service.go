package service

import (
	"context"
	"micro-go/security/model"
	"net/http"
)

/**
根据授权类型使用不同方式
对用户和客户端信息进行认证
生成并管理令牌，使用TokenStore存储令牌
*/

// 错误信息
var ()

// 令牌授予
type TokenGranter interface {
	Grant(ctx context.Context, grantType string, client *model.ClientDetails, reader *http.Request) (*model.OAuth2Token, error)
}

// 令牌桶
type ComposeTokenGranter struct {
	TokenGrantDict map[string]TokenGranter
}

func (c ComposeTokenGranter) Grant(ctx context.Context, grantType string, client *model.ClientDetails, reader *http.Request) (*model.OAuth2Token, error) {
	panic("implement me")
}

func NewComposeTokenGranter(tokenGrantDict map[string]TokenGranter) TokenGranter {
	return &ComposeTokenGranter{TokenGrantDict: tokenGrantDict}
}

// 用户密码令牌验证授予
type UsernamePasswordTokenGranter struct {
	supportGrantType   string
	userDetailsService UserDetailsService
	tokenService       TokenService
}

type TokenService interface {
	// 根据访问令牌获取对应的用户信息和客户端信息
	GetOAuth2DetailsByAccessToken(tokenValue string) (*model.OAuth2Details, error)
	// 根据用户信息和客户端生成访问令牌
	CreateAccessToken(oauth2Details *model.OAuth2Details) (*model.OAuth2Token, error)
	// 根据刷新令牌获取访问令牌
	RefreshAccessToken(refreshTokenValue string) (*model.OAuth2Token, error)
	// 根据用户信息和客户端信息获取已生成访问令牌
	GetAccessToken(details *model.OAuth2Details) (*model.OAuth2Token, error)
	// 根据访问令牌获取访问令牌结构体
	ReadAccessToken(tokenValue string) (*model.OAuth2Token, error)
}
