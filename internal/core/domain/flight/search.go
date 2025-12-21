package modelFlight

import "time"

type ReqSearch struct {
	Origin        string
	Destination   string
	DepartureDate time.Time
	ArrivalDate   time.Time
	Passengers    int
	CabinClass    string
	Sort          Sort
	Filters       Filters
}

type Sort struct {
	Key   string
	Order string
}
type DepartureTime struct {
	From time.Time
	To   time.Time
}
type ArrivalTime struct {
	From time.Time
	To   time.Time
}
type Filters struct {
	MinPrice           int
	MaxPrice           int
	Stops              []int
	Airlines           []string
	MaxDurationMinutes int
}
