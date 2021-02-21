package main

import (
	"context"
	"flag"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"micro-go/common/discover"
	"micro-go/security/config"
	"micro-go/security/endpoint"
	"micro-go/security/model"
	"micro-go/security/service"
	"micro-go/security/transport"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	var (
		servicePort = flag.Int("service.port", 10098, "service port")
		serviceHost = flag.String("service.host", "127.0.0.1", "service host")

		consulPort = flag.Int("consul.port", 8500, "consul port")
		consulHost = flag.String("consul.port", "127.0.0.1", "consul host")

		serviceName = flag.String("service.name", "oauth", "service name")
	)

	flag.Parse()

	ctx := context.Background()
	errChan := make(chan error)

	// 服务发现与注册
	var discoveryClient discover.DiscoveryClient
	discoveryClient, err := discover.NewKitDiscoverClient(*consulHost, *consulPort)
	if err != nil {
		config.Logger.Println("Get Consul Client failed")
		os.Exit(-1)
	}

	var (
		tokenService         service.TokenService
		tokenGranter         service.TokenGranter
		tokenEnhancer        service.TokenEnhancer
		tokenStore           service.TokenStore
		userDetailsService   service.UserDetailsService
		clientDetailsService service.ClientDetailsService
		srv                  service.Service
	)

	tokenEnhancer = service.NewJwtTokenEnhancer("secret")
	tokenStore = service.NewJwtTokenStore(tokenEnhancer.(*service.JwtTokenEnhancer))
	tokenService = service.NewTokenService(tokenStore, tokenEnhancer)

	userDetailsService = service.NewInMemoryUserDetailsService([]*model.UserDetails{
		{Username: "simple", Password: "123456", UserId: 1, Authorities: []string{"Simple"}},
		{Username: "admin", Password: "123456", UserId: 1, Authorities: []string{"Admin"}},
	})

	clientDetailsService = service.NewInMemoryClientDetailService([]*model.ClientDetails{
		{ClientId: "clientId", ClientSecret: "clientSecret",
			AccessTokenValiditySeconds: 1800, RefreshTokenValiditySeconds: 18000,
			RegisteredRedirectUri: "http://127.0.0.1", AuthorizedGrantTypes: []string{"password", "refresh_token"}},
	})

	tokenGranter = service.NewComposeTokenGranter(map[string]service.TokenGranter{
		"password":      service.NewUsernamePasswordTokenGranter("password", userDetailsService, tokenService),
		"refresh_token": service.NewRefreshGranter("refresh_token", userDetailsService, tokenService),
	})

	// endpoint
	tokenEndpoint := endpoint.MakeTokenEndpoint(tokenGranter, clientDetailsService)
	tokenEndpoint = endpoint.MakeClientAuthorizationMiddleware(config.KitLogger)(tokenEndpoint)
	checkTokenEndpoint := endpoint.MakeCheckTokenEndpoint(tokenService)
	checkTokenEndpoint = endpoint.MakeClientAuthorizationMiddleware(config.KitLogger)(checkTokenEndpoint)

	srv = service.NewCommonService()

	simpleEndpoint := endpoint.MakeSimpleEndpoint(srv)
	simpleEndpoint = endpoint.MakeOAuth2AuthorizationMiddleware(config.KitLogger)(simpleEndpoint)
	adminEndpoint := endpoint.MakeAdminEndpoint(srv)
	adminEndpoint = endpoint.MakeOAuth2AuthorizationMiddleware(config.KitLogger)(adminEndpoint)
	adminEndpoint = endpoint.MakeAuthorityAuthorizationMiddleware("Admin", config.KitLogger)(adminEndpoint)

	// 创建健康检查的Endpoint
	healthEndpoint := endpoint.MakeHealthCheckEndpoint(srv)

	endpts := endpoint.OAuth2Endpoints{
		TokenEndpoint:       tokenEndpoint,
		CheckTokenEndpoint:  checkTokenEndpoint,
		HealthCheckEndpoint: healthEndpoint,
		SimpleEndpoint:      simpleEndpoint,
		AdminEndpoint:       adminEndpoint,
	}

	// 根据transport 创建http.Handler
	r := transport.MakeHttpHandler(ctx, endpts, tokenService, clientDetailsService, config.KitLogger)

	instanceId := *serviceName + "-" + uuid.NewV4().String()

	// http server
	go func() {
		config.Logger.Println("Http Server start at port:" + strconv.Itoa(*servicePort))
		// 启动前执行注册
		if !discoveryClient.Register(*serviceName, instanceId, "/health", *serviceHost, *servicePort, nil, config.Logger) {
			config.Logger.Printf("use-string-service for servcie %s failed.", serviceName)
			// 注册失败，服务启动失败
			os.Exit(-1)
		}
		handler := r
		errChan <- http.ListenAndServe(":"+strconv.Itoa(*servicePort), handler)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	error := <-errChan

	// 服务退出取消注册
	discoveryClient.DeRegister(instanceId, config.Logger)
	config.Logger.Println(error)
}
