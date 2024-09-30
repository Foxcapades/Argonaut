package argo

// region Command

type Command interface {
	// Name returns the name of the command.
	Name() string

	FlagGroupContainer

	ArgumentContainer

	UnmappedInputContainer

	PassthroughContainer
}

// endregion Command

// region CommandTree

type CommandTree interface {
	CommandNode
	ParentCommand
}

type CommandTreeSpec interface {
	// Name returns the name of the command.
	Name() string

	// Description returns the description text for the command.
	Description() string

	// HasDescription returns an indicator as to whether this CommandSpec has a
	// description value set.
	//
	// If this method returns `false`, the Description method will return an empty
	// string.
	HasDescription() bool

	FlagGroupSpecContainer

	ArgumentSpecContainer

	UnmappedInputContainerSpec
}

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

// endregion CommandTree

// region CommandBranch

type CommandBranch interface {
	Name() string

	FlagGroupContainer
}

type CommandBranchSpec interface {
	Name() string

	FlagGroupSpecContainer
}

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

// endregion CommandBranch

// region CommandLeaf

type CommandLeaf interface {
	CommandNode

	SubCommand

	ArgumentContainer

	UnmappedInputContainer

	PassthroughContainer
}

type CommandLeafSpec interface {
	CommandNodeSpec

	SubCommandSpec

	ArgumentSpecContainer

	UnmappedInputContainerSpec

	PassthroughContainerSpec
}

type CommandLeafSpecBuilder interface {
	SubCommandSpecBuilder[CommandLeafSpecBuilder]

	WithPositionalArgument(arg ArgumentSpecBuilder) CommandLeafSpecBuilder

	WithUnmappedInputLabel(label string) CommandLeafSpecBuilder

	WithUnmappedInputDescription(description string) CommandLeafSpecBuilder

	WithPassthroughsDisabled() CommandLeafSpecBuilder

	Build(config Config) (CommandLeafSpec, error)
}

// endregion CommandLeaf

// region CommandNode

type CommandNode interface {
	Name() string

	FlagGroupContainer
}

// endregion CommandNode

// region CommandNode

type CommandNodeSpec interface {
	Name() string

	FlagGroupSpecContainer
}

// endregion CommandNode

// region SubCommand

type SubCommand interface {
	CommandNode
	Parent() ParentCommand
}

type SubCommandSpec interface {
	CommandNodeSpec
	Parent() ParentCommandSpec
}

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

// endregion SubCommand

// region ParentCommand

type ParentCommand interface {
	CommandNode

	SubCommandContainer
}

type ParentCommandSpec interface {
	CommandNodeSpec

	SubCommandSpecContainer
}

// endregion ParentCommand

type SubCommandContainer interface {
	SubCommands() []SubCommand

	FindSubCommand(name string) SubCommand

	Branches() []CommandBranch

	FindBranch(name string) CommandBranch

	HasBranches() bool

	Leaves() []CommandLeaf

	FindLeaf(name string) CommandLeaf

	HasLeaves() bool
}

type SubCommandSpecContainer interface {
	SubCommands() []SubCommandSpec

	FindSubCommand(name string) SubCommandSpec

	Branches() []CommandBranchSpec

	FindBranch(name string) CommandBranchSpec

	HasBranches() bool

	Leaves() []CommandLeafSpec

	FindLeaf(name string) CommandLeafSpec

	HasLeaves() bool
}

type CommandGroup interface {
	Name() string

	SubCommandContainer
}

type CommandGroupSpec interface {
	Name() string

	Description() string

	HasDescription() bool

	SubCommandSpecContainer
}

type CommandGroupSpecBuilder interface {
	WithName(name string) CommandGroupSpecBuilder

	WithDescription(description string) CommandGroupSpecBuilder

	WithBranch(branch CommandBranchSpecBuilder) CommandGroupSpecBuilder

	WithLeaf(leaf CommandLeafSpecBuilder) CommandGroupSpecBuilder

	Build(config Config) (CommandBranchSpec, error)
}

//

type CommandCallback = func(command Command)

func SimpleCommandCallback(callback func()) CommandCallback {
	return func(Command) { callback() }
}

//

type UnmappedInputContainer interface {
	UnmappedInputs() []string

	HasUnmappedInputs() bool
}

//

type UnmappedInputContainerSpec interface {
	UnmappedInputContainer

	UnmappedLabel() string

	HasUnmappedLabel() bool
}

//

type PassthroughContainer interface {
	PassthroughInputs() []string

	HasPassthroughInputs() bool
}

//

type PassthroughContainerSpec interface {
	PassthroughContainer

	AppendPassthrough(passthrough string)
}

//

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
