package itmo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const schedulePath = "/internal/v1/schedule"

type Client struct {
	baseURL    string
	apiToken   string
	httpClient *http.Client
}

func NewClient(baseURL, apiToken string, timeout time.Duration) *Client {
	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		apiToken:   apiToken,
		httpClient: &http.Client{Timeout: timeout},
	}
}

func (c *Client) GetSchedule(ctx context.Context, from, to time.Time) ([]Event, error) {
	if from.After(to) {
		return nil, ErrInvalidPeriod
	}

	endpoint, err := url.Parse(c.baseURL + schedulePath)
	if err != nil {
		return nil, fmt.Errorf("parse adapter url: %w", err)
	}

	query := endpoint.Query()
	query.Set("from", from.Format(time.DateOnly))
	query.Set("to", to.Format(time.DateOnly))

	endpoint.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create schedule request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send schedule request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: status=%s", ErrUnexpectedResponse, resp.Status)
	}

	var scheduleResponse ScheduleResponse
	if err := json.NewDecoder(resp.Body).Decode(&scheduleResponse); err != nil {
		return nil, fmt.Errorf("decode schedule response: %w", err)
	}

	events := make([]Event, 0, len(scheduleResponse.Events))

	for _, raw := range scheduleResponse.Events {
		var event Event

		if err := json.Unmarshal(raw, &event); err != nil {
			return nil, fmt.Errorf("decode schedule event: %w", err)
		}
		event.RawPayload = append(json.RawMessage(nil), raw...)
		events = append(events, event)
	}
	return events, nil
}
