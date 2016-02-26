package tracerr

import (
	"bytes"
	"fmt"
	"runtime"
)

// Error is used to combine multiple errors into one containing stack trace information
type Error struct {
	msgs  []string
	err   error
	trace stackTrace
}

// Error returns the error message
func (we *Error) Error() string {
	return we.err.Error()
}

// Unwrap returns the original error
func (we *Error) Unwrap() error {
	return we.err
}

// TraceString returns a string representing the stack trace of the error
func (we *Error) TraceString() string {
	return fmt.Sprintf(
		"error:\n%s\n\ntrace:\n%s\n\nmessages:\n%v",
		we.Error(),
		we.trace.String(),
		we.msgsString(),
	)
}

func (we *Error) msgsString() string {
	var buffer bytes.Buffer
	for i := len(we.msgs) - 1; i >= 0; i-- {
		buffer.WriteString(we.msgs[i])
		buffer.WriteString("\n")
	}
	return buffer.String()
}

// Wrap returns a new error that wraps the given one adding stack trace info to it
func Wrap(err error, message string) *Error {
	if werr, ok := err.(*Error); ok {
		werr.msgs = append(werr.msgs, message)
		return werr
	}

	return &Error{
		err:   err,
		msgs:  []string{message},
		trace: newStackTrace(1),
	}
}

type stackFrame struct {
	pc   uintptr
	file string
	line int
}

type stackTrace []stackFrame

func newStackTrace(skip int) stackTrace {
	var st stackTrace
	var n int
	for {
		pc, file, line, b := runtime.Caller(skip + 1)
		if !b {
			return st
		}
		st = append(st, stackFrame{
			pc:   pc,
			file: file,
			line: line,
		})
		skip++
		n++
	}
}

// String returns a string representation of the stackTrace
func (st stackTrace) String() string {
	var str string
	for _, sf := range st {
		str += fmt.Sprintf("%s:%d\n", sf.file, sf.line)
	}
	return str
}
