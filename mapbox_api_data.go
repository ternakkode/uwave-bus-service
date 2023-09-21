package main

type MapBoxDistanceAPIMinifiedResponse struct {
	Routes []MapBoxRouteMinified `json:"routes"`
}

type MapBoxRouteMinified struct {
	Distance float64             `json:"distance"`
	Legs     []MapboxLegMinified `json:"legs"`
}

type MapboxLegMinified struct {
	Distance float64 `json:"distance"`
}
