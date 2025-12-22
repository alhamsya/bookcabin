package rest

import (
	"github.com/alhamsya/bookcabin/internal/core/port/service"
	"github.com/alhamsya/bookcabin/pkg/manager/config"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	Cfg        *config.Application
	App        *fiber.App
	Interactor *Interactor
}

type Interactor struct {
	FlightService port.FlightService
}

func New(this *Handler) *Handler {
	return &Handler{
		Cfg:        this.Cfg,
		App:        this.App,
		Interactor: this.Interactor,
	}
}
