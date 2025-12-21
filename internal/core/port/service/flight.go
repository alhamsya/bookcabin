package port

import (
	"context"

	"github.com/alhamsya/bookcabin/internal/core/domain/request"
	"github.com/alhamsya/bookcabin/internal/core/domain/response"
)

// FlightService is an interface for interacting with flight business logic
type FlightService interface {
	Search(ctx context.Context, param *modelRequest.ReqSearchFlight) (modelResponse.Common, error)
}
