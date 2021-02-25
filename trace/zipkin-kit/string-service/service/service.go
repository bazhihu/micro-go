package service

import (
	"context"
	"errors"
	"strings"
)

// Service constants
const (
	StrMaxSize = 1024
)

// Service errors
var (
	ErrMaxSize = errors.New("maximum size of 1024 bytes exceeded")
	//ErrStrValue = errors.New("string len too long")
)

type Service interface {
	// Concat a and b
	Concat(a, b string) (string, error)

	// a,b pkg string value
	Diff(ctx context.Context, a, b string) (string, error)

	// HealthCheck check service Health status
	HealthCheck() bool
}

type StringRequest struct {
	A string
	B string
}

type StringResponse struct {
	Result string
	Error  string
}

type StringService struct {
}

func (s StringService) Concat(a, b string) (string, error) {
	// test for length overflow
	if len(a)+len(b) > StrMaxSize {
		return "", ErrMaxSize
	}
	return a + b, nil
}

func (s StringService) Diff(ctx context.Context, a, b string) (string, error) {
	if len(a) < 1 || len(b) < 1 {
		return "", nil
	}
	var res, aa, bb string
	if len(a) >= len(b) {
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

// ServiceMiddleware define service middleware
type ServiceMiddleware func(Service) Service
