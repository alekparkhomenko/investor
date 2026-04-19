package ingestor

import "errors"

var (
	ErrMOEXUnavailable = errors.New("MOEX API unavailable")
	ErrInvalidResponse = errors.New("invalid response from MOEX")
	ErrTimeout         = errors.New("request timeout")
)
