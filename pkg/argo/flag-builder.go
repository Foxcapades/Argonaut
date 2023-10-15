package argo

type FlagBuilder interface {

	// WithShortForm sets the short-form flag character that the flag may be
	// referenced by on the CLI.
	//
	// Short-form flags consist of a single character preceded by either a
	// prefix/leader character, or by one or more other short-form flags.
	//
	// Short-form flags must be alphanumeric.
	//
	// Examples:
	//     # Single, unchained short flags.
	//     -f
	//     -f bar
	//     -f=bar
	//     # Multiple short flags chained.  In these examples, the short flag '-c'
	//     # takes an optional string argument, which will be "def" in the
	//     # examples.
	//     -abc
	//     -abc def
	//     -abc=def
	WithShortForm(char byte) FlagBuilder

	HasShortForm() bool

	GetShortForm() byte

	// WithLongForm sets the long-form flag name that the flag may be referenced
	// by on the CLI.
	//
	// Long-form flags consist of one or more characters preceded immediately by
	// two prefix/leader characters (typically dashes).
	//
	// Long-form flags must start with an alphanumeric character and may only
	// consist of alphanumeric characters, dashes, and/or underscores.
	//
	// Example long-form flags:
	//     # The '--foo' flag takes an optional string argument
	//     --foo
	//     --foo bar
	//     --foo=bar
	WithLongForm(form string) FlagBuilder

	HasLongForm() bool

	GetLongForm() string

	// WithDescription sets an optional description value for the Flag being
	// built.
	//
	// The description value is used for rendering help text.
	WithDescription(desc string) FlagBuilder

	HasDescription() bool

	// WithOnHitCallback provides a function that will be called when a Flag is
	// hit while parsing the CLI inputs.
	//
	// The given function will be called after parsing has completed, regardless
	// of whether there were parsing errors.
	//
	// Flag on-hit callbacks will be executed in the order that the flags appeared
	// in the CLI call.
	WithOnHitCallback(fn func(Flag)) FlagBuilder

	// WithArgument attaches the given argument to the Flag being built.
	//
	// Only one argument may be set on a Flag at a time.
	WithArgument(arg ArgumentBuilder) FlagBuilder

	// WithBinding is a shortcut method for attaching an argument and binding it
	// to the given pointer.
	//
	// Bind is equivalent to calling one of the following:
	//    Argument(NewArgument().Bind(ptr))
	//    // or
	//    Argument(NewArgument().Bind(ptr).Require())
	WithBinding(pointer any, required bool) FlagBuilder

	// WithBindingAndDefault is a shortcut method for attaching an argument,
	// binding it to the given pointer, and setting a default on that argument.
	//
	// BindWithDefault is equivalent to calling one of the following:
	//     WithArgument(Argument().WithBinding(ptr).WithDefault(something))
	//     // or
	//     WithArgument(Argument().WithBinding(ptr).WithDefault(something).Require())
	WithBindingAndDefault(pointer, def any, required bool) FlagBuilder

	// Require marks this Flag as being required.
	//
	// If this flag is not present in the CLI call, an error will be returned when
	// parsing the CLI input.
	Require() FlagBuilder

	// Build builds a new Flag instance constructed from the components set on
	// this FlagBuilder.
	Build() (Flag, error)
}
