package log

// NilLogger is a nil implementation of the Logger interface
type NilLogger struct{}

// Error is a no-op implementation of the Logger.Error method
func (l NilLogger) Error(v ...interface{}) {}

// Errorf is a no-op implementation of the Logger.Error method
func (l NilLogger) Errorf(format string, v ...interface{}) {}

// Info is a no-op implementation of the Logger.Info method
func (l NilLogger) Info(v ...interface{}) {}

// Infof is a no-op implementation of the Logger.Infof method
func (l NilLogger) Infof(format string, v ...interface{}) {}

// Debug is a no-op implementation of the Logger.Debug method
func (l NilLogger) Debug(v ...interface{}) {}

// Debugf is a no-op implementation of the Logger.Debugf method
func (l NilLogger) Debugf(format string, v ...interface{}) {}
