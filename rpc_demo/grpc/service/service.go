package service

import (
	"context"
	"errors"
	"micro-go/rpc_demo/pd"
	"strings"
)

const (
	StrMaxSize = 1024
)

// Service errors
var (
	ErrMaxSize  = errors.New("maximum size of 1024 bytes exceeded")
	ErrStrValue = errors.New("maximum size of 1024 bytes exceeded")
)

type StringService struct{}

func (s *StringService) Concat(ctx context.Context, req *pd.StringRequest) (*pd.StringResponse, error) {
	if len(req.A)+len(req.B) > StrMaxSize {
		response := pd.StringResponse{}
		return &response, nil
	}
	response := pd.StringResponse{Ret: req.A + req.B}
	return &response, nil
}

func (s *StringService) Diff(ctx context.Context, req *pd.StringRequest) (*pd.StringResponse, error) {
	if len(req.A) < 1 || len(req.B) < 1 {
		return &pd.StringResponse{}, nil
	}
	var res, aa, bb string
	if len(req.A) > len(req.B) {
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
	return &pd.StringResponse{Ret: res}, nil
}
