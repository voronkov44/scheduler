package itmo

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"scheduler/internal/calendar"
	"strconv"
	"strings"
	"time"
)

type Service struct {
	client   *Client
	location *time.Location
}

func NewService(client *Client, location *time.Location) *Service {
	return &Service{
		client:   client,
		location: location,
	}
}

func (s *Service) Fetch(ctx context.Context, from, to time.Time) ([]calendar.ExternalEvent, error) {
	events, err := s.client.GetSchedule(ctx, from, to)
	if err != nil {
		return nil, fmt.Errorf("get ITMO schedule events: %w", err)
	}

	result := make([]calendar.ExternalEvent, 0, len(events))

	for _, event := range events {
		externalEvent, err := s.toExternalEvent(event)
		if err != nil {
			return nil, fmt.Errorf("convert ITMO event pair_id=%d: %w", event.PairID, err)
		}
		result = append(result, externalEvent)
	}
	return result, nil
}

func (s *Service) toExternalEvent(event Event) (calendar.ExternalEvent, error) {
	startsAt, err := parseDateTime(event.Date, event.TimeStart, s.location)
	if err != nil {
		return calendar.ExternalEvent{}, fmt.Errorf("parse start time: %w", err)
	}

	endsAt, err := parseDateTime(event.Date, event.TimeEnd, s.location)
	if err != nil {
		return calendar.ExternalEvent{}, fmt.Errorf("parse end time: %w", err)
	}

	hash := sha256.Sum256(event.RawPayload)

	return calendar.ExternalEvent{
		ExternalID:      strconv.FormatInt(event.PairID, 10),
		ExternalVersion: hex.EncodeToString(hash[:]),
		Title:           event.Subject,
		Description:     buildDescription(event),
		Location:        buildLocation(event),
		StartsAt:        startsAt,
		EndsAt:          endsAt,
		AllDay:          false,
		RawPayload:      event.RawPayload,
	}, nil
}

func parseDateTime(date string, clock string, location *time.Location) (time.Time, error) {
	value := date + " " + clock

	layouts := []string{
		"2006-01-02 15:04",
		"2006-01-02 15:04:05",
	}

	for _, layout := range layouts {
		result, err := time.ParseInLocation(layout, value, location)
		if err == nil {
			return result, nil
		}
	}

	return time.Time{}, fmt.Errorf("unsupported datetime: %q", value)
}

func buildDescription(event Event) string {
	values := []string{
		event.Type,
		event.WorkType,
		event.TeacherName,
		event.Group,
	}

	result := make([]string, 0, len(values))

	for _, value := range values {
		value = strings.TrimSpace(value)

		if value != "" {
			result = append(result, value)
		}
	}

	return strings.Join(result, " · ")
}

func buildLocation(event Event) string {
	switch {
	case event.Room != "" && event.Building != "":
		return event.Building + ", " + event.Room

	case event.Room != "":
		return event.Room

	default:
		return event.Building
	}
}
