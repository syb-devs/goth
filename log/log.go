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

// Error logs an error in the global logger
func Error(v ...interface{}) {
	std.Error(v...)
}

// Errorf logs a formatted error in the global logger
func Errorf(format string, v ...interface{}) {
	std.Errorf(format, v...)
}

// Info logs an info message in the global logger
func Info(v ...interface{}) {
	std.Info(v...)
}

// Infof logs a formatted info message in the global logger
func Infof(format string, v ...interface{}) {
	std.Infof(format, v...)
}

// Debug logs a debug message in the global logger
func Debug(v ...interface{}) {
	std.Debug(v...)
}

// Debugf logs a formatted debug message in the global logger
func Debugf(format string, v ...interface{}) {
	std.Debugf(format, v...)
}
