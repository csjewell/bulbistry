package bulbistry

import (
	"bytes"
	"errors"
	"io"
	"log"
	"os"

	"github.com/google/uuid"
)

type Logger struct {
	logUUID *uuid.UUID
	*log.Logger
}

func NewLogger(w io.Writer, logUUID *uuid.UUID) *Logger {
	if w == nil {
		w = os.Stderr
	}
	if logUUID == nil {
		u := uuid.New()
		logUUID = &u
	}
	prefix := "[" + logUUID.String() + "] "
	l := log.New(w, prefix, log.Ldate|log.Ltime|log.Lmsgprefix)
	return &Logger{
		logUUID,
		l,
	}
}

type DebugLogger struct {
	activated   bool
	deactivated bool
	*bytes.Buffer
	*Logger
}

func NewDebugLogger(l *Logger) *DebugLogger {
	buf := bytes.NewBuffer(make([]byte, 0, 8192))
	return &DebugLogger{
		false,
		false,
		buf,
		NewLogger(buf, l.logUUID),
	}
}

func (dl *DebugLogger) TurnOn() error {
	if dl.activated {
		return errors.New("debug log already activated, cannot re-activate")
	}

	if dl.deactivated {
		return errors.New("debug log already deactivated, cannot activate")
	}

	_, err := dl.WriteTo(os.Stdout)
	if err != nil {
		return err
	}
	dl.SetOutput(os.Stdout)
	dl.activated = true
	return nil
}

func (dl *DebugLogger) TurnOff() error {
	if dl.activated {
		return errors.New("debug log already activated, cannot deactivate")
	}

	if dl.deactivated {
		return errors.New("debug log already deactivated, cannot re-deactivate")
	}

	dl.Truncate(0)
	dl.SetOutput(io.Discard)
	dl.deactivated = true
	return nil
}

func (dl *DebugLogger) Handled() bool {
	if dl.activated || dl.deactivated {
		return true
	}

	return false
}
