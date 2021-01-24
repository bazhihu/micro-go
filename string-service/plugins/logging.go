package plugins

import (
	"github.com/go-kit/kit/log"
	"micro-go/string-service/service"
	"time"
)

// 组装了服务接口和日志接口
type loggingMiddleware struct {
	service.Service
	logger log.Logger
}

func LoggingMiddleware(logger log.Logger) service.ServiceMiddleware {
	return func(i service.Service) service.Service {
		return loggingMiddleware{i, logger}
	}
}

func (mw loggingMiddleware) Concat(a, b string) (ret string, err error) {
	// 函数执行完后打印日志
	defer func(begin time.Time) {
		mw.logger.Log(
			"function", "Concat",
			"a", a,
			"b", b,
			"result", ret,
			"took", time.Since(begin),
		)
	}(time.Now())

	ret, err = mw.Service.Concat(a, b)
	return
}

func (mw loggingMiddleware) Diff(a, b string) (ret string, err error) {
	// 函数执行完后打印日志
	defer func(begin time.Time) {
		mw.logger.Log(
			"function", "Diff",
			"a", a,
			"b", b,
			"result", ret,
			"took", time.Since(begin),
		)
	}(time.Now())
	ret, err = mw.Service.Diff(a, b)
	return
}

func (mw loggingMiddleware) HealthCheck() (result bool) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"function", "HealthCheck",
			"result", result,
			"took", time.Since(begin),
		)
	}(time.Now())
	result = mw.Service.HealthCheck()
	return
}
