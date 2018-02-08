package log

import (
	"fmt"
	"strings"
)

// Level ...
type Level int

const (
	// DebugLevel debug
	DebugLevel Level = iota
	// InfoLevel info
	InfoLevel
	// WarningLevel warning
	WarningLevel
	// ErrorLevel error
	ErrorLevel
	// FatalLevel fatal
	FatalLevel
)

func (l Level) string() (lvl string) {
	switch l {
	case DebugLevel:
		lvl = "DEBUG"
	case InfoLevel:
		lvl = "INFO"
	case WarningLevel:
		lvl = "WARNING"
	case ErrorLevel:
		lvl = "ERROR"
	case FatalLevel:
		lvl = "FATAL"
	default:
		lvl = "UNKNOWN"
	}

	return
}

func parseLevel(lvl string) (level Level, err error) {

	switch strings.ToLower(lvl) {
	case "debug":
		level = DebugLevel
	case "info":
		level = InfoLevel
	case "warning":
		level = WarningLevel
	case "error":
		level = ErrorLevel
	case "fatal":
		level = FatalLevel
	default:
		err = fmt.Errorf("invalid log level: %s", lvl)
	}

	return
}
