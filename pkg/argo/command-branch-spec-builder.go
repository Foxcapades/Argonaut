package argo

type CommandBranchSpecBuilder interface {
	SubCommandSpecBuilder[CommandBranchSpecBuilder]

	// WithCommandGroup appends a new, custom CommandGroup to the CommandBranch
	// being defined.
	WithCommandGroup(group CommandGroupSpecBuilder) CommandBranchSpecBuilder

	// WithBranch appends the given CommandBranch specification to the default
	// CommandGroup attached to the CommandBranch being built.
	WithBranch(branch CommandBranchSpecBuilder) CommandBranchSpecBuilder

	// WithLeaf appends the given CommandLeaf specification to the default
	// CommandGroup attached to the CommandBranch being built.
	WithLeaf(leaf CommandLeafSpecBuilder) CommandBranchSpecBuilder

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
	WithIncompleteHandler(handler IncompleteCommandHandler) CommandBranchSpecBuilder

	Build(config Config) (CommandBranchSpec, error)
}
