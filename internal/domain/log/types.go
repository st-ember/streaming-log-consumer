package log

type LogLevel string

const (
	LevelInfo  LogLevel = "info"
	LevelWarn  LogLevel = "warn"
	LevelError LogLevel = "error"
)

func (lt LogLevel) IsValid() bool {
	switch lt {
	case LevelInfo:
		return true
	case LevelWarn:
		return true
	case LevelError:
		return true
	default:
		return false
	}
}

type LogCategory string

const (
	CategoryDefault LogCategory = "default"
	CategoryVideo   LogCategory = "video"
	CategoryJob     LogCategory = "job"
)

func (lc LogCategory) IsValid() bool {
	switch lc {
	case CategoryDefault:
		return true
	case CategoryVideo:
		return true
	case CategoryJob:
		return true
	default:
		return false
	}
}
