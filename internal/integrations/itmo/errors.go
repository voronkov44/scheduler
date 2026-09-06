package itmo

import "errors"

var (
	ErrInvalidPeriod      = errors.New("invalid schedule period")
	ErrUnexpectedResponse = errors.New("unexpected response from ITMO adapter")
)
