package argo

// A CommandBranchBuilder instance may be used to configure a new CommandBranch
// instance to be built.
//
// CommandBranches are intermediate steps between the root of the CommandTree
// and the CommandLeaf instances.
//
// For example, given the following command example, the tree is "foo", the
// branch is "bar", and the leaf is "fizz":
//     ./foo bar fizz
type CommandBranchBuilder interface {

	// GetName returns the name assigned to this CommandBranchBuilder.
	//
	// CommandBranch names are required and thus are set at the time that the
	// CommandBranchBuilder instance is constructed.
	GetName() string

	// WithAliases appends the given alias strings to this CommandBranchBuilder's
	// alias list.
	//
	// Aliases must be unique and must not conflict with any other command branch
	// or leaf names or aliases at a given command tree level.
	//
	// Example:
	//     cli.CommandBranch("service").
	//         WithAliases("svc")
	WithAliases(aliases ...string) CommandBranchBuilder

	// GetAliases returns the aliases assigned to this CommandBranchBuilder.
	GetAliases() []string

	// Parent sets the parent CommandNode for the CommandBranch being built.
	//
	// Values set using this method before build time will be disregarded.
	Parent(node CommandNode)

	// WithDescription sets a description value for the CommandBranch being built.
	//
	// Descriptions are used when rendering help text.
	WithDescription(desc string) CommandBranchBuilder

	HasDescription() bool

	GetDescription() string

	// WithHelpDisabled disables the automatic '-h | --help' flag that is enabled
	// by default.
	WithHelpDisabled() CommandBranchBuilder

	// WithBranch appends a child branch to the default CommandGroup for this
	// CommandBranchBuilder.
	WithBranch(branch CommandBranchBuilder) CommandBranchBuilder

	// WithLeaf appends a child leaf to the default CommandGroup for this
	// CommandBranchBuilder.
	WithLeaf(leaf CommandLeafBuilder) CommandBranchBuilder

	// WithCommandGroup appends a custom CommandGroup to this
	// CommandBranchBuilder.
	//
	// CommandGroups are used to organize subcommands into named categories that
	// are primarily used for rendering help text.
	WithCommandGroup(group CommandGroupBuilder) CommandBranchBuilder

	// WithFlag appends the given FlagBuilder to the default FlagGroup attached to
	// this CommandBranchBuilder.
	WithFlag(flag FlagBuilder) CommandBranchBuilder

	// WithFlagGroup appends the given custom FlagGroupBuilder to this
	// CommandBranchBuilder instance.
	//
	// Custom flag groups are primarily used for categorizing flags in the
	// rendered help text.
	WithFlagGroup(flagGroup FlagGroupBuilder) CommandBranchBuilder

	WithCallback(cb CommandBranchCallback) CommandBranchBuilder

	HasCallback() bool

	GetCallback() CommandBranchCallback

	// Build attempts to construct a CommandBranch configuration set on this
	// CommandBranchBuilder instance.
	//
	// At minimum, for the building of a CommandBranch to be successful, the
	// branch must have a non-blank name value and at least one child node which
	// may be one or more leaf nodes or branch nodes.
	Build() (CommandBranch, error)
}
