package argo

// TODO: Define an expected behavior for flags without explicitly defined
//       arguments.  Perhaps they should have a "magic" catch-all argument that
//       can be used as either a toggle or a container for raw values (when the
//       cli context makes it unequivocally clear that the argument belongs to
//       the flag, such as "--foo=bar")

type FlagSpecBuilder interface {

	// WithLongForm sets the long-form flag name that the flag may be referenced
	// by on the CLI.
	//
	// Long-form flags consist of one or more characters preceded immediately by
	// two prefix/leader characters (typically dashes).
	//
	// Long-form flags must start with an alphanumeric character and may only
	// consist of alphanumeric characters, dashes, and/or underscores.
	WithLongForm(name string) FlagSpecBuilder

	// WithShortForm sets the short-form flag character that the flag may be
	// referenced by on the CLI.
	//
	// Short-form flags consist of a single character preceded by either a
	// prefix/leader character, or by one or more other short-form flags.
	WithShortForm(name byte) FlagSpecBuilder

	// WithDescription sets a description for the flag that will be used when
	// printing command help text.
	WithDescription(description string) FlagSpecBuilder

	// WithArgument sets a custom argument for the flag that will handle arguments
	// passed for the flag in CLI calls.
	WithArgument(arg ArgumentSpecBuilder) FlagSpecBuilder

	// WithLazyCallback sets a callback on the FlagSpec being built that will be
	// executed on successful parsing of a valid CLI call that included this flag.
	//
	// Callbacks are executed in the order they are configured.
	WithLazyCallback(callback FlagCallback) FlagSpecBuilder

	// WithImmediateCallback sets a callback on the FlagSpec being built that will
	// be executed as soon as the flag is encountered, every time it is
	// encountered while parsing the CLI call, regardless of whether the CLI call
	// was valid.
	//
	// Immediate callbacks will be given the flag state as it was when the flag
	// was encountered, which may not include the final usage count if more flag
	// usages appear later in the CLI call.
	//
	// Callbacks are executed in the order they are configured.
	//
	// If an unrecoverable error is encountered while parsing a CLI call, Argonaut
	// will halt parsing and no further callbacks will be executed regardless of
	// whether there are more flag instances in the CLI call text.
	WithImmediateCallback(callback FlagCallback) FlagSpecBuilder

	// Require marks the flag being built as being required.
	//
	// If the flag is marked as required and is not present in a CLI call, an
	// error will be generated when parsing the CLI input.
	Require() FlagSpecBuilder

	// Build builds a new FlagSpec instance from the configuration set on this
	// FlagSpecBuilder instance.
	Build(config Config) (FlagSpec, error)
}
