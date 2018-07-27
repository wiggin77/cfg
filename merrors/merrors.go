package merrors

import (
	"bytes"
	"strconv"
)

// MultiError represents zero or more errors that can be
// accumulated via the `Append` method.
type MultiError struct {
	max      int
	errors   []error
	overflow int
}

// New returns a new instance of MultiError. `max` determines
// the maximum number of errors that can be appended, after which
// only the overflow counter will be incremented.
func New(max int) *MultiError {
	me := MultiError{}
	me.max = max
	me.errors = make([]error, 0, 10)
	me.overflow = 0
	return &me
}

// Append adds an error to the combined error list
func (me *MultiError) Append(err error) {
	if err == nil {
		return
	}
	if len(me.errors) >= me.max {
		me.overflow++
	} else {
		me.errors = append(me.errors, err)
	}
}

// Len returns the number of errors that have been appended.
func (me *MultiError) Len() int {
	return len(me.errors)
}

// Error returns a string representation of this MultiError.
// All the errors are converted to string, delimited with newlines,
// and the number of overflow errors appended to the end.
func (me *MultiError) Error() string {
	buf := bytes.Buffer{}
	buf.WriteString("MultiError:\n")
	for _, err := range me.errors {
		buf.WriteString(err.Error())
		buf.WriteString("\n")
	}
	if me.overflow > 0 {
		buf.WriteString("... ")
		buf.WriteString(strconv.Itoa(me.overflow))
		buf.WriteString(" errors truncated.\n")
	}
	return buf.String()
}
