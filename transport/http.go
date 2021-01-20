package transport

import (
	"context"
	"errors"
	"log"
	"net/http"
)

/**
项目提供的服务方式
*/

var (
	ErrorBadRequest = errors.New("invalid request parameter")
)

func MakeHttpHandler(ctx context.Context, logger log.Logger) http.Handler {

	return nil
}
