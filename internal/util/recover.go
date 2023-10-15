package util

import "errors"

// Recovery catches panics and set the panic value to the given error pointer if
// the panic value is an instance of error or string.
func Recovery(err *error) {
	rec := recover()

	// No panics, nothing to do
	if rec == nil {
		return
	}

	// If the panic was due to an error, pass it up and  return.
	if tmp, ok := rec.(error); ok {
		*err = tmp
		return
	}

	// If the panic was a string, convert it to an error, pass it up and return.
	if tmp, ok := rec.(string); ok {
		*err = errors.New(tmp)
		return
	}

	// If the panic was some unknown type, it didn't come from us, panic with it
	// again.
	panic(rec)
}
