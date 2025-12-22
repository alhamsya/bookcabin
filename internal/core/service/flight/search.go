package flight

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/alhamsya/bookcabin/internal/core/domain/constant"
	"github.com/alhamsya/bookcabin/internal/core/domain/flight"
	"github.com/alhamsya/bookcabin/internal/core/domain/response"
	"github.com/alhamsya/bookcabin/pkg/util"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"net/http"
	"slices"
	"sort"
	"strings"
	"sync"
)

func (s *Service) Search(ctx context.Context, param *modelFlight.ReqSearch) (modelResponse.Common, error) {
	// 1. check cache in L1 in-memory

	// 2. call providers airline using asynchronous
	// 		2.1. check data in cache in L2 redis
	// 		2.2. call provider using retry mechanism
	// 		2.3. save cache in L2 redis if provider success
	listFlight, err := s.callProviders(ctx)
	if err != nil {
		return modelResponse.Common{}, errors.Wrap(err, "failed call providers")
	}

	// 3. filter results providers
	out := make([]modelFlight.Info, 0, len(listFlight))
	for _, info := range listFlight {
		if !applyFilter(info, param) {
			continue
		}

		// 4. calculate best value score (ranking)
		applyRanking(&info)

		out = append(out, info)
	}

	// 5. sorting best value score
	sort.Slice(out, func(i, j int) bool {
		return out[i].BestValueScore > out[j].BestValueScore
	})

	// 6. set metadata

	// 7. save to cache in L1 in-memory

	return modelResponse.Common{
		HttpCode: http.StatusOK,
		Data:     out,
	}, nil
}

func (s *Service) callProviders(ctx context.Context) ([]modelFlight.Info, error) {
	var (
		wg     sync.WaitGroup
		mu     sync.Mutex
		out    []modelFlight.Info
		failed int
	)

	providers := map[string]func(context.Context) ([]modelFlight.Info, error){
		"airasia": s.AirAsiaRepo.GetFlight,
		"batik":   s.BatikRepo.GetFlight,
		"garuda":  s.GarudaRepo.GetFlight,
		"lion":    s.LionRepo.GetFlight,
	}

	for name, fn := range providers {
		wg.Add(1)
		go func(name string, fn func(context.Context) ([]modelFlight.Info, error)) {
			defer wg.Done()

			log := zerolog.Ctx(ctx).With().Str("provider", name).Logger()

			cached, err := s.Cache.Get(ctx, fmt.Sprintf(constant.RedisKeyProvider, name))
			if err != nil {
				log.Warn().Err(err).Msg("redis get failed")
			} else if len(cached) > 0 {
				var cachedFlights []modelFlight.Info
				if err = json.Unmarshal(cached, &cachedFlights); err != nil {
					log.Warn().Err(err).Msg("redis cache unmarshal failed")
				} else if len(cachedFlights) > 0 {
					mu.Lock()
					out = append(out, cachedFlights...)
					mu.Unlock()
					return
				}
			}

			flights, err := fn(ctx)
			if err != nil {
				mu.Lock()
				failed++
				mu.Unlock()

				return
			}

			dataFLight, _ := json.Marshal(flights)
			if err = s.Cache.Set(ctx, fmt.Sprintf(constant.RedisKeyProvider, name), dataFLight, constant.DurationRedisProvider); err != nil {
				zerolog.Ctx(ctx).Warn().
					Err(err).
					Str("provider", name).
					Msg("failed save provider flights to redis")
			}

			mu.Lock()
			out = append(out, flights...)
			mu.Unlock()
		}(name, fn)
	}
	wg.Wait()

	if len(out) == 0 && failed == len(providers) {
		return nil, errors.New("all providers failed")
	}

	return out, nil
}

func applyFilter(info modelFlight.Info, req *modelFlight.ReqSearch) bool {
	// origin & destination
	if !strings.EqualFold(strings.TrimSpace(info.Route.Origin), strings.TrimSpace(req.Origin)) {
		return false
	}
	if !strings.EqualFold(strings.TrimSpace(info.Route.Destination), strings.TrimSpace(req.Destination)) {
		return false
	}

	// Departure date (YYYY-MM-DD)
	if req.DepartureDate.IsZero() {
		if info.Schedule.DepartureTime != req.DepartureDate {
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

	// Arrival time window (uses only time-of-day)
	if !req.DepartureDate.IsZero() && !req.ArrivalDate.IsZero() {
		if !util.WithinTimeWindow(info.Schedule.DepartureTime, req.DepartureDate, req.ArrivalDate) {
			return false
		}
	}
	return true
}

func applyRanking(info *modelFlight.Info) {
	info.BestValueScore = float64(info.Price.Amount) / float64(info.Duration.TotalMinutes)
}
