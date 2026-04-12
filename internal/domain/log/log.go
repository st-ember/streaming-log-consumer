package log

import (
	"fmt"
	"time"
)

type Log struct {
	ID        string
	SourceID  string
	Level     LogLevel
	Category  LogCategory
	Message   string
	CreatedAt time.Time
}

func NewLog(id, sourceID string, logLevel LogLevel, category LogCategory, msg string) (*Log, error) {
	if id == "" {
		return nil, ErrLogIDEmpty
	}

	if !logLevel.IsValid() {
		return nil, fmt.Errorf("%w: %s", ErrLogLevelInvalid, logLevel)
	}

	if !category.IsValid() {
		return nil, fmt.Errorf("%w: %s", ErrLogCategoryInvalid, category)
	}

	if category != CategoryDefault && sourceID == "" {
		return nil, ErrSourceIDEmpty
	}

	if msg == "" {
		return nil, ErrMessageEmpty
	}

	return &Log{
		ID:        id,
		SourceID:  sourceID,
		Level:     logLevel,
		Category:  category,
		Message:   msg,
		CreatedAt: time.Now(),
	}, nil
}
