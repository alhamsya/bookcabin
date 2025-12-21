package xhttp

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/pkg/errors"
)

type retryableStatusError struct {
	Code int
}

func (e retryableStatusError) Error() string {
	return fmt.Sprintf("retryable status: %d", e.Code)
}

func NewClient(cfg *UpstreamConfig, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{cfg: cfg, http: httpClient}
}

func (c *Client) Do(ctx context.Context, call Call) (*Result, error) {
	ep, ok := c.cfg.Endpoint[call.Endpoint]
	if !ok {
		return nil, fmt.Errorf("endpoint %q not found (upstream=%q)", call.Endpoint, call.Upstream)
	}

	rr := resolveRetry(ep.Retry, c.cfg.Global)

	req, err := c.buildRequest(call, ep)
	if err != nil {
		return nil, errors.Wrap(err, "failed buildRequest")
	}

	var out Result
	err = retry.Do(
		func() error {
			res, errDo := c.doOnce(ctx, req, rr)
			if errDo != nil {
				return errDo
			}
			out = *res
			return nil
		},
		c.retryOptions(ctx, rr)...,
	)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) buildRequest(call Call, ep Endpoint) (*http.Request, error) {
	fullURL, err := buildURL(c.cfg.Host+ep.Path+call.PathSuffix, call.Query)
	if err != nil {
		return nil, err
	}
	return newRetryRequest(call.Method, fullURL, call.Body, call.Headers)
}

func (c *Client) doOnce(ctx context.Context, req *http.Request, rr RetryConfig) (*Result, error) {
	attemptCtx := ctx
	cancel := func() {}
	if rr.Timeout > 0 {
		attemptCtx, cancel = context.WithTimeout(ctx, rr.Timeout)
	}
	defer cancel()

	r2, err := cloneForRetry(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed cloneForRetry")
	}

	resp, err := c.http.Do(r2.WithContext(attemptCtx))
	if err != nil {
		return nil, errors.Wrap(err, "failed do")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= http.StatusInternalServerError {
		_, _ = io.Copy(io.Discard, resp.Body)
		return nil, retryableStatusError{Code: resp.StatusCode}
	}

	resultBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &Result{
		Status: resp.StatusCode,
		Header: resp.Header.Clone(),
		Body:   resultBody,
	}, nil
}

func (c *Client) retryOptions(ctx context.Context, rr RetryConfig) []retry.Option {
	return []retry.Option{
		retry.Context(ctx),
		retry.Attempts(rr.Attempt),
		retry.Delay(rr.Delay),
		retry.DelayType(func(n uint, err error, cfg *retry.Config) time.Duration {
			if rr.Type == RetryTypeExponential {
				return retry.BackOffDelay(n, err, cfg)
			}
			return retry.FixedDelay(n, err, cfg)
		}),
		retry.LastErrorOnly(true),
		retry.RetryIf(isRetryable),
	}
}

func isRetryable(err error) bool {
	if err == nil || errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}
	var ne net.Error
	if errors.As(err, &ne) {
		return true
	}
	var se retryableStatusError
	return errors.As(err, &se)
}

func resolveRetry(endpoint, global *RetryConfig) RetryConfig {
	out := RetryConfig{
		Attempt: 3,
		Delay:   100 * time.Millisecond,
		Timeout: 0,
		Type:    RetryTypeFixed,
	}

	chosen := endpoint
	if chosen == nil {
		chosen = global
	}
	if chosen == nil {
		return out
	}

	if chosen.Attempt > 0 {
		out.Attempt = chosen.Attempt
	}
	if chosen.Type != "" {
		out.Type = chosen.Type
	}
	if chosen.Delay != 0 {
		out.Delay = chosen.Delay
	}
	if chosen.Timeout != 0 {
		out.Timeout = chosen.Timeout
	}

	return out
}

func buildURL(raw string, q map[string]string) (string, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", err
	}
	if len(q) > 0 {
		qq := u.Query()
		for k, v := range q {
			qq.Set(k, v)
		}
		u.RawQuery = qq.Encode()
	}
	return u.String(), nil
}

func newRetryRequest(method, urlStr string, body []byte, headers map[string]string) (*http.Request, error) {
	var rdr io.Reader
	if body != nil {
		rdr = bytes.NewReader(body)
	}

	req, err := http.NewRequest(method, urlStr, rdr)
	if err != nil {
		return nil, err
	}

	// Make body replayable for retries.
	if body != nil {
		req.GetBody = func() (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader(body)), nil
		}
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return req, nil
}

func cloneForRetry(req *http.Request) (*http.Request, error) {
	clone := req.Clone(req.Context())

	if req.Body == nil || req.Body == http.NoBody {
		clone.Body = http.NoBody
		return clone, nil
	}
	if req.GetBody == nil {
		return nil, errors.New("request body is not retryable (GetBody is nil)")
	}

	rc, err := req.GetBody()
	if err != nil {
		return nil, err
	}
	clone.Body = rc
	return clone, nil
}
