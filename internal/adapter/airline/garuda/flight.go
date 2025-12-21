package garuda

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"

	"github.com/alhamsya/bookcabin/common/xhttp"
	"github.com/alhamsya/bookcabin/internal/core/domain/airline"
)

func (a *Airline) GetFlight(ctx context.Context) (*airline.GarudaFlightResp, error) {
	client := xhttp.NewClient(a.CfgUpstream.Garuda, http.DefaultClient)
	resp, err := client.Do(ctx, xhttp.Call{
		Upstream: "garuda",
		Endpoint: "flight",
		Method:   http.MethodGet,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed call flight")
	}

	var respGaruda airline.GarudaFlightResp
	if err = json.Unmarshal(resp.Body, &respGaruda); err != nil {
		return nil, errors.Wrap(err, "failed json unmarshal")
	}

	return &respGaruda, nil
}
