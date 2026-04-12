package logretry

import (
	"fmt"
	"time"

	"github.com/st-ember/streaming-log-consumer/internal/domain/log"
)

// LogRetry tracks the errors encountered when storing a log.
// If a log ultimately succeeds, the error is deleted with a cron job.
type LogRetry struct {
	ID            string
	OriginalLog   *log.Log
	ErrorMsg      string
	RetryCount    int
	LastAttempt   time.Time
	NextAttemptAt time.Time
	Status        LogRetryStatus
	CreatedAt     time.Time
}

var RetryLimit = 5

func NewLogRetry(id string, ogLog *log.Log, errMsg string, nextAttemptAt time.Time, status LogRetryStatus) (*LogRetry, error) {
	if id == "" {
		return nil, ErrLogRetryIDEmpty
	}

	if ogLog == nil {
		return nil, ErrOriginalLogEmpty
	}

	if errMsg == "" {
		return nil, ErrMessageEmpty
	}

	if !status.IsValidForNew() {
		return nil, fmt.Errorf("%w: %s", ErrStatusInvalid, status)
	}

	return &LogRetry{
		ID:            id,
		OriginalLog:   ogLog,
		ErrorMsg:      errMsg,
		RetryCount:    1,
		NextAttemptAt: nextAttemptAt,
		LastAttempt:   time.Now(),
		Status:        status,
		CreatedAt:     time.Now(),
	}, nil
}

func (lr *LogRetry) Update(errMsg string, status LogRetryStatus, nextAttemptAt time.Time) error {
	if lr.Status == StatusDiscarded || lr.Status == StatusSucceeded {
		return fmt.Errorf("%w: %s", ErrStatusInvalidForUpdate, lr.Status)
	}

	lr.LastAttempt = time.Now()

	// Ignore updated retry count when setting to succeed
	if status == StatusSucceeded {
		lr.Status = status
		return nil
	}

	lr.ErrorMsg = errMsg
	lr.RetryCount++
	if lr.RetryCount >= RetryLimit {
		lr.Status = StatusDiscarded
		return nil
	}

	lr.Status = status
	lr.NextAttemptAt = nextAttemptAt

	return nil
}
