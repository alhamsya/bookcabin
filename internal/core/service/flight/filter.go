package flight

import (
	"slices"
	"strings"

	"github.com/alhamsya/bookcabin/internal/core/domain/constant"
	"github.com/alhamsya/bookcabin/internal/core/domain/flight"
	"github.com/alhamsya/bookcabin/internal/core/domain/request"
	"github.com/alhamsya/bookcabin/lib/util"
)

func applyFilter(info modelFlight.Info, req *modelRequest.ReqSearchFlight) bool {
	// origin & destination
	if !strings.EqualFold(strings.TrimSpace(info.Route.Origin), strings.TrimSpace(req.Origin)) {
		return false
	}
	if !strings.EqualFold(strings.TrimSpace(info.Route.Destination), strings.TrimSpace(req.Destination)) {
		return false
	}

	// Departure date (YYYY-MM-DD)
	if req.DepartureDate != "" {
		if info.Schedule.DepartureTime.Format(constant.DateOnly) != req.DepartureDate {
			return false
		}
	}

	f := req.Filters
	// Price range
	if f.MinPrice > 0 && info.Price.Amount < f.MinPrice {
		return false
	}
	if f.MaxPrice > 0 && info.Price.Amount > f.MaxPrice {
		return false
	}

	// Stops (list)
	if len(f.Stops) > 0 && !slices.Contains(f.Stops, info.Stops) {
		return false
	}

	// Airlines (by code)
	if len(f.Airlines) > 0 && !slices.Contains(f.Airlines, info.Airline.Code) {
		return false
	}

	// Duration
	if f.MaxDurationMinutes > 0 && info.Duration.TotalMinutes > f.MaxDurationMinutes {
		return false
	}

	// Departure time window (uses only time-of-day)
	if !f.DepartureTime.From.IsZero() && !f.DepartureTime.To.IsZero() {
		if !util.WithinTimeWindow(info.Schedule.DepartureTime, f.DepartureTime.From, f.DepartureTime.To) {
			return false
		}
	}

	// Arrival time window (uses only time-of-day)
	if !f.ArrivalTime.From.IsZero() && !f.ArrivalTime.To.IsZero() {
		if !util.WithinTimeWindow(info.Schedule.ArrivalTime, f.ArrivalTime.From, f.ArrivalTime.To) {
			return false
		}
	}
	return true
}
