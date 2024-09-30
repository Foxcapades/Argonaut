package argo

type SubCommandSpecBuilder[T any] interface {
	// WithName sets the name for the command being defined.
	WithName(name string) T

	// WithSummary sets a short summary for the command being defined.
	//
	// If a summary is not provided, but a description is provided a portion of
	// the description will be used as the summary.
	//
	// Summary text is used when rendering lists of subcommands, where longform
	// descriptions may be undesirable for readability purposes.
	WithSummary(summary string) T

	// WithDescription sets a longform description of the command being defined.
	//
	// If a description is not provided, but a summary is provided, the summary
	// will also be used for the description.
	//
	// Description text is used when rendering help text specifically for the
	// sub-command that it describes.
	WithDescription(description string) T

	// WithFlag appends a new flag to the default FlagGroup attached to the
	// command being defined.
	WithFlag(flag FlagSpecBuilder) T

	// WithFlagGroup appends a new, custom FlagGroup to the command being defined.
	WithFlagGroup(group FlagGroupSpecBuilder) T

	// WithCallback sets a callback for the command being defined.
	//
	// If set, this callback will be executed on parsing success.  Each level of
	// a CommandTree may have a callback.  The callbacks are called in the order
	// the command segments appear in the CLI call.
	WithCallback(callback CommandCallback) T
}
