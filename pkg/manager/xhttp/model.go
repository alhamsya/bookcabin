package xhttp

import (
	"net/http"
	"time"
)

const (
	RetryTypeExponential = "exponential"
	RetryTypeFixed       = "fixed"
)

type UpstreamType string

const (
	UpstreamTypeREST   UpstreamType = "rest"   // REST API
	UpstreamTypeGRPC   UpstreamType = "grpc"   // gRPC
	UpstreamTypeSocket UpstreamType = "socket" // socket protocol
)

type UpstreamConfig struct {
	Type     UpstreamType        `mapstructure:"type" json:"type" validated:"required"`
	Host     string              `mapstructure:"host" json:"host" validated:"required"`
	Port     int                 `mapstructure:"port" json:"port" validated:"required"`
	Endpoint map[string]Endpoint `mapstructure:"endpoint" json:"endpoint" validated:"required"`

	// Optional
	Extra     map[string]string `mapstructure:"extra" json:"extra"`
	Global    *RetryConfig      `mapstructure:"global-retry" json:"global-retry"`
	EnableTLS bool              `mapstructure:"enable-tls" json:"enable-tls"`
}

type Endpoint struct {
	Path  string       `mapstructure:"path" json:"path" validated:"required"`
	Retry *RetryConfig `mapstructure:"retry" json:"retry"`
}

type RetryConfig struct {
	Attempt uint          `mapstructure:"attempt" json:"attempt" validated:"required"`
	Delay   time.Duration `mapstructure:"delay" json:"delay"`
	Timeout time.Duration `mapstructure:"timeout" json:"timeout"`
	Type    string        `mapstructure:"type" json:"type" validated:"required"`
}

type Call struct {
	Upstream string
	Endpoint string
	Method   string

	PathSuffix string            // optional
	Query      map[string]string // optional
	Headers    map[string]string // optional
	Body       []byte            // optional (nil for no body)
}

type Result struct {
	Status int
	Header http.Header
	Body   []byte
}

type Client struct {
	cfg  *UpstreamConfig
	http *http.Client
}
