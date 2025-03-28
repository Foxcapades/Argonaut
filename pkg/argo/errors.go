package argo

import (
	"reflect"
)

// A MultiError is an error instance that contains a collection of sub-errors
// of any type except MultiError.
type MultiError interface {
	error

	// Errors returns a slice of all the distinct errors collected into this
	// MultiError instance.
	Errors() []error

	// AppendError appends the given error to this MultiError.  If the given error
	// is itself a MultiError, it will be unpacked into this MultiError instance.
	AppendError(err error)
}

type UnmarshallingError interface {
	error

	Value() reflect.Value
}

type FormatError interface {
	UnmarshallingError

	Argument() string

	Kind() reflect.Kind

	Root() error
}
