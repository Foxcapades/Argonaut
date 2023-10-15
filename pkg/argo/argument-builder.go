package argo

// An ArgumentBuilder instance is used to construct a CLI argument that may be
// attached to a Flag or CommandLeaf.
type ArgumentBuilder interface {

	// WithName sets the name for this argument.
	//
	// The name value is used when rendering help information about this argument.
	WithName(name string) ArgumentBuilder

	HasName() bool

	GetName() string

	// WithDescription sets the description of this argument to be shown in rendered
	// help text.
	WithDescription(desc string) ArgumentBuilder

	HasDescription() bool

	GetDescription() string

	// Require marks the output Argument as being required.
	Require() ArgumentBuilder

	IsRequired() bool

	// WithBinding sets the pointer into which the value will be parsed into on
	// parse.
	//
	// The type of this pointer must match the type of the value provided with
	// `Default()` if default values are used.
	//     foo := myStruct{}
	//
	//     cli.Argument().
	//         WithBinding(&foo.bar)
	WithBinding(pointer any) ArgumentBuilder

	HasBinding() bool

	// GetBinding returns the binding value that was set on this ArgumentBuilder.
	GetBinding() any

	// WithDefault sets the default value for the argument to be used if the
	// argument is not provided on the command line.
	//
	// Setting this value without providing a binding value using `Bind()` will
	// mean that the given default will not be set to anything when the CLI input
	// is parsed.
	//
	// When used, the type of this value must meet one of the following criteria:
	//   1. `val` is compatible with the type of the value used with `Bind()`
	//   2. `val` is a function which returns a type that is compatible with the
	//      type of the value used with `Bind()`
	//   3. `val` is a function which returns a type that is compatible with the
	//      type of the value used with `Bind()` in addition to returning an error
	//      as the second return value.
	//
	// Examples:
	//     arg.WithBinding(&fooString).WithDefault(3)   // Type mismatch
	//
	//     arg.WithBinding(&fooInt).WithDefault(3)      // OK
	//
	//     arg.WithBinding(&fooInt).
	//       WithDefault(func() int {return 3})         // OK
	//
	//     arg.WithBinding(&fooInt).
	//       WithDefault(func() (int, error) {
	//         return 3, nil
	//       })                                         // OK
	//
	// If the value provided to this method is a pointer to the type of the bind
	// value it will be dereferenced to set the bind value.
	WithDefault(def any) ArgumentBuilder

	HasDefault() bool

	// GetDefault returns the default value, if any, that was set on this
	// ArgumentBuilder instance.
	GetDefault() any

	// WithUnmarshaler allows providing a custom ValueUnmarshaler instance that
	// will be used to unmarshal string values into the binding type.
	//
	// If no binding is set on this argument, the provided ValueUnmarshaler will
	// not be used.
	//
	// If a custom unmarshaler is not provided by way of this method, then the
	// internal magic unmarshaler will be used to parse raw arguments.
	WithUnmarshaler(fn ValueUnmarshaler) ArgumentBuilder

	// Build constructs an Argument instance from the parameters set on this
	// ArgumentBuilder.
	Build() (Argument, error)
}
