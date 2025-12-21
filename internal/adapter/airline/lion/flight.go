package lion

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/alhamsya/bookcabin/common/xhttp"
	"github.com/alhamsya/bookcabin/lib/util"
	"github.com/pkg/errors"

	modelLion "github.com/alhamsya/bookcabin/internal/core/domain/airline/lion"
	modelFlight "github.com/alhamsya/bookcabin/internal/core/domain/flight"
)

func (a *Airline) GetFlight(ctx context.Context) ([]modelFlight.Info, error) {
	client := xhttp.NewClient(a.CfgUpstream.Lion, http.DefaultClient)
	resp, err := client.Do(ctx, xhttp.Call{
		Upstream: "lion",
		Endpoint: "flight",
		Method:   http.MethodGet,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed call flight lion")
	}

	var respLion modelLion.FlightResp
	if err = json.Unmarshal(resp.Body, &respLion); err != nil {
		return nil, errors.Wrap(err, "failed json unmarshal model airasia")
	}

	out := make([]modelFlight.Info, 0, len(respLion.Data.AvailableFlights))
	for _, data := range respLion.Data.AvailableFlights {
		locDep, err := time.LoadLocation(data.Schedule.DepartureTimezone)
		if err != nil {
			return nil, errors.Wrapf(err, "lion invalid departure_timezone (%s)", data.ID)
		}
		locArr, err := time.LoadLocation(data.Schedule.ArrivalTimezone)
		if err != nil {
			return nil, errors.Wrapf(err, "lion invalid arrival_timezone (%s)", data.ID)
		}

		dep, err := time.ParseInLocation(time.RFC3339, data.Schedule.Departure, locDep)
		if err != nil {
			return nil, errors.Wrapf(err, "lion invalid departure time (%s)", data.ID)
		}
		arr, err := time.ParseInLocation(time.RFC3339, data.Schedule.Arrival, locArr)
		if err != nil {
			return nil, errors.Wrapf(err, "lion invalid arrival time (%s)", data.ID)
		}

		totalMin := data.FlightTime
		if totalMin <= 0 {
			totalMin = int(arr.Sub(dep).Minutes())
		}
		if totalMin < 0 {
			return nil, errors.Errorf("lion invalid schedule: arrival before departure (%s)", data.ID)
		}

		stops := 0
		if !data.IsDirect {
			stops = 1
		}

		airlineCode := ""
		if len(data.Carrier.Iata) >= 2 {
			airlineCode = data.Carrier.Iata[:2]
		}
		out = append(out, modelFlight.Info{
			ID:       fmt.Sprintf("%s_LION", data.ID),
			Provider: modelFlight.ProviderLion,
			Airline: modelFlight.Airline{
				Name: data.Carrier.Name,
				Code: airlineCode,
			},
			Route: modelFlight.Route{
				Origin:      data.Route.From.Code,
				Destination: data.Route.To.Code,
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
			Stops: stops,
			Price: modelFlight.Price{
				Amount:   data.Pricing.Total,
				Currency: data.Pricing.Currency,
			},
			SeatsAvailable: data.SeatsLeft,
		})
	}

	return out, nil
}
