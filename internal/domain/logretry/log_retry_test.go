package logretry_test

import (
	"errors"
	"testing"
	"time"

	"github.com/st-ember/streaming-log-consumer/internal/domain/log"
	"github.com/st-ember/streaming-log-consumer/internal/domain/logretry"
)

func TestNewLogRetry(t *testing.T) {
	mockLog, _ := log.NewLog("l-1", "s-1", log.LevelInfo, log.CategoryDefault, "msg")

	t.Run("should create a valid log retry", func(t *testing.T) {
		id := "lr-1"
		errMsg := "db timeout"
		next := time.Now().Add(time.Minute)
		status := logretry.StatusPending

		lr, err := logretry.NewLogRetry(id, mockLog, errMsg, next, status)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if lr.ID != id {
			t.Errorf("expected ID %s, got %s", id, lr.ID)
		}
		if lr.OriginalLog != mockLog {
			t.Errorf("expected OriginalLog to match")
		}
		if lr.RetryCount != 1 {
			t.Errorf("expected initial RetryCount to be 1, got %d", lr.RetryCount)
		}
	})

	t.Run("should fail if ID is empty", func(t *testing.T) {
		_, err := logretry.NewLogRetry("", mockLog, "err", time.Now(), logretry.StatusPending)
		if !errors.Is(err, logretry.ErrLogRetryIDEmpty) {
			t.Errorf("expected ErrLogRetryIDEmpty, got %v", err)
		}
	})

	t.Run("should fail if OriginalLog is nil", func(t *testing.T) {
		_, err := logretry.NewLogRetry("id", nil, "err", time.Now(), logretry.StatusPending)
		if !errors.Is(err, logretry.ErrOriginalLogEmpty) {
			t.Errorf("expected ErrOriginalLogEmpty, got %v", err)
		}
	})

	t.Run("should fail if errMsg is empty", func(t *testing.T) {
		_, err := logretry.NewLogRetry("id", mockLog, "", time.Now(), logretry.StatusPending)
		if !errors.Is(err, logretry.ErrMessageEmpty) {
			t.Errorf("expected ErrMessageEmpty, got %v", err)
		}
	})

	t.Run("should fail if status is invalid for new", func(t *testing.T) {
		_, err := logretry.NewLogRetry("id", mockLog, "err", time.Now(), logretry.StatusSucceeded)
		if !errors.Is(err, logretry.ErrStatusInvalid) {
			t.Errorf("expected ErrStatusInvalid, got %v", err)
		}
	})
}

func TestUpdate(t *testing.T) {
	mockLog, _ := log.NewLog("l-1", "s-1", log.LevelInfo, log.CategoryDefault, "msg")

	t.Run("should update retry count and status on failure", func(t *testing.T) {
		lr, _ := logretry.NewLogRetry("id", mockLog, "err1", time.Now(), logretry.StatusPending)
		
		newErr := "err2"
		next := time.Now().Add(time.Minute)
		err := lr.Update(newErr, logretry.StatusPending, next)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if lr.RetryCount != 2 {
			t.Errorf("expected RetryCount 2, got %d", lr.RetryCount)
		}
		if lr.ErrorMsg != newErr {
			t.Errorf("expected ErrorMsg %s, got %s", newErr, lr.ErrorMsg)
		}
	})

	t.Run("should set status to succeeded without incrementing retry count", func(t *testing.T) {
		lr, _ := logretry.NewLogRetry("id", mockLog, "err1", time.Now(), logretry.StatusPending)
		
		err := lr.Update("", logretry.StatusSucceeded, time.Now())

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if lr.Status != logretry.StatusSucceeded {
			t.Errorf("expected status succeeded, got %s", lr.Status)
		}
		if lr.RetryCount != 1 {
			t.Errorf("expected RetryCount to stay 1, got %d", lr.RetryCount)
		}
	})

	t.Run("should discard if retry limit is reached", func(t *testing.T) {
		lr, _ := logretry.NewLogRetry("id", mockLog, "err1", time.Now(), logretry.StatusPending)
		
		// Move to limit (RetryLimit = 5)
		// Initial is 1. Updates to 2, 3, 4, 5.
		for i := 0; i < 4; i++ {
			lr.Update("err", logretry.StatusPending, time.Now())
		}

		if lr.Status != logretry.StatusDiscarded {
			t.Errorf("expected StatusDiscarded, got %s", lr.Status)
		}
	})

	t.Run("should fail if updating a terminal state", func(t *testing.T) {
		lr, _ := logretry.NewLogRetry("id", mockLog, "err1", time.Now(), logretry.StatusPending)
		lr.Update("", logretry.StatusSucceeded, time.Now())

		err := lr.Update("err", logretry.StatusPending, time.Now())
		if !errors.Is(err, logretry.ErrStatusInvalidForUpdate) {
			t.Errorf("expected ErrStatusInvalidForUpdate, got %v", err)
		}
	})
}
