package service

import (
	"errors"
	"strings"
)

// Service constants
const (
	StrMaxSize = 1024
)

// Service errors
var (
	ErrMaxSize  = errors.New("maximum size of 1024 bytes exceeded")
	ErrStrValue = errors.New("string len too long")
)

type Service interface {
	// Concat a and b
	Concat(req StringRequest, ret *string) error
	// a,b common string value
	Diff(req StringRequest, ret *string) error

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

func (s StringService) Concat(req StringRequest, ret *string) error {
	// test for length overflow
	if len(req.A)+len(req.B) > StrMaxSize {
		*ret = ""
		return ErrMaxSize
	}
	*ret = req.A + req.B
	return nil
}

func (s StringService) Diff(req StringRequest, ret *string) error {
	if len(req.A) < 1 || len(req.B) < 1 {
		*ret = ""
		return nil
	}

	var res, aa, bb string
	if len(req.A) >= len(req.B) {
		aa = req.A
		bb = req.B
	} else {
		aa = req.B
		bb = req.A
	}

	for _, char := range bb {
		if strings.Contains(aa, string(char)) {
			res = res + string(char)
		}
	}

	*ret = res
	return nil
}

func (s StringService) HealthCheck() bool {
	return true
}

// ServiceMiddleware define service middleware
type ServiceMiddleware func(Service) Service
