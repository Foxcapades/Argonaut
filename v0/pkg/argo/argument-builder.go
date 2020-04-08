package argo

type ArgumentBuilder interface {
	Name(string) ArgumentBuilder

	// TypeHint sets the hint value for this argument.
	//
	// Used when rendering a usage example.
	TypeHint(string) ArgumentBuilder

	// Default sets the default value for the argument to be
	// used if the argument is not provided on the command
	// line.
	//
	// The type of this value must match the type of the value
	// used with `Bind()`
	//
	// If the value provided to this method is a pointer to
	// the type of the bind value it will be dereferenced to
	// set the bind value.
	Default(interface{}) ArgumentBuilder

	GetDefault() interface{}

	HasDefault() bool

	HasDefaultProvider() bool

	// Bind sets the pointer into which the value will be.
	// parsed into on parse.
	//
	// The type of this pointer must match the type of the
	// value provided with `Default()` if default values are
	// used.
	Bind(ptr interface{}) ArgumentBuilder

	GetBinding() interface{}

	HasBinding() bool

	// Description sets the description of this argument to be
	// shown in rendered help text.
	Description(string) ArgumentBuilder

	// Require marks this argument as required.
	Require() ArgumentBuilder

	// Required sets whether or not this argument is required
	// based on the provided value.
	Required(bool) ArgumentBuilder

	Build() (Argument, error)

	MustBuild() Argument

	Parent(interface{}) ArgumentBuilder
}
