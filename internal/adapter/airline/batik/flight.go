package batik

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/alhamsya/bookcabin/common/xhttp"
	"github.com/alhamsya/bookcabin/lib/util"
	"github.com/pkg/errors"

	modelBatik "github.com/alhamsya/bookcabin/internal/core/domain/airline/batik"
	modelFlight "github.com/alhamsya/bookcabin/internal/core/domain/flight"
)

func (a *Airline) GetFlight(ctx context.Context) ([]modelFlight.Info, error) {
	client := xhttp.NewClient(a.CfgUpstream.Batik, http.DefaultClient)
	resp, err := client.Do(ctx, xhttp.Call{
		Upstream: "batik",
		Endpoint: "flight",
		Method:   http.MethodGet,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed call flight batik")
	}

	var respBatik modelBatik.FlightResp
	if err = json.Unmarshal(resp.Body, &respBatik); err != nil {
		return nil, errors.Wrap(err, "failed json unmarshal model batik")
	}

	out := make([]modelFlight.Info, 0, len(respBatik.Results))
	for _, data := range respBatik.Results {
		dep, err := time.Parse(time.RFC3339, data.DepartureDateTime)
		if err != nil {
			return nil, errors.Wrapf(err, "batik invalid departureDateTime (%s)", data.FlightNumber)
		}
		arr, err := time.Parse(time.RFC3339, data.ArrivalDateTime)
		if err != nil {
			return nil, errors.Wrapf(err, "batik invalid arrivalDateTime (%s)", data.FlightNumber)
		}

		totalMin := int(arr.Sub(dep).Minutes())
		if totalMin < 0 {
			return nil, errors.Errorf("batik invalid schedule: arrival before departure (%s)", data.FlightNumber)
		}

		airlineCode := ""
		if len(data.AirlineIATA) >= 2 {
			airlineCode = data.AirlineIATA[:2]
		}

		out = append(out, modelFlight.Info{
			ID:       fmt.Sprintf("%s_BATIK", data.FlightNumber),
			Provider: modelFlight.ProviderBatik,
			Airline: modelFlight.Airline{
				Name: data.AirlineName,
				Code: airlineCode,
			},
			Route: modelFlight.Route{
				Origin:      data.Origin,
				Destination: data.Destination,
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
			Stops: data.NumberOfStops,
			Price: modelFlight.Price{
				Amount:   data.Fare.TotalPrice,
				Currency: data.Fare.CurrencyCode,
			},
			SeatsAvailable: data.SeatsAvailable,
		})
	}

	return out, nil
}
