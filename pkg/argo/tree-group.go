package argo

// A CommandGroup is an organizational category containing one or more commands.
//
// Subcommands that are appended to a parent directly will be placed in a
// default CommandGroup that will not have a display name unless custom command
// groups are added to that parent.
type CommandGroup interface {

	// Name returns the custom name for the CommandGroup.
	Name() string

	// HasDescription indicates whether a description value was set on this
	// CommandGroup.
	HasDescription() bool

	// Description returns the description value set on this CommandGroup.
	//
	// If no description was set on this CommandGroup, this method will return an
	// empty string.
	Description() string

	// Branches returns the BranchCommand nodes attached to this CommandGroup.
	Branches() []BranchCommand

	// HasBranches indicates whether this CommandGroup contains any branch nodes.
	HasBranches() bool

	// Leaves returns the LeafCommand nodes attached to this CommandGroup.
	Leaves() []LeafCommand

	// HasLeaves indicates whether this CommandGroup contains any leaf nodes.
	HasLeaves() bool

	// FindChild searches this CommandGroup instance for a BranchCommand or
	// LeafCommand node that matches the given string.
	//
	// Commands may match on either their name or one of their aliases.
	FindChild(name string) ChildNode
}

// A CommandGroupBuilder is used to construct a CommandGroup instance.
type CommandGroupBuilder interface {
	Name() string

	// WithDescription sets a description value for this CommandGroupBuilder.
	//
	// Descriptions are used when rendering help text.
	WithDescription(desc string) CommandGroupBuilder

	HasDescription() bool

	Description() string

	// WithBranch appends the given branch builder to this command group builder.
	WithBranch(branch BranchCommandBuilder) CommandGroupBuilder

	WithBranches(branches ...BranchCommandBuilder) CommandGroupBuilder

	HasBranches() bool

	Branches() []BranchCommandBuilder

	// WithLeaf appends the given leaf builder to this command group builder.
	WithLeaf(leaf LeafCommandBuilder) CommandGroupBuilder

	WithLeaves(leaves ...LeafCommandBuilder) CommandGroupBuilder

	HasLeaves() bool

	Leaves() []LeafCommandBuilder

	HasSubcommands() bool

	// Build attempts to build a new CommandGroup instance from the set
	// configuration.
	Build(warnings *WarningContext) (CommandGroup, error)
}
