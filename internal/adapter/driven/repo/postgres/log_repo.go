package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/st-ember/streaming-log-consumer/internal/application/ports/repo"
	"github.com/st-ember/streaming-log-consumer/internal/domain/log"
)

type PostgresLogRepo struct {
	tx *sql.Tx
}

func NewPostgresLogRepo(tx *sql.Tx) *PostgresLogRepo {
	return &PostgresLogRepo{tx}
}

// Save inserts one log item into repo
func (r *PostgresLogRepo) Save(ctx context.Context, log *log.Log) error {
	query := `
		INSERT INTO logs (id, source_id, level, category, message, created_at)
		VALUES($1, $2, $3, $4, $5, $6)
	`

	_, err := r.tx.ExecContext(ctx, query, log.ID, log.SourceID, log.Level, log.Category, log.Message, log.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case CodeUniqueViolation, CodeNotNullViolation,
				CodeCheckViolation, CodeInvalidText, CodeDataException:
				return repo.ErrData
			case CodeDeadlockDetected, CodeSerializationFailure, CodeLockNotAvailable:
				return repo.ErrConflict
			case CodeTooManyConnections, CodeInsufficientResources, CodeOutOfMemory,
				CodeConfigLimitExceeded, CodeProgramLimitExceeded:
				return repo.ErrRateLimit
			case CodeConnectionException, CodeConnectionDoesNotExist, CodeConnectionFailure,
				CodeAdminShutdown, CodeCannotConnectNow, CodeReadOnlyTransaction:
				return repo.ErrConnection
			}
		}
		return fmt.Errorf("save log %s: %w", log.ID, err)
	}

	return nil
}
