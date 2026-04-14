package repo

import (
	"context"

	"github.com/st-ember/streaming-log-consumer/internal/domain/logretry"
)

type LogRetryRepo interface {
	// Save upserts one logRetry item into repo
	Save(ctx context.Context, logRetry *logretry.LogRetry) error
	// FindCanRetry returns a limited slice of logRetry items that can be retried
	FindCanRetry(ctx context.Context, limit int) ([]*logretry.LogRetry, error)
	// RemoveSucceeded batch deletes logRetry items with status of "succeeded"
	RemoveSucceeded(ctx context.Context) error
}
