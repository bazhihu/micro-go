package service

import (
	"errors"
	"strings"
)

// service constants
const (
	StrMaxSize = 1024
)

// service errors
var (
	ErrMaxSize  = errors.New("maximum size of 1024 bytes exceeded")
	ErrStrValue = errors.New("maximum size of 1024 bytes exceeded")
)

// service interface
type Service interface {
	// 连接字符串 a,b
	Concat(a, b string) (string, error)
	// 字符串a, b 公共字符
	Diff(a, b string) (string, error)
	// 健康检查
	HealthCheck() bool
}

// ArithmeticService implement Service interface
type StringService struct {
}

func (s StringService) Concat(a, b string) (string, error) {
	if len(a)+len(b) > StrMaxSize {
		return "", ErrMaxSize
	}
	return a + b, nil
}

func (s StringService) Diff(a, b string) (string, error) {
	if len(a) < 1 || len(b) < 1 {
		return "", nil
	}
	var res, aa, bb string
	if len(a) > len(b) {
		aa = a
		bb = b
	} else {
		aa = b
		bb = a
	}

	for _, char := range bb {
		if strings.Contains(aa, string(char)) {
			res = res + string(char)
		}
	}

	return res, nil
}

func (s StringService) HealthCheck() bool {
	return true
}

type ServiceMiddleware func(Service) Service
