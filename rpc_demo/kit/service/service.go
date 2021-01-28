package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

// constants
const (
	StrMaxSize = 1024
)

// service errors
var (
	ErrMaxSize  = errors.New("maximum size of 1024 bytes exceeded")
	ErrStrValue = errors.New("maximum size of 1024 bytes exceeded")
)

type Service interface {
	// concat a and b
	Concat(ctx context.Context, a, b string) (string, error)

	// a,b pkg string value
	Diff(ctx context.Context, a, b string) (string, error)
}

// service
type StringService struct {
}

// concat
func (s StringService) Concat(ctx context.Context, a, b string) (string, error) {
	if len(a)+len(b) > StrMaxSize {
		return "", ErrMaxSize
	}
	fmt.Printf("SstringService Concat return %s", a+b)
	return a + b, nil
}

// diff
func (s StringService) Diff(ctx context.Context, a, b string) (string, error) {
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

// ServiceMiddleware define service middleware
type ServiceMiddleware func(Service) Service
