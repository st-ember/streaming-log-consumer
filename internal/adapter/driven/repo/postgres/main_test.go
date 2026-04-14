package postgres_test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	embeddedpgx "github.com/fergusstrange/embedded-postgres"
	"github.com/stretchr/testify/require"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	config := embeddedpgx.DefaultConfig().
		Port(5433).
		Logger(nil)

	// Create embdedded postgres
	postgres := embeddedpgx.NewDatabase(config)
	if err := postgres.Start(); err != nil {
		log.Fatalf("start embdedded postgres: %v", err)
	}

	// Connect to embedded postgres
	connString := "postgres://postgres:postgres@localhost:5433/postgres?sslmode=disable"
	var err error
	testDB, err = sql.Open("pgx", connString)
	if err != nil {
		log.Fatalf("connect to embdedded postgres: %v", err)
	}

	// Create database schema
	createTablesSQL := `
        CREATE TABLE IF NOT EXISTS logs (
			id TEXT PRIMARY KEY, source_id TEXT, level TEXT, category TEXT, message TEXT, created_at TIMESTAMPTZ
        );
        CREATE TABLE IF NOT EXISTS log_retries (
			id TEXT PRIMARY KEY, original_log JSONB, error_msg TEXT, retry_count INT, 
			last_attempt TIMESTAMPTZ, next_attempt_at TIMESTAMPTZ, status TEXT, created_at TIMESTAMPTZ
        );
	`
	_, err = testDB.Exec(createTablesSQL)
	if err != nil {
		log.Fatalf("create tables in embdedded postgres: %v", err)
	}

	// Run all the tests in the package
	code := m.Run()

	// Stop embdedded postgres
	if err := postgres.Stop(); err != nil {
		log.Fatalf("stop embdedded postgres: %v", err)
	}

	if err := testDB.Close(); err != nil {
		log.Fatalf("close test db connection: %v", err)
	}

	// Exit with the tests' exit code
	os.Exit(code)
}

func beginTx(t *testing.T) *sql.Tx {
	tx, err := testDB.BeginTx(t.Context(), nil)
	require.NoError(t, err)

	_, err = tx.ExecContext(t.Context(), "TRUNCATE logs, log_retries RESTART IDENTITY CASCADE;")
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = tx.Rollback()
	})

	return tx
}
