package flight

import (
	"context"
	"github.com/alhamsya/bookcabin/internal/core/domain/request"
	"sync"

	"github.com/alhamsya/bookcabin/internal/core/domain/response"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	modelFlight "github.com/alhamsya/bookcabin/internal/core/domain/flight"
)

func (s *Service) Search(ctx context.Context, param *modelRequest.ReqSearchFlight) (modelResponse.Common, error) {
	// 1. check cache in L1 in-memory
	// 2. call providers airline using asynchronous
	// 		2.1. check data in redis
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

		out = append(out, info)
	}

	// 4. calculate best value score (ranking)
	// 5. sort result
	// 6. set metadata
	// 7. save to cache in L1 in-memory

	return modelResponse.Common{}, nil
}

func (s *Service) callProviders(ctx context.Context) ([]modelFlight.Info, error) {
	var (
		wg   sync.WaitGroup
		mu   sync.Mutex
		out  []modelFlight.Info
		errs = make(map[string]error)
	)

	run := func(name string, fn func(context.Context) ([]modelFlight.Info, error)) {
		wg.Add(1)
		go func() {
			defer wg.Done()

			flights, err := fn(ctx)

			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				errs[name] = err
				return
			}
			out = append(out, flights...)
		}()
	}

	run("airasia", s.AirAsiaRepo.GetFlight)
	run("batik", s.BatikRepo.GetFlight)
	run("garuda", s.GarudaRepo.GetFlight)
	run("lion", s.LionRepo.GetFlight)

	wg.Wait()
	zerolog.Ctx(ctx).Warn().
		Interface("errors", errs).
		Msg("test")
	return out, nil
}
