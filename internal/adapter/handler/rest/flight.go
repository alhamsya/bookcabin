package rest

import (
	"net/http"

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

	resp, err := h.Interactor.FlightService.Search(ctx.Context(), req)
	if err != nil {
		return response.New(ctx).SetHttpCode(resp.HttpCode).
			SetErr(err).SetMessage("failed searching flight").Send()
	}
	return response.New(ctx).SetHttpCode(resp.HttpCode).
		SetData(resp.Data).SetMessage("success searching flight").Send()
}
