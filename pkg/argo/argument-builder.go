package argo

// An ArgumentBuilder instance is used to construct a CLI argument that may be
// attached to a Flag or LeafCommand.
type ArgumentBuilder interface {

	// WithName sets the name for this argument.
	//
	// The name value is used when rendering help information about this argument.
	WithName(name string) ArgumentBuilder

	HasName() bool

	Name() string

	// WithDescription sets the description of this argument to be shown in
	// rendered help text.
	WithDescription(desc string) ArgumentBuilder

	HasDescription() bool

	Description() string

	// Require marks the output Argument as being required.
	Require() ArgumentBuilder

	IsRequired() bool

	// WithBinding sets the bind value for the Argument.
	//
	// The bind value may be one of a value pointer, a consumer function, or an
	// Unmarshaler instance.  For demonstrations of each, see the examples below.
	//
	// If the bind value is a pointer, the Argument's value unmarshaler will be
	// called to unmarshal the raw string value into a value of the type passed
	// to this method.
	//
	// If the bind value is a consumer function, that function will be called with
	// the parsed value from the CLI.  The consumer function may optionally return
	// an error which, if not nil, will be passed up as a parsing error.
	//
	// If the bind value is an Unmarshaler instance, that instance's Unmarshal
	// method will be called with the raw input from the CLI.
	//
	// Setting this value to anything other than a pointer or an Unmarshaler
	// instance will result in an error being returned when building the argument
	// is attempted.
	//
	// Example 1 (a simple var binding):
	//     var myValue time.Duration
	//     cli.Argument.WithBinding(&myValue)
	//
	// Example 2 (an unmarshaler func):
	//     cli.Argument.WithBinding(UnmarshalerFunc(func(raw string) error {
	//         fmt.Println(raw)
	//         return nil
	//     }))
	//
	// Example 3 (lets get silly with it):
	//     var myValue map[bool]**string
	//     cli.Argument().WithBinding(&myValue)
	//
	// Example 4 (custom type)
	//     type Foo struct {
	//         // some fields
	//     }
	//
	//     func (f *Foo) Unmarshal(raw string) error {
	//         // parse the given string
	//         return nil
	//     }
	//
	//     func main() {
	//         var foo Foo
	//         cli.Argument().WithBinding(&foo)
	//     }
	//
	// Example 5 (plain consumer func which returns an error):
	//     cli.Argument().WithBinding(func(value int) error { do something })
	//
	// Example 6 (plain consumer func which returns nothing):
	//     cli.Argument().WithBinding(func(value int) { do something })
	//
	WithBinding(pointer any) ArgumentBuilder

	HasBinding() bool

	Binding() ArgumentBinding

	// WithDefault sets the default value for the argument to be used if the
	// argument is not provided on the command line.
	//
	// Setting this value without providing a binding value using `Bind()` will
	// mean that the given default will not be set to anything when the CLI input
	// is parsed.
	//
	// When used, the type of this value must meet one of the following criteria:
	//   1. `val` is compatible with the type of the value used with
	//      WithBinding.
	//   2. `val` is a string that may be parsed into a value of the type used
	//      with WithBinding.
	//   3. `val` is a function which returns a type that is compatible with the
	//      type of the value used with WithBinding
	//   4. `val` is a function which returns a type that is compatible with the
	//      type of the value used with WithBinding in addition to returning an
	//      error as the second return value.
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

	Default() ArgumentDefault

	// WithUnmarshaler allows providing a custom ValueUnmarshaler instance that
	// will be used to unmarshal string values into the binding type.
	//
	// If no binding is set on this argument, the provided ValueUnmarshaler will
	// not be used.
	//
	// If a custom unmarshaler is not provided by way of this method, then the
	// internal magic unmarshaler will be used to parse raw argument values.
	WithUnmarshaler(fn ValueUnmarshaler) ArgumentBuilder

	HasUnmarshaler() bool

	Unmarshaler() ValueUnmarshaler

	// WithValidator appends the given validator function to the argument's
	// internal slice of validators.
	//
	// There are 2 types of validators that may be set here, each of which going
	// to a separate slice.  Type 1 is a pre-parse validator which will be called
	// when an argument is first hit, but before it is parsed.  Type 2 is a
	// post-parse validator which will be called immediately after an argument is
	// parsed to validate the parsed value.
	//
	// When appending a validator function, if it is of type 1 it will go to the
	// pre-parse validator slice, and if it is of type 2 it will go to the
	// post-parse validator slice.
	//
	// Pre-parse (type 1) validators must match the following function signature:
	//     func(string) error
	//
	// The value that is passed to the function will be the raw value that was
	// passed to the command on the CLI.  If an error is returned, CLI parsing
	// will halt, and the returned error will be passed up.
	//
	// Post-parse (type 2) validators must match the following function signature:
	//     func(T, string) error
	//
	// Two values are passed to the function, the parsed value, and the raw value
	// that was passed to the command ont he CLI.  If an error is returned, CLI
	// parsing will halt, and the returned error will be passed up.
	//
	// Validators will be executed in the order they are appended.
	WithValidator(validatorFn any) ArgumentBuilder

	HasValidators() bool

	Validators() []any
}
