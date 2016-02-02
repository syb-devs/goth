package log

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	// LevelError is used for error messages
	LevelError = iota
	// LevelInfo is used for informational messages
	LevelInfo
	// LevelDebug is used for debug messages
	LevelDebug

	colorBlack   = 30
	colorRed     = 31
	colorGreen   = 32
	colorYellow  = 33
	colorBlue    = 34
	colorMagenta = 35
	colorCyan    = 36
	colorWhite   = 97

	colorReset = "\033[0m"
)

var logLevelColors map[int]string

func init() {
	logLevelColors = getLevelColors()
}

// NowFunc is a type of function that returns a time. Useful for unit testing
type NowFunc func() time.Time

// WLogger implements the WLogger interface using a Writer to log to
type WLogger struct {
	mu       sync.Mutex
	writer   io.Writer
	level    int
	prefix   string
	pattern  string
	coloring bool
	nowFunc  NowFunc
}

// check statically that WLogger implements the Logger interface
var _ Logger = (*WLogger)(nil)

// New returns a new WLogger, which uses a writer to write the messages
func New(w io.Writer) *WLogger {
	return &WLogger{
		writer:   w,
		level:    LevelDebug,
		pattern:  "{{ color }}{{ time }} {{ prefix }} [{{ level_literal }}] {{ message }}{{ color_reset }}\n",
		coloring: true,
		nowFunc:  time.Now,
	}
}

// SetLevel sets the threshold level for the logger.
// Only messages with a level lower or equal to the threshold level will be written.
// You can use the defined LevelXXX constants to set it.
// Ex: logger.SetLevel(log.LevelDebug)
func (l *WLogger) SetLevel(level int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// SetPrefix sets the prefix for the log lines.
// This is helpful to filter log contents.
func (l *WLogger) SetPrefix(prefix string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.prefix = prefix
}

// SetPattern sets the log line patter for the logger.
// Defined tokens are:
// {{ time }} - the actual time of the logged event
// {{ level_literal }} - the literal representation of the severity level
// {{ level }} - the numeric severity level
// {{ message }} - the message beign logged
// {{ prefix }} - the prefix set to the logger (if any)
// {{ color }} - the terminal escape sequence for the color assigned to the log level
// {{ color_reset }} - the terminal escape sequence for resetting the coloring (foreground and background)
func (l *WLogger) SetPattern(pattern string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.pattern = pattern
}

// SetColoring sets the coloring on/off
func (l *WLogger) SetColoring(b bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.coloring = b
}

// SetNowFunc sets a custom function for getting the log event time
func (l *WLogger) SetNowFunc(nowFunc NowFunc) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.nowFunc = nowFunc
}

func (l *WLogger) log(level int, message string) error {
	if level > l.level {
		return nil
	}

	var colorSeq, colorOff string
	if l.coloring {
		colorSeq = levelColor(level)
		colorOff = colorReset
	}

	line := l.pattern
	line = strings.Replace(line, "{{ time }}", l.now(), -1)
	line = strings.Replace(line, "{{ level }}", strconv.Itoa(level), -1)
	line = strings.Replace(line, "{{ level_literal }}", strings.ToUpper(levelString(level)), -1)
	line = strings.Replace(line, "{{ prefix }}", l.prefix, -1)
	line = strings.Replace(line, "{{ message }}", message, -1)
	line = strings.Replace(line, "{{ color }}", colorSeq, -1)
	line = strings.Replace(line, "{{ color_reset }}", colorOff, -1)

	_, err := l.writer.Write([]byte(line))
	return err
}

func (l *WLogger) now() string {
	return l.nowFunc().Format(time.RFC3339)
}

// Error logs the given message(s) with error level
func (l *WLogger) Error(v ...interface{}) {
	l.log(LevelError, fmt.Sprint(v...))
}

// Errorf formats the message and logs it with error level
func (l *WLogger) Errorf(format string, v ...interface{}) {
	l.log(LevelError, fmt.Sprintf(format, v...))
}

// Info logs the given message(s) with info level
func (l *WLogger) Info(v ...interface{}) {
	l.log(LevelInfo, fmt.Sprint(v...))
}

// Infof formats the message and logs it with info level
func (l *WLogger) Infof(format string, v ...interface{}) {
	l.log(LevelInfo, fmt.Sprintf(format, v...))
}

// Debug logs the given message(s) with debug level
func (l *WLogger) Debug(v ...interface{}) {
	l.log(LevelDebug, fmt.Sprint(v...))
}

// Debugf formats the message and logs it with debug level
func (l *WLogger) Debugf(format string, v ...interface{}) {
	l.log(LevelDebug, fmt.Sprintf(format, v...))
}

func getLevelColors() map[int]string {
	return map[int]string{
		LevelError: colorEscape(colorRed, false),
		LevelInfo:  colorEscape(colorGreen, false),
		LevelDebug: colorEscape(colorCyan, false),
	}
}

func colorEscape(color int, bold bool) string {
	if bold {
		return fmt.Sprintf("\033[%d;1m", color)
	}
	return fmt.Sprintf("\033[%dm", color)
}

func levelColor(level int) string {
	return logLevelColors[level]
}

func levelString(level int) string {
	switch level {
	case LevelError:
		return "error"
	case LevelInfo:
		return "info"
	case LevelDebug:
		return "debug"
	default:
		return ""
	}
}
