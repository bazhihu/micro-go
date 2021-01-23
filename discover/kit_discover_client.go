package discover

import (
	"github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
)

type KitDiscoverClient struct {
	Host   string // Consul Host
	Port   int    // Consul Port
	client consul.Client
}

func NewKitDiscoverClient(consulHost string, consulPort int) (DiscoveryClient, error) {
	// 通过Consul Host和Consul Port创建一个consul.Client
	consulConfig := api.DefaultConfig()

	return
}
