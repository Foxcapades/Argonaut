package argo

type CommandTreeSpecBuilder interface {
	// WithDescription sets a description value for the root of the command tree.
	WithDescription(description string) CommandTreeSpecBuilder

	// WithFlag appends a new flag to the default FlagGroup attached to the root
	// of the CommandTree.
	WithFlag(flag FlagSpecBuilder) CommandTreeSpecBuilder

	// WithFlagGroup appends a new, custom FlagGroup to the root of the
	// CommandTree.
	WithFlagGroup(group FlagGroupSpecBuilder) CommandTreeSpecBuilder

	// WithCallback sets a callback for the CommandTree.
	//
	// If set, this callback will be executed on parsing success.  Each level of
	// a CommandTree may have a callback.  The callbacks are called in the order
	// the command segments appear in the CLI call.
	WithCallback(callback CommandCallback) CommandTreeSpecBuilder

	// WithCommandGroup appends a new, custom CommandGroup to the root of the
	// CommandTree.
	WithCommandGroup(group CommandGroupSpecBuilder) CommandTreeSpecBuilder

	// WithBranch appends the given CommandBranch specification to the default
	// CommandGroup attached to the root of the CommandTree.
	WithBranch(branch CommandBranchSpecBuilder) CommandTreeSpecBuilder

	// WithLeaf appends the given CommandLeaf specification to the default
	// CommandGroup attached to the root of the CommandTree.
	WithLeaf(leaf CommandLeafSpecBuilder) CommandTreeSpecBuilder

	// WithIncompleteHandler sets a custom handler for cases when a CLI call is
	// incomplete.
	//
	// The given IncompleteCommandHandler is called when a CommandTree is called,
	// but a CommandLeaf is not reached in the CLI call.
	//
	// If no handler is provided, the default behavior is to print the help text
	// for the furthest command node reached, then exit with code `1`.
	//
	// Child nodes will inherit this handler if they do not have their own handler
	// set.
	WithIncompleteHandler(handler IncompleteCommandHandler) CommandTreeSpecBuilder

	Build(config Config) (CommandTreeSpec, error)
}
