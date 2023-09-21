package main

import (
	"context"
	"encoding/json"
	"io"
	"time"

	"github.com/gojek/heimdall/httpclient"
)

const UWAVE_BASE_URL = "https://test.uwave.sg"

func GetBusLinesUwave(ctx context.Context) ([]UwaveBusLineDetail, error) {
	client := httpclient.NewClient(
		httpclient.WithHTTPTimeout(5 * time.Second),
	)

	res, err := client.Get(UWAVE_BASE_URL+"/busLines", nil)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var responses UwaveBaseResponse[[]UwaveBusLineDetail]
	if err := json.Unmarshal(body, &responses); err != nil {

		return nil, err
	}

	return responses.Payload, nil
}

func GetBusPositionUwave(ctx context.Context, busLineID string) ([]UwaveBusPosition, error) {
	client := httpclient.NewClient(
		httpclient.WithHTTPTimeout(5 * time.Second),
	)

	res, err := client.Get(UWAVE_BASE_URL+"/busPositions/"+busLineID, nil)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var responses UwaveBaseResponse[[]UwaveBusPosition]
	if err := json.Unmarshal(body, &responses); err != nil {

		return nil, err
	}

	return responses.Payload, nil
}
