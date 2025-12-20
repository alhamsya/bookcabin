package port

import (
	"context"

	"github.com/alhamsya/bookcabin/internal/core/domain/response"
)

// FlightRepository is an interface for interacting with flight database
type FlightRepository interface {
}

// FlightService is an interface for interacting with flight business logic
type FlightService interface {
	Search(ctx context.Context) (*response.Common, error)
}
