package argo

type IncompleteCommandHandler interface {
	// OnIncomplete handles the case when a CommandTree is called, but no
	// CommandLeaf is reached in the CLI call.
	//
	// If this handler returns, Argonaut will end the process immediately with an
	// exit code of `1`.
	//
	// The given command parameter will be either a CommandTree or CommandBranch
	// instance depending on whether a subcommand was reached in the CLI call.
	OnIncomplete(command ParentCommand)
}

type IncompleteCommandHandlerFn func(command ParentCommand)

func (i IncompleteCommandHandlerFn) OnIncomplete(command ParentCommand) {
	i(command)
}
