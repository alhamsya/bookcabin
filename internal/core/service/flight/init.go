package flight

import "github.com/alhamsya/bookcabin/internal/core/port"

type FlightParam struct {
}
type FlightService struct {
	cache port.CacheRepo
}

func NewFlightService(cache port.CacheRepo) *FlightService {
	return &FlightService{
		cache,
	}
}
