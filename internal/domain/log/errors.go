package log

import "errors"

var (
	ErrLogIDEmpty         = errors.New("log id cannot be empty")
	ErrLogLevelInvalid    = errors.New("log level is invalid")
	ErrLogCategoryInvalid = errors.New("log category is invalid")
	ErrSourceIDEmpty      = errors.New("source id cannot be empty when category is not default")
	ErrMessageEmpty       = errors.New("log message cannot be empty")
)
