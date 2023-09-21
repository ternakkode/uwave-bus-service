package main

import "time"

type BusResponse struct {
	PlatNumber                   string
	Latitude                     float64
	Longitude                    float64
	Bearing                      float64
	CrodwLevel                   string
	DistanceInMeter              *float64
	EstimatedTimeArrivalInSecond *float64
	ETA                          *time.Time
}

type BusLineInformation struct {
	ID             string
	FullName       string
	ShortName      string
	Origin         string
	AvailableBuses []BusResponse
}

type GetBusStopResponse struct {
	ID        string
	Name      string
	Latitude  float64
	Longitude float64
	Lines     []BusLineInformation `json:",omitempty"`
}

func NewGetBusResponses(busStops []BusStopDB) []GetBusStopResponse {
	busStopResult := make([]GetBusStopResponse, 0)
	for _, busStop := range busStops {
		busStopResult = append(busStopResult, GetBusStopResponse{
			ID:        busStop.ID,
			Name:      busStop.Name,
			Latitude:  busStop.Latitude,
			Longitude: busStop.Longitude,
		})
	}

	return busStopResult
}
