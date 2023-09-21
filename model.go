package main

import (
	"github.com/uptrace/bun"
)

type BusLineDB struct {
	bun.BaseModel `bun:"bus_lines"`

	ID         string `bun:",pk" json:"id"`
	ExternalID string
	FullName   string
	ShortName  string
	Origin     string
}

type BusStopDB struct {
	bun.BaseModel `bun:"bus_stops"`

	ID         string `bun:",pk" json:"id"`
	ExternalID string
	Name       string
	Latitude   float64
	Longitude  float64
}

type BusLinePathDB struct {
	bun.BaseModel `bun:"bus_line_paths"`

	ID                 string `bun:",pk" json:"id"`
	Sequence           int64
	BusLineID          string
	Latitude           float64
	Longitude          float64
	DistanceToNextPath float64
}

type BusLineStopDB struct {
	bun.BaseModel `bun:"bus_line_stops"`

	ID                   string `bun:",pk" json:"id"`
	BusLineID            string
	NearestBusLinePathID string
	BusStopID            string
}

type BusInformationDB struct {
	bun.BaseModel `bun:"bus_informations"`

	ID             string `bun:",pk" json:"id"`
	PlatNumber     string
	CurrentLineID  string
	LastLocationID string
}

type BusLocationHistoryDB struct {
	bun.BaseModel `bun:"bus_location_histories"`

	ID         string `bun:",pk" json:"id"`
	BusID      string
	Latitude   float64
	Longitude  float64
	Bearing    float64
	CrowdLevel string
}

type GetBusLineStopQueryResult struct {
	NearestBusLinePathID       string
	NearestBusLinePathSequence int64
	ID                         string
	ExternalID                 string
	FullName                   string
	ShortName                  string
	Origin                     string
}

type GetBusInfoQueryResult struct {
	ID            string
	PlatNumber    string
	CurrentLineID string
	Latitude      float64
	Longitude     float64
	Bearing       float64
	CrowdLevel    string
}
