/*
Copyright © 2023 Curtis Jewell <swordsman@curtisjewell.name>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package bulbistry

import (
	"bytes"
	"errors"
	"io"
	"log"
	"os"

	"github.com/google/uuid"
)

// Bulbistry's extension of log.Logger to include the UUID in the logs.
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

// A logger at debugging level.
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
