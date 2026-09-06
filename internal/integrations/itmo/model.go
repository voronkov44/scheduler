package itmo

import (
	"encoding/json"
	"time"
)

type ScheduleResponse struct {
	Source    string            `json:"source"`
	From      string            `json:"from"`
	To        string            `json:"to"`
	FetchedAt time.Time         `json:"fetched_at"`
	Events    []json.RawMessage `json:"events"`
}

type Event struct {
	PairID int64 `json:"pair_id"`

	Subject string `json:"subject"`

	Date      string `json:"date"`
	TimeStart string `json:"time_start"`
	TimeEnd   string `json:"time_end"`

	Room     string `json:"room"`
	Building string `json:"building"`

	Type     string `json:"type"`
	WorkType string `json:"work_type"`

	TeacherName string `json:"teacher_name"`
	Group       string `json:"group"`

	RawPayload json.RawMessage `json:"-"`
}
