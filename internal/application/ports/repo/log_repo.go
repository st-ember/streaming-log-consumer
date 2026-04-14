package repo

import (
	"context"

	"github.com/st-ember/streaming-log-consumer/internal/domain/log"
)

type LogRepo interface {
	// Save inserts one log item into repo
	Save(ctx context.Context, log *log.Log) error
}
