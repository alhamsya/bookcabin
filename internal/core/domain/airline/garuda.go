package airline

import "time"

type GarudaFlightResp struct {
	Status  string    `json:"status"`
	Flights []Flights `json:"flights"`
}
type Departure struct {
	Airport  string    `json:"airport"`
	City     string    `json:"city"`
	Time     time.Time `json:"time"`
	Terminal string    `json:"terminal"`
}
type Arrival struct {
	Airport  string    `json:"airport"`
	City     string    `json:"city"`
	Time     time.Time `json:"time"`
	Terminal string    `json:"terminal"`
}
type Price struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
}
type Baggage struct {
	CarryOn int `json:"carry_on"`
	Checked int `json:"checked"`
}

type Segments struct {
	FlightNumber string `json:"flight_number"`
	Departure    struct {
		Airport string    `json:"airport"`
		Time    time.Time `json:"time"`
	} `json:"departure"`
	Arrival struct {
		Airport string    `json:"airport"`
		Time    time.Time `json:"time"`
	} `json:"arrival"`
	DurationMinutes int `json:"duration_minutes"`
	LayoverMinutes  int `json:"layover_minutes,omitempty"`
}
type Flights struct {
	FlightID        string     `json:"flight_id"`
	Airline         string     `json:"airline"`
	AirlineCode     string     `json:"airline_code"`
	Departure       Departure  `json:"departure"`
	Arrival         Arrival    `json:"arrival"`
	DurationMinutes int        `json:"duration_minutes"`
	Stops           int        `json:"stops"`
	Aircraft        string     `json:"aircraft"`
	Price           Price      `json:"price"`
	AvailableSeats  int        `json:"available_seats"`
	FareClass       string     `json:"fare_class"`
	Baggage         Baggage    `json:"baggage"`
	Amenities       []string   `json:"amenities,omitempty"`
	Segments        []Segments `json:"segments,omitempty"`
}
