package argo

import (
	"fmt"
	R "reflect"
)

//┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓//
//┃     Invalid Flag Config     ┃//
//┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛//

type InvalidFlagErrorType uint8

const (
	InvalidFlagNoFlags InvalidFlagErrorType = iota
	InvalidFlagBadShortFlag
	InvalidFlagBadLongFlag
)

type InvalidFlagError interface {
	error
	Type() InvalidFlagErrorType
}

func NewInvalidFlagError(errType InvalidFlagErrorType) error {
	return &invalidFlagError{eType: errType}
}

type invalidFlagError struct {
	eType InvalidFlagErrorType
}

func (i *invalidFlagError) Type() InvalidFlagErrorType {
	return i.eType
}

func (i *invalidFlagError) Error() string {
	switch i.eType {
	case InvalidFlagNoFlags:
		return `Flags must have a short and/or long flag set`
	case InvalidFlagBadLongFlag:
		return `Long flags may only contain alphanumeric characters, underscores, and dashes`
	case InvalidFlagBadShortFlag:
		return `Short flags must be alphanumeric`
	}
	panic(fmt.Errorf("invalid flag error type %d", i.eType))
}

//┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓//
//┃     Nil or Non-Ptr          ┃//
//┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛//

type InvalidUnmarshalError struct {
	Value    R.Value
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
	Value R.Value
}

func (i *InvalidTypeError) Error() string {
	return fmt.Sprintf("Cannot unmarshal type %s", i.Value.Type())
}

//┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓//
//┃     Invalid Format          ┃//
//┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛//

type FormatError struct {
	Value    R.Value
	Argument string
	Kind     R.Kind
	Root     error
}

func (f *FormatError) Error() string {
	return fmt.Sprintf(errFormat, f.Kind)
}

//┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓//
//┃     Missing Argument        ┃//
//┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛//

type MissingRequiredArgumentError interface {
	error
	Argument() Argument
	Flag() Flag
	HasFlag() bool
}

func NewMissingRequiredPositionalArgumentError(a Argument, c Command) error {
	return &missingArgError{arg: a, com: c}
}

func NewMissingRequiredFlagArgumentError(a Argument, f Flag, c Command) error {
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

func (m *missingArgError) Error() string {
	if m.flag != nil {
		return fmt.Sprintf("Missing required argument for flag %s", printFlag(m.flag))
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
	panic("implement me")
}

//┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓//
//┃                                                                          ┃//
//┃      Internals                                                           ┃//
//┃                                                                          ┃//
//┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛//

const (
	errInvalidFlag     = `Cannot set flag key "%s": %s`
	errCannotBuildFlag = `Cannot build flag: %s`
	errFormat          = "Format error in input text, could not unmarshal to type %s"
)

func printFlag(f Flag) string {
	if f.HasLong() {
		if f.HasShort() {
			return fmt.Sprintf("--%s / -%c", f.Long(), f.Short())
		}
		return fmt.Sprintf("--%s", f.Long())
	}
	return fmt.Sprintf("%c", f.Short())
}
