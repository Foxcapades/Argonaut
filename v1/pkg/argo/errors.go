package argo

import (
	"fmt"
	"reflect"
)

// ┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓ //
// ┃     Nil or Non-Ptr          ┃ //
// ┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛ //

type InvalidUnmarshalError struct {
	Value    reflect.Value
	Argument string
}

func (i *InvalidUnmarshalError) Error() string {
	if i.Value.IsNil() {
		return "Attempted to unmarshal into nil"
	}
	return "Attempted to unmarshal into a non-pointer"
}

// ┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓ //
// ┃     Invalid Type            ┃ //
// ┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛ //

type InvalidTypeError struct {
	Value reflect.Value
}

func (i *InvalidTypeError) Error() string {
	return fmt.Sprintf("Cannot unmarshal type %s", i.Value.Type())
}

// ┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓ //
// ┃     Invalid Format          ┃ //
// ┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛ //

type FormatError struct {
	Value    reflect.Value
	Argument string
	Kind     reflect.Kind
	Root     error
}

func (f *FormatError) Error() string {
	return fmt.Sprintf(errFormat, f.Kind)
}

// ┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓ //
// ┃     Missing Argument        ┃ //
// ┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛ //

type MissingRequiredArgumentError interface {
	error
	Argument() Argument
	Flag() Flag
	HasFlag() bool
}

type UnexpectedArgumentError interface {
	error
	RawValue() string
	Flag() Flag
	HasFlag() bool
}

// ////////////////////////////////////////////////////////////////////////// //
//                                                                            //
//    Missing Flag Error                                                      //
//                                                                            //
// ////////////////////////////////////////////////////////////////////////// //

// A MissingFlagError is returned on CLI parse when a flag that has been marked
// as being required was not found to be present in the CLI call.
//
// MissingFlagError is a hard error that will be returned regardless of whether
// the parser is operating in strict mode.
type MissingFlagError interface {
	error
	Flag() Flag
}

// ////////////////////////////////////////////////////////////////////////// //
//                                                                            //
//    Missing Argument Error                                                  //
//                                                                            //
// ////////////////////////////////////////////////////////////////////////// //

// An IncompleteCommandError is returned when the CLI call did not reach a
// command leaf.
type IncompleteCommandError interface {
	error
	LastReached() CommandNode
}

// ┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓ //
// ┃                                                                        ┃ //
// ┃     Internals                                                          ┃ //
// ┃                                                                        ┃ //
// ┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛ //

const (
	errInvalidFlag     = `Cannot set flag key "%s": %s`
	errCannotBuildFlag = `Cannot build flag: %s`
	errFormat          = "Format error in input text, could not unmarshal to type %s"
)
