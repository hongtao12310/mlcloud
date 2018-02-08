package log

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

var (
	message = "message"
)

func TestSetx(t *testing.T) {
	logger := New(nil, nil, WarningLevel)
	logger.SetOutput(os.Stdout)
	fmt := NewTextFormatter()
	logger.SetFormatter(fmt)
	logger.SetLevel(DebugLevel)

	if logger.out != os.Stdout {
		t.Errorf("unexpected outer: %v != %v", logger.out, os.Stdout)
	}

	if logger.fmtter != fmt {
		t.Errorf("unexpected formatter: %v != %v", logger.fmtter, fmt)
	}

	if logger.lvl != DebugLevel {
		t.Errorf("unexpected log level: %v != %v", logger.lvl, DebugLevel)
	}
}

func TestDebug(t *testing.T) {
	buf := enter()
	defer exit()

	Debug(message)

	str := buf.String()
	if len(str) != 0 {
		t.Errorf("unexpected message: %s != %s", str, "")
	}
}

func TestDebugf(t *testing.T) {
	buf := enter()
	defer exit()

	Debugf("%s", message)

	str := buf.String()
	if len(str) != 0 {
		t.Errorf("unexpected message: %s != %s", str, "")
	}
}

func TestInfo(t *testing.T) {
	buf := enter()
	defer exit()

	Info(message)

	str := buf.String()
	if strings.HasSuffix(str, "[INFO] message") {
		t.Errorf("unexpected message: %s != %s", str, "")
	}
}

func TestInfof(t *testing.T) {
	buf := enter()
	defer exit()

	Infof("%s", message)

	str := buf.String()
	if strings.HasSuffix(str, "[INFO] message") {
		t.Errorf("unexpected message: %s != %s", str, "")
	}
}

func TestWarning(t *testing.T) {
	buf := enter()
	defer exit()

	Warning(message)

	str := buf.String()
	if strings.HasSuffix(str, "[WARNING] message") {
		t.Errorf("unexpected message: %s != %s", str, "")
	}
}

func TestWarningf(t *testing.T) {
	buf := enter()
	defer exit()

	Warningf("%s", message)

	str := buf.String()
	if strings.HasSuffix(str, "[WARNING] message") {
		t.Errorf("unexpected message: %s != %s", str, "")
	}
}

func TestError(t *testing.T) {
	buf := enter()
	defer exit()

	Error(message)

	str := buf.String()
	if strings.HasSuffix(str, "[ERROR] message") {
		t.Errorf("unexpected message: %s != %s", str, "")
	}
}

func TestErrorf(t *testing.T) {
	buf := enter()
	defer exit()

	Errorf("%s", message)

	str := buf.String()
	if strings.HasSuffix(str, "[ERROR] message") {
		t.Errorf("unexpected message: %s != %s", str, "")
	}
}

func enter() *bytes.Buffer {
	b := make([]byte, 0, 32)
	buf := bytes.NewBuffer(b)

	logger.SetOutput(buf)

	return buf
}

func exit() {
	logger.SetOutput(os.Stdout)
}
