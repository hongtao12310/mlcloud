package log

import (
	"fmt"
	"time"
)

var defaultTimeFormat = time.RFC3339 // 2006-01-02T15:04:05Z07:00

// TextFormatter represents a kind of formatter that formats the logs as plain text
type TextFormatter struct {
	timeFormat string
}

// NewTextFormatter returns a TextFormatter, the format of time is time.RFC3339
func NewTextFormatter() *TextFormatter {
	return &TextFormatter{
		timeFormat: defaultTimeFormat,
	}
}

// Format formats the logs as "time [level] line message"
func (t *TextFormatter) Format(r *Record) (b []byte, err error) {
	s := fmt.Sprintf("%s [%s] ", r.Time.Format(t.timeFormat), r.Lvl.string())

	if len(r.Line) != 0 {
		s = s + r.Line + " "
	}

	if len(r.Msg) != 0 {
		s = s + r.Msg
	}

	b = []byte(s)

	if len(b) == 0 || b[len(b)-1] != '\n' {
		b = append(b, '\n')
	}

	return
}

// SetTimeFormat sets time format of TextFormatter if the parameter fmt is not null
func (t *TextFormatter) SetTimeFormat(fmt string) {
	if len(fmt) != 0 {
		t.timeFormat = fmt
	}
}
