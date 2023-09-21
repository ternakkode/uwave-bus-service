package main

type UwaveBaseResponse[T any] struct {
	Status  int64 `json:"status"`
	Payload T     `json:"payload"`
}

type UwaveBusLineDetail struct {
	ID        string             `json:"id"`
	FullName  string             `json:"fullName"`
	Origin    string             `json:"origin"`
	ShortName string             `json:"shortName"`
	BusStops  []UwaveBusStop     `json:"busStops"`
	Paths     []UwaveBusLinePath `json:"path"`
}

type UwaveBusLinePath []float64

func (u UwaveBusLinePath) GetLat() float64 {
	return u[0]
}

func (u UwaveBusLinePath) GetLong() float64 {
	return u[1]
}

type UwaveBusStop struct {
	ID        string  `json:"id"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
	Name      string  `json:"name"`
}

type UwaveBusPosition struct {
	Bearing      float64 `json:"bearing"`
	CrowdLevel   string  `json:"crowdLevel"`
	Latitude     float64 `json:"lat"`
	Longitude    float64 `json:"lng"`
	VehiclePlate string  `json:"vehiclePlate"`
}
