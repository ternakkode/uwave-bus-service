package main

import "fmt"

type Coord struct {
	Latitude  float64
	Longitude float64
}

func ConvertCoordListToStringLongLat(coords ...Coord) string {
	results := ""
	for index, coord := range coords {
		results += fmt.Sprintf("%v,%v", coord.Longitude, coord.Latitude)

		isLastIndex := len(coords)-1 == index
		if !isLastIndex {
			results += ";"
		}
	}

	return results
}
