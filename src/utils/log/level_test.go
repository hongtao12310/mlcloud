package log

import (
	"testing"
)

func TestString(t *testing.T) {
	m := map[Level]string{
		DebugLevel:   "DEBUG",
		InfoLevel:    "INFO",
		WarningLevel: "WARNING",
		ErrorLevel:   "ERROR",
		FatalLevel:   "FATAL",
		-1:           "UNKNOWN",
	}

	for level, str := range m {
		if level.string() != str {
			t.Errorf("unexpected string: %s != %s", level.string(), str)
		}
	}
}

func TestParseLevel(t *testing.T) {
	m := map[string]Level{
		"DEBUG":   DebugLevel,
		"INFO":    InfoLevel,
		"WARNING": WarningLevel,
		"ERROR":   ErrorLevel,
		"FATAL":   FatalLevel,
	}

	for str, level := range m {
		l, err := parseLevel(str)
		if err != nil {
			t.Errorf("failed to parse level: %v", err)
		}
		if l != level {
			t.Errorf("unexpected level: %d != %d", l, level)
		}
	}

	if _, err := parseLevel("UNKNOWN"); err == nil {
		t.Errorf("unexpected behaviour: should be error here")
	}
}
