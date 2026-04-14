package postgres

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/st-ember/streaming-log-consumer/internal/domain/log"
)

// Postgres Error Codes (SQLSTATE)
const (
	// Data Errors (repo.ErrData)
	CodeUniqueViolation     = "23505"
	CodeNotNullViolation    = "23502"
	CodeForeignKeyViolation = "23503"
	CodeCheckViolation      = "23514"
	CodeStringTruncation    = "22001"
	CodeInvalidText         = "22P02"
	CodeDataException       = "22000"

	// Conflict Errors (repo.ErrConflict)
	CodeDeadlockDetected     = "40P01"
	CodeSerializationFailure = "40001"
	CodeLockNotAvailable     = "55P03"

	// Rate Limit Errors (repo.ErrRateLimit)
	CodeTooManyConnections    = "53300"
	CodeInsufficientResources = "53000"
	CodeOutOfMemory           = "53200"
	CodeConfigLimitExceeded   = "53400"
	CodeProgramLimitExceeded  = "54000"
	CodeQueryCanceled         = "57014"

	// Connection Errors (repo.ErrConnection)
	CodeConnectionException    = "08000"
	CodeConnectionDoesNotExist = "08003"
	CodeConnectionFailure      = "08006"
	CodeAdminShutdown          = "57P01"
	CodeCannotConnectNow       = "57P03"
	CodeReadOnlyTransaction    = "25006"
)

// LogJSON is a wrapper that implements sql.Scanner and driver.Valuer
type LogJson struct {
	*log.Log
}

// Value tells Postgres how to store the Log (converts to JSON)
func (lj *LogJson) Value() (driver.Value, error) {
	if lj.Log == nil {
		return nil, nil
	}

	return json.Marshal(lj.Log)
}

// Scan tells Postgres how to read the Log (converts from JSON)
func (lj *LogJson) Scan(value any) error {
	if value == nil {
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &lj.Log)
}
