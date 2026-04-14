package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/st-ember/streaming-log-consumer/internal/adapter/driven/repo/postgres"
	"github.com/st-ember/streaming-log-consumer/internal/domain/log"
	"github.com/st-ember/streaming-log-consumer/internal/domain/logretry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogRetryRepo(t *testing.T) {
	tx := beginTx(t)
	repo := postgres.NewPostgresLogRetryRepo(tx)
	ctx := context.Background()

	mockLog, _ := log.NewLog("l-1", "s-1", log.LevelInfo, log.CategoryDefault, "msg")

	t.Run("Save should upsert successfully", func(t *testing.T) {
		lr, _ := logretry.NewLogRetry("lr-1", mockLog, "err", time.Now().Add(time.Minute), logretry.StatusPending)

		err := repo.Save(ctx, lr)
		require.NoError(t, err)

		// Update and Save again (upsert)
		lr.ErrorMsg = "new error"
		err = repo.Save(ctx, lr)
		require.NoError(t, err)

		var errMsg string
		err = tx.QueryRowContext(ctx, "SELECT error_msg FROM log_retries WHERE id='lr-1'").Scan(&errMsg)
		require.NoError(t, err)
		assert.Equal(t, "new error", errMsg)
	})

	t.Run("FindCanRetry should return eligible logs", func(t *testing.T) {
		// Past attempt
		lr1, _ := logretry.NewLogRetry("lr-past", mockLog, "err", time.Now().Add(-time.Minute), logretry.StatusPending)
		repo.Save(ctx, lr1)

		// Future attempt
		lr2, _ := logretry.NewLogRetry("lr-future", mockLog, "err", time.Now().Add(time.Minute), logretry.StatusPending)
		repo.Save(ctx, lr2)

		retries, err := repo.FindCanRetry(ctx, 10)
		require.NoError(t, err)

		assert.Len(t, retries, 1)
		assert.Equal(t, "lr-past", retries[0].ID)
	})

	t.Run("RemoveSucceeded should delete successfully", func(t *testing.T) {
		lr, _ := logretry.NewLogRetry("lr-success", mockLog, "err", time.Now(), logretry.StatusPending)
		repo.Save(ctx, lr)

		// Update to succeeded
		lr.Update("", logretry.StatusSucceeded, time.Now())
		repo.Save(ctx, lr)

		err := repo.RemoveSucceeded(ctx)
		require.NoError(t, err)

		var exists bool
		err = tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM log_retries WHERE id='lr-success')").Scan(&exists)
		require.NoError(t, err)
		assert.False(t, exists)
	})
}
