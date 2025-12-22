package rest

import (
	"net/http"
	"time"

	"github.com/alhamsya/bookcabin/internal/core/domain/constant"
	"github.com/alhamsya/bookcabin/internal/core/domain/request"
	"github.com/alhamsya/bookcabin/pkg/manager/response"
	"github.com/gofiber/fiber/v2"

	modelFlight "github.com/alhamsya/bookcabin/internal/core/domain/flight"
)

// SearchFlight handle POST /v1/flights/search
func (h *Handler) SearchFlight(ctx *fiber.Ctx) error {
	req := new(modelRequest.ReqSearchFlight)
	err := ctx.BodyParser(req)
	if err != nil {
		return response.New(ctx).SetHttpCode(http.StatusBadRequest).
			SetErr(err).SetMessage("please check request body").Send()
	}

	var departureDate, arrivalDate time.Time
	departureDate, err = time.Parse(constant.DateOnly, req.DepartureDate)
	if err != nil {
		return response.New(ctx).SetHttpCode(http.StatusBadRequest).
			SetErr(err).SetMessage("please check request departureDate").Send()
	}
	if req.ArrivalDate != "" {
		arrivalDate, err = time.Parse(constant.DateOnly, req.ArrivalDate)
		if err != nil {
			return response.New(ctx).SetHttpCode(http.StatusBadRequest).
				SetErr(err).SetMessage("please check request arrivalDate").Send()
		}
	}

	errValidate := req.ValidateStruct()
	if errValidate != nil {
		return response.New(ctx).SetHttpCode(http.StatusBadRequest).
			SetErr(errValidate).
			SetMessage("please check your request").Send()
	}

	resp, err := h.Interactor.FlightService.Search(ctx.Context(), &modelFlight.ReqSearch{
		Origin:        req.Origin,
		Destination:   req.Destination,
		DepartureDate: departureDate,
		ArrivalDate:   arrivalDate,
		Passengers:    req.Passengers,
		CabinClass:    req.CabinClass,
		Sort: modelFlight.Sort{
			Key:   req.Sort.Key,
			Order: req.Sort.Order,
		},
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
		SetData(resp).SetMessage("success searching flight").Send()
}
