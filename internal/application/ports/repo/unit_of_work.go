package repo

import "context"

type UnitOfWork interface {
	// Repository functions returns the transactional repositories
	LogRepo() LogRepo
	LogRetryRepo() LogRetryRepo
	// Commit finalizes the transaction
	Commit(ctx context.Context) error
	// Rollback cancels the transaction
	Rollback(ctx context.Context) error
	// Close ends read-only transactions
	Close(ctx context.Context) error
}

type UnitOfWorkFactory interface {
	NewUnitOfWork(ctx context.Context) (UnitOfWork, error)
}
