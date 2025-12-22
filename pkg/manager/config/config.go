package config

import "github.com/alhamsya/bookcabin/internal/core/domain/config"

type Application struct {
	Credential config.Credential `mapstructure:"credential"`
	Static     config.Static     `mapstructure:"static"`
}
