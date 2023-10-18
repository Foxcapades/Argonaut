package argo

import (
	"errors"
	"fmt"
	"strings"
)

// A MultiError is an error instance that contains a collection of sub-errors
// of any type except MultiError.
type MultiError interface {
	error

	// Errors returns a slice of all the errors collected into this MultiError
	// instance.
	Errors() []error

	// AppendError appends the given error to this MultiError.  If the given error
	// is itself a MultiError, it will be unpacked into this MultiError instance.
	AppendError(err error)
}

func newMultiError() MultiError {
	return &multiError{
		errs:    make([]error, 0, 10),
		strings: make(map[string]bool, 10),
	}
}

type multiError struct {
	errs    []error
	strings map[string]bool
}

func (m *multiError) Error() string {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("encountered %d unique errors:", len(m.strings)))

	for k := range m.strings {
		sb.WriteString("\n  ")
		sb.WriteString(k)
	}

	return sb.String()
}

func (m *multiError) Errors() []error {
	return m.errs
}

func (m *multiError) AppendError(err error) {
	var e MultiError
	if errors.As(err, &e) {
		for _, err := range e.Errors() {
			m.AppendError(err)
		}
	} else {
		m.errs = append(m.errs, err)
		m.strings[err.Error()] = true
	}
}
