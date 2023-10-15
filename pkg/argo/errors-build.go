package argo

// A MultiError is an error instance that contains a collection of sub-errors
// of any type except MultiError.
type MultiError interface {
	error

	// Errors returns a slice of all the errors collected into this MultiError
	// instance.
	Errors() []error

	// NonStrictErrors returns a slice of only the errors that apply regardless of
	// strict mode.
	NonStrictErrors() []error

	// AppendError appends the given error to this MultiError.  If the given error
	// is itself a MultiError, it will be unpacked into this MultiError instance.
	AppendError(err error)
}
