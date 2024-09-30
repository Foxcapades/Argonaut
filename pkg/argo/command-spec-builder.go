package argo

type CommandSpecBuilder interface {
	// WithDescription sets a description value for the command being defined.
	WithDescription(description string) CommandSpecBuilder

	// WithFlag adds the given flag to the default FlagGroup for the command being
	// defined.
	WithFlag(flag FlagSpecBuilder) CommandSpecBuilder

	// WithFlagGroup adds the given flag group definition to the command being
	// defined.
	WithFlagGroup(group FlagGroupSpecBuilder) CommandSpecBuilder

	// WithPositionalArgument adds the given argument definition as a positional
	// argument for the command being defined.
	WithPositionalArgument(arg ArgumentSpecBuilder) CommandSpecBuilder

	// WithCallback adds a callback to the command being defined.
	//
	// Callbacks will be executed in the order they are attached to the command.
	WithCallback(callback CommandCallback) CommandSpecBuilder

	// Build attempts to build a CommandSpec from the current state of this
	// CommandSpecBuilder.
	Build(config Config) (CommandSpec, error)
}
