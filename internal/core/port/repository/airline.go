package port

import (
	"context"
	modelFlight "github.com/alhamsya/bookcabin/internal/core/domain/flight"
)

type AirAsiaRepo interface {
	GetFlight(ctx context.Context) ([]modelFlight.Info, error)
}

type BatikRepo interface {
	GetFlight(ctx context.Context) ([]modelFlight.Info, error)
}

type GarudaRepo interface {
	GetFlight(ctx context.Context) ([]modelFlight.Info, error)
}

type LionRepo interface {
	GetFlight(ctx context.Context) ([]modelFlight.Info, error)
}
