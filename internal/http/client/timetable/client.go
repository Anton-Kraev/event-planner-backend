package timetable

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/Anton-Kraev/event-planner-backend/internal/domain/timetable"
)

const (
	api = "%s/api"
)

type (
	eventsResp struct {
		Events []timetable.Event `json:"events"`
	}

	Client struct {
		host       string
		httpClient *http.Client
	}
)

func NewClient(host string, httpClient *http.Client) Client {
	return Client{host: host, httpClient: httpClient}
}

func (c Client) doHTTP(ctx context.Context, method string, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, http.NoBody)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API Error: %w", err)
	}

	respB, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return respB, nil
}

func (c Client) getEvents(ctx context.Context, url string) ([]timetable.Event, error) {
	respB, err := c.doHTTP(ctx, http.MethodGet, url)
	if err != nil {
		return nil, err
	}

	if respB == nil || len(respB) == 0 {
		return nil, errors.New("no events found")
	}

	var events eventsResp

	if err = json.Unmarshal(respB, &events); err != nil {
		return nil, fmt.Errorf("failed to parse response body: %w", err)
	}

	return events.Events, nil
}
