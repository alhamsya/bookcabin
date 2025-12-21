package modelFlight

import "time"

type Provider string

const (
	ProviderGaruda  Provider = "GARUDA"
	ProviderLion    Provider = "LION"
	ProviderBatik   Provider = "BATIK"
	ProviderAirAsia Provider = "AIRASIA"
)

type Info struct {
	ID       string
	Provider Provider

	Airline  Airline
	Route    Route
	Schedule Schedule

	Duration Duration
	Stops    int

	Price Price

	SeatsAvailable int

	BestValueScore float64
}

type Airline struct {
	Name string
	Code string
}

type Route struct {
	Origin      string
	Destination string
}

type Schedule struct {
	DepartureTime time.Time
	ArrivalTime   time.Time

	DepartureTs int64
	ArrivalTs   int64
}

type Duration struct {
	TotalMinutes int
	Formatted    string // "1h 45m"
}

type Price struct {
	Amount   int
	Currency string
}

type Baggage struct {
	CarryOn string
	Checked string
}
