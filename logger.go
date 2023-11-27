package bulbistry

import (
	"bytes"
	"io"
	"log"
	"os"

	"github.com/google/uuid"
)

type Logger struct {
	logUUID uuid.UUID
	*log.Logger
}

func NewLogger(w io.Writer, *logUUID uuid.UUID) Logger {
	if w == nil {
		w = os.Stderr
	}
	if logUUID == nil {
		logUUID = uuid.New()
	}
	prefix := "[" + logUUID.String() + "] "
	return &Logger{
		logUUID,
		log.New(w, prefix, log.Ldate | log.Ltime | log.Lmsgprefix)
	}
}

type DebugLogger struct {
	activated bool
	deactivated bool
	*bytes.Buffer
	*Logger
}

func NewDebugLogger(l Logger) DebugLogger {
	buf := bytes.NewBuffer(make([]bytes, 0, 8192))
	return &DebugLogger{
		false,
		false,
		buf,
		NewLogger(buf, l.logUUID)
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
	dl.activated = true;
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
	dl.SetOutput(io.Discard())
	dl.deactivated = true
	return nil
}

func (dl *DebugLogger) Handled() bool {
	if dl.activated || dl.deactivated {
		return true
	}

	return false
}


