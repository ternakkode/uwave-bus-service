package main

import (
	"context"
	"database/sql"
	"log"
	"math"

	"github.com/segmentio/ksuid"
)

func UpdateBusLocation(ctx context.Context) error {
	busLinesInfo := make([]BusLineDB, 0)
	if err := GetDB().
		NewSelect().
		Model(&busLinesInfo).
		Column("id", "external_id").
		Scan(ctx); err != nil {
		return err
	}

	busLocations := make([]BusLocationHistoryDB, 0)
	busInformations := make([]BusInformationDB, 0)
	for _, busLine := range busLinesInfo {
		busesInsideLane, err := GetBusPositionUwave(ctx, busLine.ExternalID)
		if err != nil {
			log.Println(err)
		}

		for _, busInsideLane := range busesInsideLane {
			busInformationID := ksuid.New().String()
			busLocationID := ksuid.New().String()

			busInformations = append(busInformations, BusInformationDB{
				ID:             busInformationID,
				PlatNumber:     busInsideLane.VehiclePlate,
				CurrentLineID:  busLine.ID,
				LastLocationID: busLocationID,
			})

			busLocations = append(busLocations, BusLocationHistoryDB{
				ID:         busLocationID,
				BusID:      busInformationID,
				Latitude:   busInsideLane.Latitude,
				Longitude:  busInsideLane.Longitude,
				Bearing:    busInsideLane.Bearing,
				CrowdLevel: busInsideLane.CrowdLevel,
			})

		}
	}

	if _, err := GetDB().
		NewInsert().
		Model(&busInformations).
		On("CONFLICT (plat_number) DO UPDATE").
		Exec(ctx); err != nil {
		return err
	}

	if _, err := GetDB().
		NewInsert().
		Model(&busLocations).
		Exec(ctx); err != nil {
		return err
	}

	return nil
}

