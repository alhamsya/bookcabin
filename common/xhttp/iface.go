package xhttp

import (
	"context"
	"net/http"

	"github.com/avast/retry-go/v4"
)

type retryConfig struct {
	retryOpts []retry.Option
}

type callOptions struct {
	retryCfg retryConfig
}

type CallOption func(*callOptions)

type HTTPCall interface {
	CallAPI(ctx context.Context, req *http.Request, opts ...CallOption) (int, []byte, error)
}
