package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/st-ember/streaming-log-consumer/internal/domain/logretry"
)

type PostgresLogRetryRepo struct {
	tx *sql.Tx
}

func NewPostgresLogRetryRepo(tx *sql.Tx) *PostgresLogRetryRepo {
	return &PostgresLogRetryRepo{tx}
}

// Save upserts one logRetry item into repo
func (r *PostgresLogRetryRepo) Save(ctx context.Context, logRetry *logretry.LogRetry) error {
	query := `
		INSERT INTO log_retries (id, original_log, error_msg, retry_count, 
		last_attempt, next_attempt_at, status, created_at)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (id) DO UPDATE SET
		error_msg = EXCLUDED.error_msg,
		retry_count = EXCLUDED.retry_count,
		last_attempt = EXCLUDED.last_attempt,
		next_attempt_at = EXCLUDED.next_attempt_at,
		status = EXCLUDED.status;
	`

	_, err := r.tx.ExecContext(
		ctx, query, logRetry.ID, LogJson{logRetry.OriginalLog},
		logRetry.ErrorMsg, logRetry.RetryCount, logRetry.LastAttempt,
		logRetry.NextAttemptAt, logRetry.Status, logRetry.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("save log retry %s: %w", logRetry.ID, err)
	}

	return nil
}

// FindCanRetry returns a limited slice of logRetry items that can be retried
func (r *PostgresLogRetryRepo) FindCanRetry(ctx context.Context, limit int) ([]*logretry.LogRetry, error) {
	query := `
		SELECT id, original_log, error_msg, retry_count, 
		last_attempt, next_attempt_at, status, created_at
		FROM log_retries
		WHERE (status = 'pending' or status = 'failed')
		AND next_attempt_at < CURRENT_TIMESTAMP
		ORDER BY last_attempt
		LIMIT $1
		FOR UPDATE SKIP LOCKED; -- to prevent race conditions 
	`
	rows, err := r.tx.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("find retriable logs: %w", err)
	}
	defer rows.Close()

	ls := make([]*logretry.LogRetry, 0, limit)
	for rows.Next() {
		var lj LogJson
		l := &logretry.LogRetry{}
		err := rows.Scan(
			&l.ID,
			&lj,
			&l.ErrorMsg,
			&l.RetryCount,
			&l.LastAttempt,
			&l.NextAttemptAt,
			&l.Status,
			&l.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan log retries: %w", err)
		}

		// Put scanned json back into domain struct
		l.OriginalLog = lj.Log

		ls = append(ls, l)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return ls, nil
}

// RemoveSucceeded batch deletes logRetry items with status of "succeeded"
func (r *PostgresLogRetryRepo) RemoveSucceeded(ctx context.Context) error {
	query := `
		DELETE FROM log_retries WHERE status = 'succeeded'
	`

	if _, err := r.tx.ExecContext(ctx, query); err != nil {
		return fmt.Errorf("batch remove succeeded retries: %w", err)
	}

	return nil
}
