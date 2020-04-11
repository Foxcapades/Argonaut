package argo

import "fmt"

type FlagBuilderErrorType uint8

const (
	FlagBuilderErrNoFlags FlagBuilderErrorType = iota
	FlagBuilderErrBadShortFlag
	FlagBuilderErrBadLongFlag
)

func (f FlagBuilderErrorType) String() string {
	switch f {
	case FlagBuilderErrNoFlags:
		return "No flag pattern provided"
	case FlagBuilderErrBadShortFlag:
		return "Invalid short flag character"
	case FlagBuilderErrBadLongFlag:
		return "Invalid character in long flag"
	}
	panic("Invalid flag error type")
}

// FlagBuilderError represents a configuration error
// encountered while attempting to build a Flag
type FlagBuilderError interface {
	error

	// Type returns the specific type of error encountered.
	Type() FlagBuilderErrorType

	// FlagBuilder returns the FlagBuilder instance in which
	// the error was encountered
	FlagBuilder() FlagBuilder
}

func NewFlagBuilderError(errType FlagBuilderErrorType, builder FlagBuilder) error {
	return &flagBuilderError{eType: errType, builder: builder}
}

type flagBuilderError struct {
	eType   FlagBuilderErrorType
	builder FlagBuilder
}

func (i *flagBuilderError) Type() FlagBuilderErrorType {
	return i.eType
}

func (i *flagBuilderError) FlagBuilder() FlagBuilder {
	return i.builder
}

func (i *flagBuilderError) Error() string {
	switch i.eType {
	case FlagBuilderErrNoFlags:
		return `Flags must have a short and/or long flag set`
	case FlagBuilderErrBadLongFlag:
		return `Long flags may only contain alphanumeric characters, underscores, and dashes`
	case FlagBuilderErrBadShortFlag:
		return `Short flags must be alphanumeric`
	}
	panic(fmt.Errorf("invalid flag error type %d", i.eType))
}
