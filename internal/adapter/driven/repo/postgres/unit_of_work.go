package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/st-ember/streaming-log-consumer/internal/application/ports/repo"
)

// PostgresUnitOfWork implements the UnitOfWork interface for PostgreSQL
// and holds the transaction object.
type PostgresUnitOfWork struct {
	tx *sql.Tx
}

// PostgresUnitOfWorkFactory is the factory for creating UoW instances
// and holds the connection pool
type PostgresUnitOfWorkFactory struct {
	db *sql.DB
}

// NewPostgresUnitOfWorkFactory creates a new factory instance.
func NewPostgresUnitOfWorkFactory(db *sql.DB) repo.UnitOfWorkFactory {
	return &PostgresUnitOfWorkFactory{db}
}

// NewUnitOfWork is called every time a database transaction is needed
func (f *PostgresUnitOfWorkFactory) NewUnitOfWork(ctx context.Context) (repo.UnitOfWork, error) {
	tx, err := f.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("start new transaction: %w", err)
	}

	return &PostgresUnitOfWork{tx}, nil
}

// LogRepo returns a new PostgresLogRepo that uses the UoW's transaction
func (u *PostgresUnitOfWork) LogRepo() repo.LogRepo {
	return NewPostgresLogRepo(u.tx)
}

// LogRetryRepo returns a new PostgresLogRepo that uses the UoW's transaction
func (u *PostgresUnitOfWork) LogRetryRepo() repo.LogRetryRepo {
	return NewPostgresLogRetryRepo(u.tx)
}

// Commit finalizes the transaction
func (u *PostgresUnitOfWork) Commit(ctx context.Context) error {
	return u.tx.Commit()
}

// Rollback cancels the transaction
func (u *PostgresUnitOfWork) Rollback(ctx context.Context) error {
	return u.tx.Rollback()
}

// Close ends read-only transactions
func (u *PostgresUnitOfWork) Close(ctx context.Context) error {
	return u.tx.Rollback()
}
