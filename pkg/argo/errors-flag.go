package argo

import (
	"fmt"
)

type FlagBindingError interface {
	error

	ArgumentBuilder() ArgumentBuilder
	FlagBuilder() FlagBuilder

	Unwrap() error
}

func newFlagBindingError(root ArgumentBindingError, arg ArgumentBuilder, flag FlagBuilder) FlagBindingError {
	return &flagBindingError{root.Unwrap(), arg, flag}
}

type flagBindingError struct {
	root error
	arg  ArgumentBuilder
	flag FlagBuilder
}

func (f flagBindingError) Error() string {
	if f.flag.hasLongForm() {
		if f.flag.hasShortForm() {
			return "FlagBindingError (--" + f.flag.getLongForm() + "|-" + string([]byte{f.flag.getShortForm()}) + "): " + f.root.Error()
		} else {
			return "FlagBindingError (--" + f.flag.getLongForm() + "): " + f.root.Error()
		}
	} else if f.flag.hasShortForm() {
		return "FlagBindingError (-" + string([]byte{f.flag.getShortForm()}) + "): " + f.root.Error()
	} else {
		return "FlagBindingError: " + f.root.Error()
	}
}

func (f flagBindingError) ArgumentBuilder() ArgumentBuilder {
	return f.arg
}

func (f flagBindingError) FlagBuilder() FlagBuilder {
	return f.flag
}

func (f flagBindingError) Unwrap() error {
	return f.root
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

func newMissingFlagError(flag Flag) MissingFlagError {
	return missingFlagError{flag}
}

type missingFlagError struct {
	flag Flag
}

func (m missingFlagError) Error() string {
	return fmt.Sprintf("required flag %s was missing from the CLI call", printFlagNames(m.flag))
}

func (m missingFlagError) Flag() Flag {
	return m.flag
}