// PopulateBusLineInformation function serves to retrieve bus line and bus stop data from the uWave API and save it in a database.
// Additionally, it computes the distances between points along a bus path using the MapBox Distance API.
// It also determines the nearest bus path for a given bus stop within a bus line, employing the Haversine formula.
func PopulateBusLineInformation(ctx context.Context) error {
	latestBuslineInformations, err := GetBusLinesUwave(ctx)
	if err != nil {
		return err
	}

	busLinesToDB := make([]*BusLineDB, 0, len(latestBuslineInformations))
	busStopToDB := make(map[string]*BusStopDB, 0)
	busLineStopList := make(map[string][]string)
	busPathToDB := make(map[string][]*BusLinePathDB, 0)
	for _, busLine := range latestBuslineInformations {
		busLinesToDB = append(busLinesToDB, &BusLineDB{
			ID:         ksuid.New().String(),
			ExternalID: busLine.ID,
			FullName:   busLine.FullName,
			ShortName:  busLine.ShortName,
			Origin:     busLine.Origin,
		})

		var lastSearchedBusPathIndex int
		distanceToNextBusPath := make([]float64, 0, len(busLine.Paths))
		for {
			busCoords := make([]Coord, 0, 25)
			for i := lastSearchedBusPathIndex; len(busCoords) < 25 && i < len(busLine.Paths); i++ {
				busCoords = append(busCoords, Coord{
					Latitude:  busLine.Paths[i].GetLat(),
					Longitude: busLine.Paths[i].GetLong(),
				})
			}

			mapBoxDistanceResult, err := MapBoxGetDirection(ctx, busCoords...)
			if err != nil {
				return err
			}

			for _, leg := range mapBoxDistanceResult.Routes[0].Legs {
				distanceToNextBusPath = append(distanceToNextBusPath, leg.Distance)
			}

			lastSearchedBusPathIndex += 24
			if lastSearchedBusPathIndex > len(busLine.Paths) {
				break
			}
		}

		for i := range busLine.Paths {
			var distanceToNexPath float64
			if i == len(busLine.Paths)-1 {
				distanceToNexPath = distanceToNextBusPath[i-1]
			} else {
				distanceToNexPath = distanceToNextBusPath[i]
			}

			busLinePathDB := BusLinePathDB{
				ID:                 ksuid.New().String(),
				Sequence:           int64(i + 1),
				Latitude:           busLine.Paths[i].GetLat(),
				Longitude:          busLine.Paths[i].GetLong(),
				DistanceToNextPath: distanceToNexPath,
			}

			busLinePaths, ok := busPathToDB[busLine.ID]
			if !ok {
				busLinePaths = make([]*BusLinePathDB, 0)
			}

			busLinePaths = append(busLinePaths, &busLinePathDB)
			busPathToDB[busLine.ID] = busLinePaths
		}

		for _, busStop := range busLine.BusStops {
			busLineStops, ok := busLineStopList[busLine.ID]
			if !ok {
				busLineStops = make([]string, 0)
			}

			busLineStops = append(busLineStops, busStop.ID)
			busLineStopList[busLine.ID] = busLineStops

			if _, ok := busStopToDB[busStop.ID]; ok {
				continue
			}

			busStopToDB[busStop.ID] = &BusStopDB{
				ID:         ksuid.New().String(),
				Name:       busStop.Name,
				ExternalID: busStop.ID,
				Latitude:   busStop.Latitude,
				Longitude:  busStop.Longitude,
			}
		}
	}

	tx, err := GetDB().BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	if _, err := tx.NewInsert().Model(&busLinesToDB).On("CONFLICT (external_id) DO UPDATE").Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	busLineExternalIDMap := make(map[string]string, len(busLinesToDB))
	for _, busLine := range busLinesToDB {
		busLineExternalIDMap[busLine.ExternalID] = busLine.ID
	}

	busStops := make([]*BusStopDB, 0, len(busStopToDB))
	for _, busStopDB := range busStopToDB {
		busStops = append(busStops, busStopDB)
	}

	if _, err := tx.NewInsert().Model(&busStops).On("CONFLICT (external_id) DO UPDATE").Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	finalBusPath := make([]*BusLinePathDB, 0)
	for externalBusStopID, busPathList := range busPathToDB {
		busLineID := busLineExternalIDMap[externalBusStopID]
		for _, busPath := range busPathList {
			busPath.BusLineID = busLineID
			finalBusPath = append(finalBusPath, busPath)
		}
	}

	if _, err := tx.NewTruncateTable().Model((*BusLinePathDB)(nil)).Cascade().Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.NewTruncateTable().Model((*BusLineStopDB)(nil)).Cascade().Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.NewInsert().Model(&finalBusPath).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	busLineStops := make([]BusLineStopDB, 0)
	for busLineExternalID, busStopExternalIDs := range busLineStopList {
		busLineIDFromDB := busLineExternalIDMap[busLineExternalID]
		for _, busStopExternalID := range busStopExternalIDs {
			busStopInfo := busStopToDB[busStopExternalID]

			var nearesetBusPathID string
			var nearestBusPathDistance float64 = math.MaxFloat64

			for _, busPathInLine := range busPathToDB[busLineExternalID] {
				distanceInKM := Distance(
					Coord{
						Latitude:  busStopInfo.Latitude,
						Longitude: busStopInfo.Longitude,
					},
					Coord{
						Latitude:  busPathInLine.Latitude,
						Longitude: busPathInLine.Longitude,
					},
				)

				if nearestBusPathDistance > distanceInKM {
					nearesetBusPathID = busPathInLine.ID
					nearestBusPathDistance = distanceInKM
				}
			}

			busLineStops = append(busLineStops, BusLineStopDB{
				ID:                   ksuid.New().String(),
				BusLineID:            busLineIDFromDB,
				BusStopID:            busStopInfo.ID,
				NearestBusLinePathID: nearesetBusPathID,
			})

		}
	}

	if _, err := tx.NewInsert().Model(&busLineStops).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}
