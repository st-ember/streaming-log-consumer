package repo

import "errors"

var (
	ErrData       = errors.New("malformed data")
	ErrConflict   = errors.New("temporary data conflict")
	ErrRateLimit  = errors.New("resource exhausted")
	ErrConnection = errors.New("connection failure")
)
