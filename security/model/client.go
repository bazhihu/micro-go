package model

// 客户端信息
type ClientDetails struct {
	// client的标示
	ClientId string
	// client的密钥
	ClientSecret string
	// 访问令牌有效时间 秒
	AccessTokenValiditySeconds int
	// 刷新令牌有效时间 秒
	RefreshTokenValiditySeconds int
	// 重定向地址 授权码类型中使用
	RegisteredRedirectUri string
	// 可以使用的授权类型
	AuthorizedGrantTypes []string
}
