package main

import (
	"context"
	"github.com/go-kit/kit/log"
	"micro-go/trace/zipkin-kit/string-service/service"
	"time"
)

// loggingMiddleware Make a new type
type loggingMiddleware struct {
	service.Service
	logger log.Logger
}

// make logging middleware 中间件
func LoggingMiddleware(logger log.Logger) service.ServiceMiddleware {
	return func(s service.Service) service.Service {
		return loggingMiddleware{s, logger}
	}
}

func (s loggingMiddleware) Concat(a, b string) (ret string, err error) {
	// test for length overflow
	defer func(begin time.Time) {
		s.logger.Log("function", "Concat", "a", a, "b", b, "result", ret, "took", time.Since(begin))
	}(time.Now())
	ret, err = s.Service.Concat(a, b)
	return ret, err
}

func (s loggingMiddleware) Diff(ctx context.Context, a, b string) (ret string, err error) {
	defer func(begin time.Time) {
		s.logger.Log("function", "Diff", "a", a, "b", b, "result", ret, "took", time.Since(begin))
	}(time.Now())

	ret, err = s.Service.Diff(ctx, a, b)
	return
}

func (s loggingMiddleware) HealthCheck() (result bool) {
	defer func(begin time.Time) {
		s.logger.Log("function", "HealthCheck", "result", result, "took", time.Since(begin))
	}(time.Now())
	result = s.Service.HealthCheck()
	return true
}
