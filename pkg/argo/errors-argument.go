package argo

import (
	"fmt"
	"reflect"
)

// ArgumentErrorType provides a flag indicating what kind
// of error was encountered while attempting to build or use
// an Argument.
type ArgumentErrorType uint8

const (
	// ArgErrInvalidDefault represents all cases where the
	// Default value passed to an ArgumentBuilder is the root
	// cause.
	ArgErrInvalidDefault ArgumentErrorType = 1 << iota

	// ArgErrInvalidBinding represents all cases where the
	// Binding value passed to an ArgumentBuilder is the root
	// cause.
	ArgErrInvalidBinding

	// ArgErrInvalidDefaultFn represents the case where the
	// default value provider function set on an
	// ArgumentBuilder is incompatible with the binding type
	// also set on that ArgumentBuilder.
	ArgErrInvalidDefaultFn = ArgErrInvalidDefault | 1<<iota

	// ArgErrInvalidDefaultFn represents the case where the
	// default value set on an ArgumentBuilder is incompatible
	// with the binding type also set on that ArgumentBuilder.
	ArgErrInvalidDefaultVal = ArgErrInvalidDefault | 1<<iota

	// ArgErrInvalidBindingNil <currently not used>
	// represents the case where the binding set on an
	// ArgumentBuilder is an untyped nil.
	// ArgErrInvalidBindingNil     = ArgErrInvalidBinding | 1<<iota

	// ArgErrInvalidBindingBadType represents the case where
	// the value set as the ArgumentBuilder's binding is not
	// of a type that can be unmarshaled.
	ArgErrInvalidBindingBadType = ArgErrInvalidBinding | 1<<iota
)

func (a ArgumentErrorType) String() string {
	switch a {
	case ArgErrInvalidDefault:
		return "Invalid default"
	case ArgErrInvalidBinding:
		return "Invalid binding"
	case ArgErrInvalidDefaultFn:
		return "Invalid default value provider"
	case ArgErrInvalidDefaultVal:
		return "Invalid default value"
	case ArgErrInvalidBindingBadType:
		return "Invalid binding type, not unmarshalable"
	}
	panic("Invalid ArgumentErrorType")
}

// ArgumentError represents an error encountered when attempting to build or
// handle an Argument.
type ArgumentError interface {
	error

	// Type returns the specific type of this error.
	//
	// See ArgumentErrorType for more details.
	Type() ArgumentErrorType

	// Is returns whether this error is of the type given.
	Is(errorType ArgumentErrorType) bool

	// Builder returns the ArgumentBuilder in which this error was encountered.
	Builder() ArgumentBuilder
}

// ┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓ //
// ┃     Invalid Arg Config      ┃ //
// ┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛ //

// Error message formats
const (
	errArgBinding = `Argument bound to type "%s" which cannot be unmarshaled.`
	errArgDefault = "Argument default type (%s) must be compatible with binding" +
		" type (%s)"
	errArgDefFn = "Invalid default provider func (%s): %s"
)

func newInvalidArgError(
	kind ArgumentErrorType,
	build ArgumentBuilder,
	reason string,
) error {
	return &InvalidArgError{eType: kind, build: build, reason: reason}
}

type InvalidArgError struct {
	eType  ArgumentErrorType
	build  ArgumentBuilder
	reason string
}

func (i *InvalidArgError) Type() ArgumentErrorType {
	return i.eType
}

func (i *InvalidArgError) Is(kind ArgumentErrorType) bool {
	return i.eType&kind == kind
}

func (i *InvalidArgError) Builder() ArgumentBuilder {
	return i.build
}

func (i *InvalidArgError) Error() string {
	switch i.eType {
	case ArgErrInvalidBindingBadType /*, ArgErrInvalidBindingNil*/ :
		return fmt.Sprintf(errArgBinding, reflect.TypeOf(i.build.getBinding()))
	case ArgErrInvalidDefaultVal:
		return fmt.Sprintf(errArgDefault, reflect.TypeOf(i.build.getDefault()),
			reflect.TypeOf(i.build.getBinding()))
	case ArgErrInvalidDefaultFn:
		return fmt.Sprintf(errArgDefFn, reflect.TypeOf(i.build.getDefault()), i.reason)
	}
	panic(fmt.Errorf("invalid argument error type %d", i.eType))
}

// ////////////////////////////////////////////////////////////////////////// //
//                                                                            //
//    Missing Argument Error                                                  //
//                                                                            //
// ////////////////////////////////////////////////////////////////////////// //

type MissingRequiredArgumentError interface {
	error
	Argument() Argument
	Flag() Flag
	HasFlag() bool
}

func newMissingRequiredPositionalArgumentError(a Argument, c Command) MissingRequiredArgumentError {
	return &missingArgError{arg: a, com: c}
}

func newMissingRequiredFlagArgumentError(a Argument, f Flag, c Command) MissingRequiredArgumentError {
	return &missingArgError{arg: a, flag: f, com: c}
}

type missingArgError struct {
	arg  Argument
	flag Flag
	com  Command
}

func (m *missingArgError) Flag() Flag {
	return m.flag
}

func (m *missingArgError) HasFlag() bool {
	return m.flag != nil
}

func (m *missingArgError) StrictOnly() bool {
	return false
}

func (m *missingArgError) Error() string {
	if m.flag != nil {
		return fmt.Sprintf("Missing required argument for flag %s", printFlagNames(m.flag))
	} else if m.arg.HasName() {
		return fmt.Sprintf("Missing required positional argument %s", m.arg.Name())
	} else {
		for i, a := range m.com.Arguments() {
			if m.arg == a {
				return fmt.Sprintf("Missing required positional argument #%d", i+1)
			}
		}
		return "Missing required positional argument"
	}
}

func (m *missingArgError) Argument() Argument {
	return m.arg
}
