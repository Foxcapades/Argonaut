package argo

// WARNING:
//   This is a generated file!  Edits here will be lost!

// BranchCommand represents a subcommand under a CommandTree that is an
// intermediate node between the tree root and an executable LeafCommand.
//
// CommandBranches enable the organization of subcommands into categories.
//
// Example command tree:
//
//	docker
//	 |- compose
//	 |   |- build
//	 |   |- down
//	 |   |- ...
//	 |- container
//	 |   |- exec
//	 |   |- ls
//	 |   |- ...
//	 |- ...
type BranchCommand interface {
	ChildNode
	ParentNode

	// Name returns the name of the command or subcommand.
	//
	// For root commands, this value will be the name of the command as it
	// appeared in the command line call.
	Name() string

	// Description returns the description value assigned to this node.
	//
	// Description values are used when rendering help text.
	Description() string

	// HasDescription indicates whether this Node has a description value
	// set.
	HasDescription() bool

	// FlagGroups returns the flag groups assigned to this Node.
	//
	// This method will only return flag groups that had flags assigned to them,
	// the rest of the flag groups will have been filtered out when the node was
	// built.
	FlagGroups() []FlagGroup

	// HasFlagGroups indicates whether this Node has at least one populated
	// flag group.
	HasFlagGroups() bool

	// FindShortFlag looks up a target Flag instance by its short-form character.
	//
	// If no such flag exists on this Node or any of its parents, this
	// method will return nil.
	FindShortFlag(c byte) Flag

	// FindLongFlag looks up a target Flag instance by its long-form name.
	//
	// If no such flag exists on this Node or any of its parents, this
	// method will return nil.
	FindLongFlag(name string) Flag

	HasCallback() bool

	Callback() CommandCallback[BranchCommand]

	IsHelpDisabled() bool

	// IncompleteHandler returns the incomplete command handler function that was
	// attached to this BranchCommand through the BranchCommandBuilder.
	//
	// If no incomplete command handler function was attached, this method returns
	// nil.
	IncompleteHandler() IncompleteCommandHandler[BranchCommand]
}

// A BranchCommandBuilder instance may be used to configure a new BranchCommand
// instance to be built.
//
// CommandBranches are intermediate steps between the root of the CommandTree
// and the LeafCommand instances.
//
// For example, given the following command example, the tree is "foo", the
// branch is "bar", and the leaf is "fizz":
//
//	./foo bar fizz
type BranchCommandBuilder interface {
	ChildNodeBuilder
	ParentNodeBuilder

	// WithDescription sets the description value that will be used for the built
	// cli command or subcommand.
	//
	// Descriptions are used when rendering help text.
	WithDescription(desc string) BranchCommandBuilder

	HasDescription() bool

	Description() string

	// WithHelpDisabled disables the automatic `-h` and `--help` flags for
	// rendering help text.
	WithHelpDisabled() BranchCommandBuilder

	IsHelpDisabled() bool

	// WithFlagGroup appends the given FlagGroupBuilder to this CLI component
	// builder.
	WithFlagGroup(group FlagGroupBuilder) BranchCommandBuilder

	WithFlagGroups(groups ...FlagGroupBuilder) BranchCommandBuilder

	HasFlagGroups(includeDefault bool) bool

	FlagGroups(includeDefault bool) []FlagGroupBuilder

	// WithFlag attaches the given FlagBuilder to the default FlagGroupBuilder
	// instance attached to this CLI component builder.
	WithFlag(flag FlagBuilder) BranchCommandBuilder

	WithFlags(flags ...FlagBuilder) BranchCommandBuilder

	HasFlags() bool

	WithCallback(callback CommandCallback[BranchCommand]) BranchCommandBuilder

	Callback() CommandCallback[BranchCommand]

	HasCallback() bool

	// WithBranch appends the given BranchCommandBuilder instance to this
	// BranchCommandBuilder instance.
	//
	// Branch commands added to this BranchCommandBuilder are placed in the default
	// CommandGroupBuilder.
	WithBranch(branch BranchCommandBuilder) BranchCommandBuilder

	// WithBranches appends the given BranchCommandBuilder instances to this
	// BranchCommandBuilder instance.
	//
	// Branch commands added to this BranchCommandBuilder are placed in the default
	// CommandGroupBuilder.
	WithBranches(branches ...BranchCommandBuilder) BranchCommandBuilder

	// WithLeaf appends the given LeafCommandBuilder instance to this
	// BranchCommandBuilder instance.
	//
	// Leaf commands added to this BranchCommandBuilder are placed in the default
	// CommandGroupBuilder.
	WithLeaf(leaf LeafCommandBuilder) BranchCommandBuilder

	// WithLeaves appends the given LeafCommandBuilder instances to this
	// BranchCommandBuilder instance.
	//
	// Leaf commands added to this BranchCommandBuilder are placed in the default
	// CommandGroupBuilder.
	WithLeaves(leaves ...LeafCommandBuilder) BranchCommandBuilder

	// WithCommandGroup appends the given CommandGroupBuilder instance to this
	// BranchCommandBuilder instance.
	//
	// Command groups are used for organizing subcommands into named groups that
	// are primarily used for rendering help text.
	//
	// Example usage:
	//   cli.Tree().
	//       WithCommandGroup(cli.CommandGroup("My Command Group").
	//       WithBranch(cli.Branch("foo").
	//           WithLeaf(cli.Leaf("bar"))))
	//
	// Resulting help text:
	//
	WithCommandGroup(group CommandGroupBuilder) BranchCommandBuilder

	WithCommandGroups(groups ...CommandGroupBuilder) BranchCommandBuilder

	WithIncompleteHandler(handler IncompleteCommandHandler[BranchCommand]) BranchCommandBuilder
	IncompleteHandler() IncompleteCommandHandler[BranchCommand]

	// WithAlias assigns the given alias to the target command node.
	//
	// Command aliases must be unique per level in a command tree.  This means
	// that for any given step in the tree, no alias may conflict with another
	// branch or leaf subcommand's name or aliases.
	//
	// This also applies if a subcommand node is reused at multiple levels of the
	// command tree.
	//
	// If a conflict is found between subcommand names and/or aliases, an error
	// will be returned when attempting to build the command tree.
	//
	// Example:
	//   cli.Branch("list").WithAlias("ls")
	WithAlias(alias string) BranchCommandBuilder

	// WithAliases assigns the given aliases to the target command node.
	//
	// Command aliases must be unique per level in a command tree.  This means
	// that for any given step in the tree, no alias may conflict with another
	// branch or leaf subcommand's name or aliases.
	//
	// This also applies if a subcommand node is reused at multiple levels of the
	// command tree.
	//
	// If a conflict is found between subcommand names and/or aliases, an error
	// will be returned when attempting to build the command tree.
	WithAliases(aliases ...string) BranchCommandBuilder
}
