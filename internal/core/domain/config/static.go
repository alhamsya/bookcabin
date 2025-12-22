package config

import (
	"github.com/alhamsya/bookcabin/pkg/manager/xhttp"
	"time"
)

type Static struct {
	Env      string       `mapstructure:"env" validated:"required"`
	Frontend *Frontend    `mapstructure:"frontend"`
	App      *App         `mapstructure:"app" validated:"required"`
	Upstream *Upstream    `mapstructure:"upstream" validate:"required"`
	Redis    *StaticRedis `mapstructure:"redis" validated:"required"`
}

type App struct {
	Rest Rest `mapstructure:"rest"`
}
type Rest struct {
	Port        int           `mapstructure:"port" validated:"required"`
	ReadTimeout time.Duration `mapstructure:"read-timeout"`
	IdleTimeout time.Duration `mapstructure:"idle-timeout"`
	Limiter     Limiter       `mapstructure:"limiter"`
}

type Limiter struct {
	Max        int           `mapstructure:"max"`
	Expiration time.Duration `mapstructure:"expiration"`
}

type Frontend struct {
	URL string `mapstructure:"url"`
}

type StaticRedis struct {
	Host                string        `mapstructure:"host" validated:"required"`
	Port                int           `mapstructure:"port" validated:"required"`
	DB                  int           `mapstructure:"db"`
	IdleTimeout         string        `mapstructure:"idle-timeout"`
	PoolTimeout         string        `mapstructure:"pool-timeout"`
	ReadTimeout         string        `mapstructure:"read-timeout"`
	PoolFIFO            bool          `mapstructure:"pool-fifo"`
	PoolSize            int           `mapstructure:"pool-size"`
	IdleTimeoutDuration time.Duration `mapstructure:"idle-timeout"`
	PoolTimeoutDuration time.Duration `mapstructure:"pool-timeout"`
	ReadTimeoutDuration time.Duration `mapstructure:"read-timeout"`
}

type Upstream struct {
	AirAsia *xhttp.UpstreamConfig `mapstructure:"airasia" json:"airasia" validated:"required"`
	Batik   *xhttp.UpstreamConfig `mapstructure:"batik" json:"batik" validated:"required"`
	Garuda  *xhttp.UpstreamConfig `mapstructure:"garuda" json:"garuda" validated:"required"`
	Lion    *xhttp.UpstreamConfig `mapstructure:"lion" json:"lion" validated:"required"`
}
