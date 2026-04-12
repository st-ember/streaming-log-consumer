package log_test

import (
	"errors"
	"testing"

	"github.com/st-ember/streaming-log-consumer/internal/domain/log"
)

func TestNewLog(t *testing.T) {
	t.Run("should create a valid log", func(t *testing.T) {
		id := "log-123"
		sourceID := "src-456"
		level := log.LevelInfo
		category := log.CategoryVideo
		msg := "test message"

		l, err := log.NewLog(id, sourceID, level, category, msg)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if l.ID != id {
			t.Errorf("expected ID %s, got %s", id, l.ID)
		}
		if l.SourceID != sourceID {
			t.Errorf("expected SourceID %s, got %s", sourceID, l.SourceID)
		}
		if l.Message != msg {
			t.Errorf("expected Message %s, got %s", msg, l.Message)
		}
	})

	t.Run("should fail if ID is empty", func(t *testing.T) {
		_, err := log.NewLog("", "src-1", log.LevelInfo, log.CategoryDefault, "msg")
		if !errors.Is(err, log.ErrLogIDEmpty) {
			t.Errorf("expected ErrLogIDEmpty, got %v", err)
		}
	})

	t.Run("should fail if Level is invalid", func(t *testing.T) {
		_, err := log.NewLog("id", "src-1", log.LogLevel("invalid"), log.CategoryDefault, "msg")
		if !errors.Is(err, log.ErrLogLevelInvalid) {
			t.Errorf("expected ErrLogLevelInvalid, got %v", err)
		}
	})

	t.Run("should fail if Category is invalid", func(t *testing.T) {
		_, err := log.NewLog("id", "src-1", log.LevelInfo, log.LogCategory("invalid"), "msg")
		if !errors.Is(err, log.ErrLogCategoryInvalid) {
			t.Errorf("expected ErrLogCategoryInvalid, got %v", err)
		}
	})

	t.Run("should fail if SourceID is empty for non-default category", func(t *testing.T) {
		_, err := log.NewLog("id", "", log.LevelInfo, log.CategoryVideo, "msg")
		if !errors.Is(err, log.ErrSourceIDEmpty) {
			t.Errorf("expected ErrSourceIDEmpty, got %v", err)
		}
	})

	t.Run("should allow empty SourceID for default category", func(t *testing.T) {
		l, err := log.NewLog("id", "", log.LevelInfo, log.CategoryDefault, "msg")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if l.SourceID != "" {
			t.Errorf("expected empty SourceID, got %s", l.SourceID)
		}
	})

	t.Run("should fail if Message is empty", func(t *testing.T) {
		_, err := log.NewLog("id", "src-1", log.LevelInfo, log.CategoryDefault, "")
		if !errors.Is(err, log.ErrMessageEmpty) {
			t.Errorf("expected ErrMessageEmpty, got %v", err)
		}
	})
}
