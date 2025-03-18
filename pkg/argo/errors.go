package argo

import (
	"fmt"
	"reflect"
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

const (
	unmarshalErrFormat = "Format error in input text, could not unmarshal to type %s"
)

type UnmarshallingError interface {
	error

	Value() reflect.Value
}

// ┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓ //
// ┃     Invalid Format          ┃ //
// ┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛ //

type FormatError interface {
	UnmarshallingError

	Argument() string

	Kind() reflect.Kind

	Root() error
}

type formatError struct {
	Value    reflect.Value
	Argument string
	Kind     reflect.Kind
	Root     error
}

func (f formatError) Error() string {
	return fmt.Sprintf(unmarshalErrFormat, f.Kind)
}
