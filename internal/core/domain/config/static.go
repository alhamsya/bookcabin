package config

import "time"

type Static struct {
	Env      string      `mapstructure:"env"`
	Frontend Frontend    `mapstructure:"frontend"`
	App      App         `mapstructure:"app"`
	Upstream Upstream    `mapstructure:"upstream" validate:"required"`
	Redis    StaticRedis `mapstructure:"redis"`
}

type App struct {
	Rest Rest `mapstructure:"rest"`
}
type Rest struct {
	Port        int           `mapstructure:"port"`
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
	Host                string        `mapstructure:"host"`
	Port                int           `mapstructure:"port"`
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
}
