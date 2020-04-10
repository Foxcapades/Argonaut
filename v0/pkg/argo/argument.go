package argo

import (
	"fmt"
	"reflect"
)

// Argument represents a raw argument value for a flag or
// command.
//
// An argument is any string passed to a command that is not
// a flag or subcommand.
type Argument interface {
	fmt.Stringer
	described
	named

	// Default returns the default value assigned to this
	// argument.
	Default() interface{}

	// HasDefault returns whether or not this argument has a
	// default value assigned.
	//
	// When parsing, if an argument has no default value and
	// is not present in the command line, no attempt to
	// assign a value to a binding will be made.
	HasDefault() bool

	// DefaultType returns the reflected type of this
	// argument's default value, or nil if this argument does
	// not have a default value assigned.
	DefaultType() reflect.Type

	// RawValue returns the raw value this argument was
	// assigned when parsing the CLI input.
	//
	// If this argument was built but not parsed, this will
	// return an empty string.
	RawValue() string

	// Required returns whether or not this argument has been
	// marked as required.
	Required() bool

	// SetRawValue sets the input value from the CLI that has
	// been matched to this argument by flag or position.
	SetRawValue(string)

	// Binding returns the pointer that was assigned to this
	// argument into which the CLI input should be
	// unmarshaled.
	Binding() interface{}

	// HasBinding returns whether or not this argument has a
	// binding attached to it.
	HasBinding() bool

	// BindingType returns the root type of the pointer that
	// was bound to this Argument.  If no binding is present
	// this method returns nil.
	BindingType() reflect.Type

	// Parent returns the parent Flag or Command that contains
	// this argument.
	Parent() interface{}

	// IsFlagArg returns whether or not the parent of this
	// Argument is a Flag.
	IsFlagArg() bool

	// IsPositionalArg returns whether or not the parent of
	// this Argument is a Command.
	IsPositionalArg() bool
}
