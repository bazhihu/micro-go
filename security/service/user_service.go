package service

/**
用于获取用户信息
*/

import (
	"context"
	"errors"
	"micro-go/security/model"
)

var (
	ErrUserNotExist = errors.New("username is not exist")
	ErrPassword     = errors.New("invalid password")
)

// define user service interface
type UserDetailsService interface {
	// get UserDetails by username
	GetUserDetailByUsername(ctx context.Context, username, password string) (*model.UserDetails, error)
}

// implement Service interface
type InMemoryUserDetailsService struct {
	userDetailsDict map[string]*model.UserDetails
}

func (service *InMemoryUserDetailsService) GetUserDetailByUsername(ctx context.Context, username, password string) (*model.UserDetails, error) {
	// 根据username 获取用户信息
	userDetails, ok := service.userDetailsDict[username]
	if ok {
		// 比较 password 是否匹配
		if userDetails.Password == password {
			return userDetails, nil
		} else {
			return nil, ErrPassword
		}
	} else {
		return nil, ErrUserNotExist
	}
}

func NewInMemoryUserDetailsService(userDetailsList []*model.UserDetails) *InMemoryUserDetailsService {
	userDetailsDict := make(map[string]*model.UserDetails)

	if userDetailsDict != nil {
		for _, value := range userDetailsDict {
			userDetailsDict[value.Username] = value
		}
	}
	return &InMemoryUserDetailsService{userDetailsDict: userDetailsDict}
}
