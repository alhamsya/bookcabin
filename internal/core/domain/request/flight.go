package modelRequest

import (
	"github.com/alhamsya/bookcabin/pkg/util"
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
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value,omitempty"`
}

type ValidationError struct {
	Errors []ErrorResponse `json:"errors"`
}

func (e ValidationError) Error() string {
	return "validation failed"
}

func (req ReqSearchFlight) ValidateStruct() error {
	v := validator.New()

	err := v.Struct(req)
	if err == nil {
		return nil
	}

	ve, ok := err.(validator.ValidationErrors)
	if !ok {
		return ValidationError{Errors: []ErrorResponse{{
			Field: "request",
			Tag:   "invalid",
		}}}
	}

	out := make([]ErrorResponse, 0, len(ve))
	for _, fe := range ve {
		out = append(out, ErrorResponse{
			Field: util.GetFieldName(req, fe.StructField(), "json"),
			Tag:   fe.Tag(),
			Value: fe.Param(),
		})
	}

	return ValidationError{Errors: out}
}
