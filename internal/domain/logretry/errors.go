package logretry

import "errors"

var (
	ErrLogRetryIDEmpty        = errors.New("log retry id cannot be empty")
	ErrLogIDEmpty             = errors.New("log id cannot be empty")
	ErrOriginalLogEmpty       = errors.New("original log cannot be empty")
	ErrMessageEmpty           = errors.New("error message cannot be empty")
	ErrStatusInvalid          = errors.New("status is invalid")
	ErrStatusInvalidForUpdate = errors.New("log retry status invalid for update")
	ErrRetryLimitExceeded     = errors.New("retry limit exceeded")
)
