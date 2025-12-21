package airasia

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/alhamsya/bookcabin/common/xhttp"
	"github.com/alhamsya/bookcabin/internal/core/domain/constant"
	"github.com/alhamsya/bookcabin/lib/util"
	"github.com/pkg/errors"

	modelAirAsia "github.com/alhamsya/bookcabin/internal/core/domain/airline/airasia"
	modelFlight "github.com/alhamsya/bookcabin/internal/core/domain/flight"
)

func (a *Airline) GetFlight(ctx context.Context) ([]modelFlight.Info, error) {
	client := xhttp.NewClient(a.CfgUpstream.AirAsia, http.DefaultClient)
	resp, err := client.Do(ctx, xhttp.Call{
		Upstream: "airasia", Endpoint: "flight", Method: http.MethodGet,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed call flight airasia")
	}
	var respAirAsia modelAirAsia.FlightResp
	if err = json.Unmarshal(resp.Body, &respAirAsia); err != nil {
		return nil, errors.Wrap(err, "failed json unmarshal model airasia")
	}
	out := make([]modelFlight.Info, 0, len(respAirAsia.Flights))
	for _, data := range respAirAsia.Flights {
		dep, err := time.Parse(time.RFC3339, data.DepartTime)
		if err != nil {
			return nil, errors.Wrapf(err, "airasia invalid depart_time for %s", data.FlightCode)
		}
		arr, err := time.Parse(time.RFC3339, data.ArriveTime)
		if err != nil {
			return nil, errors.Wrapf(err, "airasia invalid arrive_time for %s", data.FlightCode)
		}
		totalMin := int(arr.Sub(dep).Minutes())
		if totalMin < 0 {
			return nil, errors.Errorf("airasia invalid schedule: arrival before departure (%s)", data.FlightCode)
		}
		stops := 0
		if !data.DirectFlight {
			stops = len(data.Stops)
			if stops == 0 {
				stops = 1
			}
		}

		airlineCode := ""
		if len(data.FlightCode) >= 2 {
			airlineCode = data.FlightCode[:2]
		}

		out = append(out, modelFlight.Info{
			ID:       fmt.Sprintf("%s_AIRASIA", data.FlightCode),
			Provider: modelFlight.ProviderAirAsia,
			Airline: modelFlight.Airline{
				Name: data.Airline,
				Code: airlineCode,
			},
			Route: modelFlight.Route{
				Origin:      data.FromAirport,
				Destination: data.ToAirport,
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
				Amount:   data.PriceIdr,
				Currency: constant.CurrencyIDR,
			},
			SeatsAvailable: data.Seats,
		})
	}
	return out, nil
}
