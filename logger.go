/*
Copyright © 2023 Curtis Jewell <bulbistry@curtisjewell.name>

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
