package main

import (
	"database/sql"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

const (
	BUS_AVERAGE_SPEED_IN_M_PER_SECOND float64 = (20 * 1000) / 60 / 60
)

func GetBusStops(ctx *gin.Context) {
	busStops := make([]BusStopDB, 0)
	if err := GetDB().NewSelect().Model(&busStops).Scan(ctx); err != nil {
		ctx.Abort()
	}

	ctx.JSON(200, NewGetBusResponses(busStops))
}

func GetBusStopsByID(ctx *gin.Context) {
	var busStop BusStopDB
	if err := GetDB().NewSelect().Model(&busStop).Where("id = ?", ctx.Param("id")).Scan(ctx); err != nil {
		ctx.Abort()
	}

	busLines := make([]GetBusLineStopQueryResult, 0)
	if err := GetDB().NewRaw(
		`
		select b_line.*, b_line_stop.nearest_bus_line_path_id, l_path.sequence nearest_bus_line_path_sequence from bus_lines b_line
		inner join bus_line_stops b_line_stop on b_line.id  = b_line_stop.bus_line_id
		inner join bus_line_paths l_path on l_path.id = b_line_stop.nearest_bus_line_path_id
		where b_line_stop.bus_stop_id = ?
		`, busStop.ID,
	).Scan(ctx, &busLines); err != nil {
		ctx.Abort()
	}

	res := GetBusStopResponse{
		ID:        busStop.ID,
		Name:      busStop.Name,
		Latitude:  busStop.Latitude,
		Longitude: busStop.Longitude,
		Lines:     make([]BusLineInformation, 0, len(busLines)),
	}

	busLineIDs := make([]string, 0, len(busLines))
	busLineSequenceMap := make(map[string]int64, len(busLines))
	for _, busLine := range busLines {
		busLineSequenceMap[busLine.ID] = busLine.NearestBusLinePathSequence
		busLineIDs = append(busLineIDs, busLine.ID)
		busLine := BusLineInformation{
			ID:             busLine.ID,
			FullName:       busLine.FullName,
			ShortName:      busLine.ShortName,
			Origin:         busLine.Origin,
			AvailableBuses: make([]BusResponse, 0),
		}

		res.Lines = append(res.Lines, busLine)
	}

	busInfos := make([]GetBusInfoQueryResult, 0)
	if err := GetDB().NewRaw(
		`
		select b_info.id, b_info.current_line_id, b_info.plat_number, b_loc_hist.latitude, b_loc_hist.longitude, b_loc_hist.bearing, b_loc_hist.crowd_level from bus_informations b_info
		inner join bus_location_histories b_loc_hist on b_info.last_location_id = b_loc_hist.id
		where b_info.current_line_id IN (?)`, bun.In(busLineIDs),
	).Scan(ctx, &busInfos); err != nil {
		ctx.Abort()
	}

	busInsideLineMap := make(map[string][]BusResponse)
	for _, busInfo := range busInfos {
		busInformation := BusResponse{
			PlatNumber: busInfo.PlatNumber,
			Latitude:   busInfo.Latitude,
			Longitude:  busInfo.Longitude,
			Bearing:    busInfo.Bearing,
			CrodwLevel: busInfo.CrowdLevel,
		}

		// TODO: Calculate the distance once the bus has already gone beyond the bus stop.
		var busFinalTotalDistance float64
		err := GetDB().NewRaw(
			`With bus_to_closest_path as (
				select blp.bus_line_id , blp."sequence", haversine_distance(?, ?, blp.latitude, blp.longitude) * 1000 as distance
				from bus_line_paths blp  
				where bus_line_id = ?
				order by distance asc
				limit 1
			)
			SELECT (SUM(blp.distance_to_next_path)) + bus_to_closest_path.distance
			FROM bus_line_paths blp 
			inner join bus_to_closest_path on bus_to_closest_path.bus_line_id = blp.bus_line_id 
			where blp."sequence" >= bus_to_closest_path.sequence 
			and blp."sequence" <= ? 
			and blp.bus_line_id = ?
			group  by (blp.bus_line_id, bus_to_closest_path.distance)`,
			busInformation.Latitude, busInformation.Longitude,
			busInfo.CurrentLineID,
			busLineSequenceMap[busInfo.CurrentLineID],
			busInfo.CurrentLineID,
		).Scan(ctx, &busFinalTotalDistance)

		var noDataFound bool // no data found possibly because bus already pass the bus stop since we didnt handle it.
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				noDataFound = true
			} else {
				ctx.Abort()
				return
			}
		}

		if !noDataFound {
			busInformation.DistanceInMeter = &busFinalTotalDistance

			etaInSecond := busFinalTotalDistance / BUS_AVERAGE_SPEED_IN_M_PER_SECOND
			etaInTime := time.Now().Add(time.Second * time.Duration(etaInSecond)).UTC()
			busInformation.EstimatedTimeArrivalInSecond = &etaInSecond
			busInformation.ETA = &etaInTime
		}

		lineBuses, ok := busInsideLineMap[busInfo.CurrentLineID]
		if !ok {
			lineBuses = make([]BusResponse, 0)
		}

		busInsideLineMap[busInfo.CurrentLineID] = append(lineBuses, busInformation)
	}

	for i := range res.Lines {
		res.Lines[i].AvailableBuses = busInsideLineMap[res.Lines[i].ID]
	}

	ctx.JSON(200, res)
}
