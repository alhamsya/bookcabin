package flight

import (
	"github.com/alhamsya/bookcabin/internal/core/port/repository"
	"github.com/alhamsya/bookcabin/pkg/manager/config"
	"github.com/rs/zerolog"
)

type Service struct {
	Cfg   *config.Application
	Cache port.CacheRepo
	Log   zerolog.Logger

	AirAsiaRepo port.AirAsiaRepo
	BatikRepo   port.BatikRepo
	GarudaRepo  port.GarudaRepo
	LionRepo    port.LionRepo
}

func NewFlightService(param *Service) *Service {
	return &Service{
		Cfg:         param.Cfg,
		Cache:       param.Cache,
		Log:         param.Log,
		AirAsiaRepo: param.AirAsiaRepo,
		BatikRepo:   param.BatikRepo,
		GarudaRepo:  param.GarudaRepo,
		LionRepo:    param.LionRepo,
	}
}
