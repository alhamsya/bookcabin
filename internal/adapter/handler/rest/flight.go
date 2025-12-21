package rest

import (
	"github.com/alhamsya/bookcabin/internal/core/domain/constant"
	modelFlight "github.com/alhamsya/bookcabin/internal/core/domain/flight"
	"net/http"
	"time"

	"github.com/alhamsya/bookcabin/internal/core/domain/request"
	"github.com/alhamsya/bookcabin/lib/manager/response"
	"github.com/gofiber/fiber/v2"
)

// GetSearchFlight handle POST /v1/flights/search
func (h *Handler) GetSearchFlight(ctx *fiber.Ctx) error {
	req := new(modelRequest.ReqSearchFlight)
	err := ctx.BodyParser(req)
	if err != nil {
		return response.New(ctx).SetHttpCode(http.StatusBadRequest).
			SetErr(err).SetData("please check request body").Send()
	}

	var departureDate, arrivalDate time.Time
	departureDate, err = time.Parse(constant.DateOnly, req.DepartureDate)
	if err != nil {
		return response.New(ctx).SetHttpCode(http.StatusBadRequest).
			SetErr(err).SetData("please check request departureDate").Send()
	}
	if req.ArrivalDate != "" {
		arrivalDate, err = time.Parse(constant.DateOnly, req.ArrivalDate)
		if err != nil {
			return response.New(ctx).SetHttpCode(http.StatusBadRequest).
				SetErr(err).SetData("please check request arrivalDate").Send()
		}
	}

	resp, err := h.Interactor.FlightService.Search(ctx.Context(), &modelFlight.ReqSearch{
		Origin:        req.Origin,
		Destination:   req.Destination,
		DepartureDate: departureDate,
		ArrivalDate:   arrivalDate,
		Passengers:    req.Passengers,
		CabinClass:    req.CabinClass,
		Sort:          modelFlight.Sort{},
		Filters: modelFlight.Filters{
			MinPrice:           req.Filters.MinPrice,
			MaxPrice:           req.Filters.MaxPrice,
			Stops:              req.Filters.Stops,
			Airlines:           req.Filters.Airlines,
			MaxDurationMinutes: req.Filters.MaxDurationMinutes,
		},
	})
	if err != nil {
		return response.New(ctx).SetHttpCode(resp.HttpCode).
			SetErr(err).SetMessage("failed searching flight").Send()
	}
	return response.New(ctx).SetHttpCode(resp.HttpCode).
		SetData(resp.Data).SetMessage("success searching flight").Send()
}
