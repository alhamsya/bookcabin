package modelLion

type FlightResp struct {
	Success bool `json:"success"`
	Data    Data `json:"data"`
}
type Carrier struct {
	Name string `json:"name"`
	Iata string `json:"iata"`
}
type From struct {
	Code string `json:"code"`
	Name string `json:"name"`
	City string `json:"city"`
}
type To struct {
	Code string `json:"code"`
	Name string `json:"name"`
	City string `json:"city"`
}
type Route struct {
	From From `json:"from"`
	To   To   `json:"to"`
}
type Schedule struct {
	Departure         string `json:"departure"`
	DepartureTimezone string `json:"departure_timezone"`
	Arrival           string `json:"arrival"`
	ArrivalTimezone   string `json:"arrival_timezone"`
}
type Pricing struct {
	Total    int    `json:"total"`
	Currency string `json:"currency"`
	FareType string `json:"fare_type"`
}
type BaggageAllowance struct {
	Cabin string `json:"cabin"`
	Hold  string `json:"hold"`
}
type Services struct {
	WifiAvailable    bool             `json:"wifi_available"`
	MealsIncluded    bool             `json:"meals_included"`
	BaggageAllowance BaggageAllowance `json:"baggage_allowance"`
}
type Layovers struct {
	Airport         string `json:"airport"`
	DurationMinutes int    `json:"duration_minutes"`
}
type AvailableFlights struct {
	ID         string     `json:"id"`
	Carrier    Carrier    `json:"carrier"`
	Route      Route      `json:"route"`
	Schedule   Schedule   `json:"schedule"`
	FlightTime int        `json:"flight_time"`
	IsDirect   bool       `json:"is_direct"`
	Pricing    Pricing    `json:"pricing"`
	SeatsLeft  int        `json:"seats_left"`
	PlaneType  string     `json:"plane_type"`
	Services   Services   `json:"services"`
	StopCount  int        `json:"stop_count,omitempty"`
	Layovers   []Layovers `json:"layovers,omitempty"`
}
type Data struct {
	AvailableFlights []AvailableFlights `json:"available_flights"`
}
