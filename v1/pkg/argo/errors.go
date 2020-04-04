package argo

import (
	"fmt"
	"reflect"
)

//┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓//
//┃                                                                          ┃//
//┃      Error Hints                                                         ┃//
//┃                                                                          ┃//
//┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛//

// ErrorHint is an explanation of a returned error that
// should provide a hint at how the issue can be resolved.
type ErrHint string

const (
	ErrHintShortInvalidChar = "Short flags must be alphanumeric"
	ErrHintLongInvalidChar  = `Long flags may only contain alphanumeric characters, underscores, and dashes`
	ErrHintNoFlag           = `Flags must have a short and/or long flag set`
)

//┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓//
//┃                                                                          ┃//
//┃      Error Types                                                         ┃//
//┃                                                                          ┃//
//┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛//

type InvalidFlagCharError struct {
	Flag string
	Hint ErrHint
}

func (i *InvalidFlagCharError) Error() string {
	return fmt.Sprintf(errInvalidFlag, i.Flag, i.Hint)
}

//┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓//
//┃     Invalid Flag Config     ┃//
//┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛//

type InvalidFlagError struct {
	Hint ErrHint
}

func (i *InvalidFlagError) Error() string {
	return fmt.Sprintf(errCannotBuildFlag, i.Hint)
}

//┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓//
//┃     Nil or Non-Ptr          ┃//
//┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛//

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

//┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓//
//┃     Invalid Type            ┃//
//┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛//

type InvalidTypeError struct {
	Value reflect.Value
}

func (i *InvalidTypeError) Error() string {
	return fmt.Sprintf("Cannot unmarshal type %s", i.Value.Type())
}

//┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓//
//┃     Invalid Format          ┃//
//┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛//

const errFormat = "Format error in input text, could not unmarshal to type %s"

type FormatError struct {
	Value    reflect.Value
	Argument string
	Kind     reflect.Kind
	Root     error
}

func (f *FormatError) Error() string {
	return fmt.Sprintf(errFormat, f.Kind)
}

//┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓//
//┃                                                                          ┃//
//┃      Internals                                                           ┃//
//┃                                                                          ┃//
//┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛//

const (
	errInvalidFlag     = `Cannot set flag key "%s": %s`
	errCannotBuildFlag = `Cannot build flag: %s`
)

func newInvalidLongFlagErr(val string, hint ErrHint) error {
	return &InvalidFlagCharError{Flag: val, Hint: hint}
}

func newInvalidShortFlagErr(val byte, hint ErrHint) error {
	return &InvalidFlagCharError{Flag: string([]byte{val}), Hint: hint}
}

func newFlagBuildErr(hint ErrHint) error {
	return &InvalidFlagError{Hint: hint}
}
