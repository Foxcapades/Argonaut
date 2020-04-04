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

type InvalidFlagError struct {
	Hint ErrHint
}

func (i *InvalidFlagError) Error() string {
	return fmt.Sprintf(errCannotBuildFlag, i.Hint)
}

type InvalidUnmarshalError struct {
	value    reflect.Value
	Argument string
}

func (i *InvalidUnmarshalError) Error() string {
	if i.value.IsNil() {
		return "Attempted to unmarshal into nil"
	}
	return "Attempted to unmarshal into a non-pointer"
}

type RecursivePointerError struct {
	Chain    []reflect.Value
	Argument string
}

func (r *RecursivePointerError) Error() string {
	return "Attempted to unmarshal into a recursive pointer"
}

func (r *RecursivePointerError) Depth() int {
	return len(r.Chain) - 1
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
