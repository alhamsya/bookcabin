package modelAirAsia

type FlightResp struct {
	Status  string    `json:"status"`
	Flights []Flights `json:"flights"`
}
type Stops struct {
	Airport         string `json:"airport"`
	WaitTimeMinutes int    `json:"wait_time_minutes"`
}
type Flights struct {
	FlightCode    string  `json:"flight_code"`
	Airline       string  `json:"airline"`
	FromAirport   string  `json:"from_airport"`
	ToAirport     string  `json:"to_airport"`
	DepartTime    string  `json:"depart_time"`
	ArriveTime    string  `json:"arrive_time"`
	DurationHours float64 `json:"duration_hours"`
	DirectFlight  bool    `json:"direct_flight"`
	PriceIdr      int     `json:"price_idr"`
	Seats         int     `json:"seats"`
	CabinClass    string  `json:"cabin_class"`
	BaggageNote   string  `json:"baggage_note"`
	Stops         []Stops `json:"stops,omitempty"`
}
