package postgres_test

import (
	"context"
	"testing"

	"github.com/st-ember/streaming-log-consumer/internal/adapter/driven/repo/postgres"
	"github.com/st-ember/streaming-log-consumer/internal/domain/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogRepo_Save(t *testing.T) {
	tx := beginTx(t)
	repo := postgres.NewPostgresLogRepo(tx)
	ctx := context.Background()

	t.Run("should save a log successfully", func(t *testing.T) {
		l, _ := log.NewLog("l-1", "s-1", log.LevelInfo, log.CategoryDefault, "msg")

		err := repo.Save(ctx, l)
		require.NoError(t, err)

		var exists bool
		err = tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM logs WHERE id='l-1')").Scan(&exists)
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("should fail on duplicate ID", func(t *testing.T) {
		l1, _ := log.NewLog("l-dup", "s-1", log.LevelInfo, log.CategoryDefault, "msg")
		err := repo.Save(ctx, l1)
		require.NoError(t, err)

		l2, _ := log.NewLog("l-dup", "s-2", log.LevelInfo, log.CategoryDefault, "msg")
		err = repo.Save(ctx, l2)
		assert.Error(t, err)
	})
}
