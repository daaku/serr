// Package serr provides stackful errors and nothing else.
//
// A traditional method to augment errors is stack traces. The `serr` package
// allows for programmers to add stack traces to errors without destroying the
// original error value.
//
// It provides pretty formatting and structured access to the stack traces.
// Structured access is in terms of the `[]uintptr` collected via
// `runtime.Callers`. Pretty formatting is provided in terms of support for the
// `%+v` style format modifier.
//
// This is the entire extent of this library. Anything else is out of scope.
package serr

import (
	"errors"
	"fmt"
	"io"
	"runtime"
	"strconv"
)

// Error wraps an error and augments it to include a single stack trace.
type Error struct {
	err   error
	stack []uintptr
}

// Unwrap provides suport for standard library error unwrapping.
func (e *Error) Unwrap() error {
	return e.err
}

// Error provides the error string of the underlying error without modification.
func (e *Error) Error() string {
	return e.err.Error()
}

// Format adds support for `%+v` style formatting to include stack trace
// details. This is suitable for plain text scenarios such as the terminal.
func (e *Error) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v", e.err)
			frames := runtime.CallersFrames(e.stack)
			for {
				frame, more := frames.Next()
				io.WriteString(s, "\n")
				io.WriteString(s, frame.Function)
				io.WriteString(s, "\n\t")
				io.WriteString(s, frame.File)
				io.WriteString(s, ":")
				io.WriteString(s, strconv.Itoa(frame.Line))
				if !more {
					break
				}
			}
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, e.err.Error())
	case 'q':
		fmt.Fprintf(s, "%q", e.err.Error())
	}
}

// Callers returns the callers collected at the time of Wrapping the error.
func (e *Error) Callers() []uintptr {
	return e.stack
}

// Wrap is a wrapper around `WrapSkip(err, 1)`.
func Wrap(err error) error {
	return WrapSkip(err, 1)
}

// WrapSkip returns an error including the callers respecting skip.
//
// For errors that are already wrapped, it returns a new error by extending it
// to include the additional stack frames. This may look confusing with errors
// being shuffled between goroutines with little or no shared stacks.
func WrapSkip(err error, skip int) error {
	if err == nil {
		return nil
	}

	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(skip+2, pcs[:])

	stack := pcs[0:n]
	if we, ok := errors.AsType[*Error](err); ok {
		err = we.err
		stack = mergeStack(stack, we.stack)
	}

	return &Error{
		err:   err,
		stack: stack,
	}
}

// Errorf is a wrapper around `WrapSkip(fmt.Errorf(...), 1)`.
func Errorf(format string, a ...any) error {
	return WrapSkip(fmt.Errorf(format, a...), 1)
}

func mergeStack(additional, existing []uintptr) []uintptr {
	commonLen := 0
	for i, j := len(existing)-1, len(additional)-1; i >= 0 && j >= 0; i, j = i-1, j-1 {
		if existing[i] != additional[j] {
			break
		}
		commonLen++
	}
	startOfAdditional := len(additional) - commonLen
	return append(additional[:startOfAdditional], existing...)
}
