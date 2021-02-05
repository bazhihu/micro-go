package service

import (
	"context"
	"errors"
	"micro-go/security/model"
)

/**
用于获取客户端信息
*/

var (
	ErrClientNotExist = errors.New("clientId is not exist")
	ErrClientSecret   = errors.New("invalid clientSecret")
)

// 客户端信息服务接口
type ClientDetailsService interface {
	GetClientDetailByClientId(ctx context.Context, clientId string, clientSecret string) (*model.ClientDetails, error)
}

// 客户端信息服务对象
type InMemoryClientDetailsService struct {
	clientDetailsDict map[string]*model.ClientDetails
}

func NewInMemoryClientDetailService(clientDetailsList []*model.ClientDetails) *InMemoryClientDetailsService {
	clientDetailsDict := make(map[string]*model.ClientDetails)

	if clientDetailsList != nil {
		for _, value := range clientDetailsList {
			clientDetailsDict[value.ClientId] = value
		}
	}

	return &InMemoryClientDetailsService{clientDetailsDict: clientDetailsDict}
}

// 根据客户端ID 获取信息
func (service *InMemoryClientDetailsService) GetClientDetailByClientId(ctx context.Context, clientId string, clientSecret string) (*model.ClientDetails, error) {
	// 根据clientId 获取ClientDetails
	clientDetails, ok := service.clientDetailsDict[clientId]
	if ok {
		// 比较clientSecret 是否正确
		if clientDetails.ClientSecret == clientSecret {
			return clientDetails, nil
		} else {
			return nil, ErrClientSecret
		}
	}
	return nil, ErrClientNotExist
}
