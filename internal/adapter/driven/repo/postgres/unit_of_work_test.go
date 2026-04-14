package postgres_test

import (
	"testing"
	"time"

	"github.com/st-ember/streaming-log-consumer/internal/adapter/driven/repo/postgres"
	"github.com/st-ember/streaming-log-consumer/internal/domain/log"
	"github.com/st-ember/streaming-log-consumer/internal/domain/logretry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitOfWork(t *testing.T) {
	factory := postgres.NewPostgresUnitOfWorkFactory(testDB)

	t.Run("should commit successfully", func(t *testing.T) {
		t.Parallel()

		uow, err := factory.NewUnitOfWork(t.Context())
		require.NoError(t, err)

		l, _ := log.NewLog("l-1", "s-1", log.LevelInfo, log.CategoryDefault, "msg")

		err = uow.LogRepo().Save(t.Context(), l)
		require.NoError(t, err)

		err = uow.Commit(t.Context())
		require.NoError(t, err)

		// Verify existence using the main testDB connection
		var exists bool
		err = testDB.QueryRowContext(t.Context(), "SELECT EXISTS(SELECT 1 FROM logs WHERE id='l-1')").Scan(&exists)
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("should rollback successfully", func(t *testing.T) {
		t.Parallel()

		uow, err := factory.NewUnitOfWork(t.Context())
		require.NoError(t, err)

		l, _ := log.NewLog("l-2", "s-1", log.LevelInfo, log.CategoryDefault, "msg")

		err = uow.LogRepo().Save(t.Context(), l)
		require.NoError(t, err)

		err = uow.Rollback(t.Context())
		require.NoError(t, err)

		// Verify it does NOT exist
		var exists bool
		err = testDB.QueryRowContext(t.Context(), "SELECT EXISTS(SELECT 1 FROM logs WHERE id='l-2')").Scan(&exists)
		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("should rollback on multiple repos successfully", func(t *testing.T) {
		t.Parallel()

		uow, err := factory.NewUnitOfWork(t.Context())
		require.NoError(t, err)

		l, _ := log.NewLog("l-3", "s-1", log.LevelInfo, log.CategoryDefault, "msg")
		lr, _ := logretry.NewLogRetry("lr-1", l, "err", time.Now(), logretry.StatusPending)

		err = uow.LogRepo().Save(t.Context(), l)
		require.NoError(t, err)

		err = uow.LogRetryRepo().Save(t.Context(), lr)
		require.NoError(t, err)

		err = uow.Rollback(t.Context())
		require.NoError(t, err)

		// Verify BOTH are gone
		var logExists, retryExists bool
		err = testDB.QueryRowContext(t.Context(), "SELECT EXISTS(SELECT 1 FROM logs WHERE id='l-3')").Scan(&logExists)
		require.NoError(t, err)
		err = testDB.QueryRowContext(t.Context(), "SELECT EXISTS(SELECT 1 FROM log_retries WHERE id='lr-1')").Scan(&retryExists)
		require.NoError(t, err)

		assert.False(t, logExists)
		assert.False(t, retryExists)
	})
}
