package aware

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

func (c *Client) PublishTelemetry(deviceID string, parameterName string, value interface{}, ts time.Time) error {
	data := struct {
		Timestamp     string      `json:"timestamp"`
		DeviceID      string      `json:"device"`
		ParameterName string      `json:"parameter"`
		Value         interface{} `json:"value"`
	}{formatTime(ts), deviceID, parameterName, value}

	header := Header{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	body, err := json.Marshal(&data)
	if err != nil {
		return err
	}

	res, err := c.request(context.Background(), http.MethodPost, c.server+"/v1/ingestion/ingest", body, header)
	if err != nil {
		return err
	}

	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusCreated {
		return formatUnexpectedResponse(res)
	}

	return nil
}
