package argo

type ArgumentBuilder interface {

	// Name sets the name for this argument.
	//
	// The name value is used when rendering help / manpage
	// information about this argument.
	Name(string) ArgumentBuilder

	// Default sets the default value for the argument to be
	// used if the argument is not provided on the command
	// line.
	//
	// Setting this value without providing a binding value
	// using `Bind()` will cause an error when Build is
	// called.
	//
	// When used, the type of this value must meet one of the
	// following criteria:
	//   1. `val` is compatible with the type of the value
	//      used with `Bind()`
	//   2. `val` is a function which returns a type that is
	//      compatible with the type of the value used with
	//      `Bind()`
	//   3. `val` is a function which returns a type that is
	//      compatible with the type of the value used with
	//      `Bind()` in addition to returning an error as the
	//      second return value.
	//
	// Examples:
	//     arg.Bind(fooString).Default(3)   // Type mismatch
	//
	//     arg.Bind(fooInt).Default(3)      // OK
	//
	//     arg.Bind(fooInt).
	//       Default(func() int {return 3}) // OK
	//
	//     arg.Bind(fooInt).
	//       Default(func() (int, error) {
	//         return 3, nil
	//       })                             // OK
	//
	// If the value provided to this method is a pointer to
	// the type of the bind value it will be dereferenced to
	// set the bind value.
	Default(val interface{}) ArgumentBuilder

	// GetDefault returns the default value assigned to this
	// Argument that will be used to populate the built
	// Argument if no CLI input is provided for it.
	//
	// If a default value has not been assigned to this
	// builder then this method will return nil.
	GetDefault() interface{}

	// HasDefault returns whether or not a default value has
	// been set on this builder.
	HasDefault() bool

	// HasDefaultProvider returns whether or not a default
	// value exists that is of a function type.
	HasDefaultProvider() bool

	// Bind sets the pointer into which the value will be.
	// parsed into on parse.
	//
	// The type of this pointer must match the type of the
	// value provided with `Default()` if default values are
	// used.
	Bind(ptr interface{}) ArgumentBuilder

	// GetBinding returns the pointer into which the CLI input
	// will be unmarshaled.
	//
	// If no binding is set, returns nil
	GetBinding() interface{}

	// HasBinding returns whether or not this builder has a
	// bind pointer set.
	HasBinding() bool

	// Description sets the description of this argument to be
	// shown in rendered help text.
	Description(string) ArgumentBuilder

	// Require marks this argument as required.
	Require() ArgumentBuilder

	// Required sets whether or not this argument is required
	// based on the provided value.
	Required(bool) ArgumentBuilder

	IsRequired() bool

	// Parent sets the parent element of type Command or Flag
	// to this builder.
	//
	// This method is intended for internal use and parent
	// values set externally may be ignored.
	Parent(interface{}) ArgumentBuilder

	// Build attempts to construct an instance of Argument
	// using the values provided using the various builder
	// methods.
	//
	// If no configuration issues are encountered, Build
	// returns a constructed Argument, else Build returns nil
	// along with the error encountered while validating the
	// configuration.
	Build() (Argument, error)

	// MustBuild calls Build and panics if an error is
	// returned.
	//
	// On success, returns the built Argument.
	MustBuild() Argument
}
