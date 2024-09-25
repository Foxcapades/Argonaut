package argo

// TODO: Define an expected behavior for flags without explicitly defined
//       arguments.  Perhaps they should have a "magic" catch-all argument that
//       can be used as either a toggle or a container for raw values (when the
//       cli context makes it unequivocally clear that the argument belongs to
//       the flag, such as "--foo=bar")

type FlagSpecBuilder interface {
	WithLongForm(name string) FlagSpecBuilder

	WithShortForm(name byte) FlagSpecBuilder

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

	WithSummary(summary string) FlagSpecBuilder

	WithDescription(description string) FlagSpecBuilder

	Require() FlagSpecBuilder

	Build(config Config) (FlagSpec, error)
}
