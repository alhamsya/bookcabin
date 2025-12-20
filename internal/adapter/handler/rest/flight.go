package rest

import (
	"github.com/alhamsya/bookcabin/lib/manager/response"
	"github.com/gofiber/fiber/v2"
)

// SearchFlights handle POST /v1/flights/search
func (h *Handler) SearchFlights(ctx *fiber.Ctx) error {
	resp, err := h.Interactor.FlightService.Search(ctx.Context())
	if err != nil {
		return response.New(ctx).
			SetErr(err).
			SetHttpCode(resp.HttpCode).
			SetMessage("failed searching flight").
			Send()
	}
	return response.New(ctx).
		SetData(resp.Data).
		SetHttpCode(resp.HttpCode).
		SetMessage("success searching flight").
		Send()
}
