package flight

import (
	"context"
	"fmt"
	"github.com/pkg/errors"

	"github.com/alhamsya/bookcabin/internal/core/domain/response"
)

func (s *Service) Search(ctx context.Context) (*response.Common, error) {
	garudaFlight, err := s.GarudaRepo.GetFlight(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed GetFlight")
	}

	fmt.Println(garudaFlight)
	return &response.Common{}, nil
}
