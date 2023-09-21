package main

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/gojek/heimdall/httpclient"
)

const MAPBOX_DIRECTION_API_URL = "https://api.mapbox.com/directions/v5/mapbox/driving"

func MapBoxGetDirection(ctx context.Context, coords ...Coord) (*MapBoxDistanceAPIMinifiedResponse, error) {
	client := httpclient.NewClient(
		httpclient.WithHTTPTimeout(5 * time.Second),
	)

	additionalPath := ConvertCoordListToStringLongLat(coords...) + "?access_token=" + os.Getenv("MAPBOX_API_KEY")
	res, err := client.Get(MAPBOX_DIRECTION_API_URL+"/"+additionalPath, nil)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var responses MapBoxDistanceAPIMinifiedResponse
	if err := json.Unmarshal(body, &responses); err != nil {

		return nil, err
	}

	return &responses, nil
}
