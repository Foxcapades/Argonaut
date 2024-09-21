package argo

import (
	"fmt"
)

type ArgumentBindingError interface {
	error

	Builder() ArgumentBuilder

	Unwrap() error
}

func newArgumentBindingError(root error, builder ArgumentBuilder) ArgumentBindingError {
	return &argumentBindingError{root, builder}
}

type argumentBindingError struct {
	root    error
	builder ArgumentBuilder
}

func (a argumentBindingError) Unwrap() error {
	return a.root
}

func (a argumentBindingError) Error() string {
	return "ArgumentBindingError: " + a.root.Error()
}

func (a argumentBindingError) Builder() ArgumentBuilder {
	return a.builder
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
