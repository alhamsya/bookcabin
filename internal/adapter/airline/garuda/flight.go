package garuda

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/alhamsya/bookcabin/common/xhttp"
	"github.com/alhamsya/bookcabin/lib/util"
	"github.com/pkg/errors"

	modelGaruda "github.com/alhamsya/bookcabin/internal/core/domain/airline/garuda"
	modelFlight "github.com/alhamsya/bookcabin/internal/core/domain/flight"
)

func (a *Airline) GetFlight(ctx context.Context) ([]modelFlight.Info, error) {
	client := xhttp.NewClient(a.CfgUpstream.Garuda, http.DefaultClient)
	resp, err := client.Do(ctx, xhttp.Call{
		Upstream: "garuda",
		Endpoint: "flight",
		Method:   http.MethodGet,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed call flight garuda")
	}

	var respGaruda modelGaruda.FlightResp
	if err = json.Unmarshal(resp.Body, &respGaruda); err != nil {
		return nil, errors.Wrap(err, "failed json unmarshal model garuda")
	}

	out := make([]modelFlight.Info, 0, len(respGaruda.Flights))
	for _, data := range respGaruda.Flights {
		dep, err := time.Parse(time.RFC3339, data.Departure.Time)
		if err != nil {
			return nil, errors.Wrapf(err, "garuda invalid departure.time (%s)", data.FlightID)
		}
		arr, err := time.Parse(time.RFC3339, data.Arrival.Time)
		if err != nil {
			return nil, errors.Wrapf(err, "garuda invalid arrival.time (%s)", data.FlightID)
		}

		totalMin := data.DurationMinutes
		if totalMin <= 0 {
			totalMin = int(arr.Sub(dep).Minutes())
		}
		if totalMin < 0 {
			return nil, errors.Errorf("garuda invalid schedule: arrival before departure (%s)", data.FlightID)
		}

		airlineCode := ""
		if len(data.AirlineCode) >= 2 {
			airlineCode = data.AirlineCode[:2]
		}

		out = append(out, modelFlight.Info{
			ID:       fmt.Sprintf("%s_GARUDA", data.FlightID),
			Provider: modelFlight.ProviderGaruda,
			Airline: modelFlight.Airline{
				Name: data.Airline,
				Code: airlineCode,
			},
			Route: modelFlight.Route{
				Origin:      data.Departure.Airport,
				Destination: data.Arrival.Airport,
			},
			Schedule: modelFlight.Schedule{
				DepartureTime: dep,
				ArrivalTime:   arr,
				DepartureTs:   dep.Unix(),
				ArrivalTs:     arr.Unix(),
			},
			Duration: modelFlight.Duration{
				TotalMinutes: totalMin,
				Formatted:    util.FormatMinutes(totalMin),
			},
			Stops: data.Stops,
			Price: modelFlight.Price{
				Amount:   data.Price.Amount,
				Currency: data.Price.Currency,
			},
			SeatsAvailable: data.AvailableSeats,
		})
	}
	return out, nil
}
