package modelRequest

import (
	"github.com/go-playground/validator/v10"
)

type ReqSearchFlight struct {
	Origin        string  `json:"origin" validate:"required,min=3,max=3"`
	Destination   string  `json:"destination" validate:"required,min=3,max=3"`
	DepartureDate string  `json:"departureDate" validate:"required"`
	ArrivalDate   string  `json:"arrivalDate"`
	Passengers    int     `json:"passengers" validate:"required,gte=1"`
	CabinClass    string  `json:"cabinClass" validate:"required"`
	Sort          Sort    `json:"sort"`
	Filters       Filters `json:"filters"`
}

type Sort struct {
	Key   string `json:"key"`
	Order string `json:"order"`
}
type DepartureTime struct {
	From string `json:"from"`
	To   string `json:"to"`
}
type ArrivalTime struct {
	From string `json:"from"`
	To   string `json:"to"`
}
type Filters struct {
	MinPrice           int      `json:"minPrice"`
	MaxPrice           int      `json:"maxPrice"`
	Stops              []int    `json:"stops"`
	Airlines           []string `json:"airlines"`
	MaxDurationMinutes int      `json:"maxDurationMinutes"`
}

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

func (req ReqSearchFlight) ValidateStruct() []*ErrorResponse {
	var errors []*ErrorResponse
	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}
