package port

import (
	"context"
	"github.com/alhamsya/bookcabin/internal/core/domain/airline"
)

type LionRepo interface {
	GetFlight(ctx context.Context) (*airline.GarudaFlightResp, error)
}

type GarudaRepo interface {
	GetFlight(ctx context.Context) (*airline.GarudaFlightResp, error)
}

type AirAsiaRepo interface {
	GetFlight(ctx context.Context) (*airline.GarudaFlightResp, error)
}

type BatikRepo interface {
	GetFlight(ctx context.Context) (*airline.GarudaFlightResp, error)
}
