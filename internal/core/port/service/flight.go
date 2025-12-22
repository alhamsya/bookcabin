package port

import (
	"context"

	"github.com/alhamsya/bookcabin/internal/core/domain/flight"
	"github.com/alhamsya/bookcabin/internal/core/domain/response"
)

//go:generate mockgen -package=mockservice -source=$GOFILE -destination=../.././../mock/service/$GOFILE
// FlightService is an interface for interacting with flight business logic
type FlightService interface {
	Search(ctx context.Context, param *modelFlight.ReqSearch) (modelResponse.Common, error)
}
