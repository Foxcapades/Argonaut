package argo

// Flag represents a single CLI flag which may have an argument.
type Flag interface {

	// ShortForm returns the short-form character representing this Flag.
	ShortForm() byte

	// HasShortForm indicates whether this Flag has a short-form character.
	HasShortForm() bool

	// LongForm returns the long-form string representing this Flag.
	LongForm() string

	// HasLongForm indicates whether this Flag has a long-form string.
	HasLongForm() bool

	// Description returns the help-text description of this flag.
	Description() string

	// HasDescription indicates whether this Flag has a help-text description.
	HasDescription() bool

	// Argument returns the argument value attached to this Flag.
	//
	// If this Flag does not have an Argument attached, this method will return
	// nil.
	Argument() Argument

	// HasArgument indicates whether this Flag accepts an argument.
	HasArgument() bool

	// IsRequired indicates whether this Flag is required.
	//
	// Required flags must be present in the CLI call.
	IsRequired() bool

	// WasHit indicates whether this flag was used in the CLI call.
	WasHit() bool

	// HitCount returns the number of times that this flag was used in the CLI
	// call.
	HitCount() int

	IncrementHitCount()

	IsHelpFlag() bool

	HasCallback() bool

	Callback() FlagCallback
}

//

// A FlagCallback is a function that, if set on a flag, will be called by the
// CLI parsing process if that flag is used in the CLI call.
//
// The flag callback will be called after CLI parsing has completed.
type FlagCallback = func(flag Flag)

//

// A FlagBuilder is used to construct a Flag instance which represents the input
// from the CLI call.
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
	//     # takes an optional string argument, which will be "def" in the last
	//     # two examples.
	//     -abc
	//     -abc def
	//     -abc=def
	WithShortForm(char byte) FlagBuilder

	HasShortForm() bool

	ShortForm() byte

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

	LongForm() string

	// WithDescription sets an optional description value for the Flag being
	// built.
	//
	// The description value is used for rendering help text.
	WithDescription(desc string) FlagBuilder

	HasDescription() bool

	Description() string

	// WithCallback provides a function that will be called when a Flag is hit
	// while parsing the CLI inputs.
	//
	// The given function will be called after parsing has completed, regardless
	// of whether there were parsing errors.
	//
	// Flag on-hit callbacks will be executed in priority order with the higher
	// priority values executing before lower priority values.  For flags that
	// have the same priority, the callbacks will be called in the order the flags
	// appeared in the CLI call.
	WithCallback(fn FlagCallback) FlagBuilder

	HasCallback() bool

	Callback() FlagCallback

	// WithArgument attaches the given argument to the Flag being built.
	//
	// Only one argument may be set on a Flag at a time.
	WithArgument(arg ArgumentBuilder) FlagBuilder

	HasArgument() bool

	Argument() ArgumentBuilder

	// WithBinding is a shortcut method for attaching an argument and binding it
	// to the given pointer.
	//
	// Bind is equivalent to calling one of the following:
	//    WithArgument(cli.Argument().Bind(ptr))
	//    // or
	//    WithArgument(cli.Argument().Bind(ptr).Require())
	WithBinding(pointer any, required bool) FlagBuilder

	// WithBindingAndDefault is a shortcut method for attaching an argument,
	// binding it to the given pointer, and setting a default on that argument.
	//
	// BindWithDefault is equivalent to calling one of the following:
	//     WithArgument(cli.Argument().WithBinding(ptr).WithDefault(something))
	//     // or
	//     WithArgument(cli.Argument().WithBinding(ptr).WithDefault(something).Require())
	WithBindingAndDefault(pointer, def any, required bool) FlagBuilder

	MarkAsHelpFlag() FlagBuilder

	IsHelpFlag() bool

	// Require marks this Flag as being required.
	//
	// If this flag is not present in the CLI call, an error will be returned when
	// parsing the CLI input.
	Require() FlagBuilder

	IsRequired() bool

	// Build builds a new Flag instance constructed from the components set on
	// this FlagBuilder.
	Build(warnings *WarningContext) (Flag, error)
}

//

type BindingError interface {
	error

	ArgumentBuilder() ArgumentBuilder
	FlagBuilder() FlagBuilder

	Unwrap() error
}

//

// A MissingFlagError is returned on CLI parse when a flag that has been marked
// as being required was not found to be present in the CLI call.
//
// MissingFlagError is a hard error that will be returned regardless of whether
// the parser is operating in strict mode.
type MissingFlagError interface {
	error
	Flag() Flag
}
