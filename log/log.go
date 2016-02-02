package log

import (
	"os"
)

var std = New(os.Stderr)

// Logger defines an interface for logging messages
type Logger interface {
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
}
