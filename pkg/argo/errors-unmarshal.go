package argo

import (
	"fmt"
	"reflect"
)

const (
	errFormat = "Format error in input text, could not unmarshal to type %s"
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
	return fmt.Sprintf(errFormat, f.Kind)
}

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
