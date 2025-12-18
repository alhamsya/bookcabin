package xhtpp

import "net/http"

type HTTPClient struct {
	Client         http.Client
	ThirdPartyName string
}

func NewHTTPClient(thirdPartyName string, roundTripper http.RoundTripper) (*HTTPClient, error) {
	return nil, nil
}
