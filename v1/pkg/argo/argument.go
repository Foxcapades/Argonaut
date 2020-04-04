package argo

import "reflect"

// Argument represents a raw argument value for a flag or
// command.
//
// An argument is any string passed to a command that is not
// a flag or subcommand.
type Argument interface {
	// Hint returns the hint text for this argument.
	//
	// Hint text is used when rendering help as a stand in for
	// this argument in examples.
	Hint() string

	// HasHint returns whether or not this argument has a
	// hint value.
	//
	// When rendering help, if an argument does not have a
	// hint value, it's type will be used to render examples
	// in generated help text.
	HasHint() bool

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

	// Description returns the description text assigned to
	// this argument.
	//
	// Description text is used when rendering help text.
	Description() string

	// HasDescription returns whether or not this argument has
	// a description assigned.
	HasDescription() bool

	RawValue() string

	Required() bool
}
