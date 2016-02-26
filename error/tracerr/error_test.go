package tracerr_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/syb-devs/tracerr"
)

var origErr = errors.New("Whoops!")

func fails() error {
	f := func() error {
		return origErr
	}
	return f()
}

func failsToo() error {
	return tracerr.Wrap(fails(), "Oh my God! This is going to explode!!")
}

func failsAgain() error {
	return tracerr.Wrap(failsToo(), "something went really wrong")
}

func middleFunc() error {
	return failsAgain()
}

func first() error {
	return tracerr.Wrap(middleFunc(), "first failed")
}

func TestTraceString(t *testing.T) {
	patterns := []string{
		"error:\nWhoops!",
		"first failed",
		"something went really wrong",
		"Oh my God! This is going to explode!!",
		"tracerr/error_test.go:21",
		"tracerr/error_test.go:25",
		"tracerr/error_test.go:29",
		"tracerr/error_test.go:33",
		"tracerr/error_test.go:49",
	}

	err := first()
	werr := err.(*tracerr.Error)

	uerr := werr.Unwrap()
	if uerr != origErr {
		t.Error("unwrapping error")
		t.Errorf("want:\n %v", origErr)
		t.Errorf("have:\n %v", uerr)
	}

	trace := werr.TraceString()
	for _, pattern := range patterns {
		i := strings.Index(trace, pattern)
		if i == -1 {
			t.Error("pattern not found in stack trace string")
			t.Errorf("want:\n %s", pattern)
			t.Errorf("stack trace string:\n %s", trace)
		}
	}
}
