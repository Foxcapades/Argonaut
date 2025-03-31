package argo

// WARNING:
//   This is a generated file!  Edits here will be lost!

// TreeCommand represents the root of a tree of subcommands.
//
// The command tree consists of branch and leaf nodes.  The branch nodes can be
// thought of as categories for containing sub-branches and/or leaves.  Leaf
// nodes are the actual callable command implementations.
//
// All levels of the command tree accept flags, with sub-node flags taking
// priority over parent node flags on flag collision.  Leaf nodes, however, are
// the only nodes that accept positional arguments, or passthroughs.
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
type TreeCommand interface {
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

	HasCallback() bool

	Callback() CommandCallback[TreeCommand]

	IsHelpDisabled() bool

	// IncompleteHandler returns the incomplete command handler function that was
	// attached to this TreeCommand through the TreeCommandBuilder.
	//
	// If no incomplete command handler function was attached, this method returns
	// nil.
	IncompleteHandler() IncompleteCommandHandler[TreeCommand]

	// SelectedCommand returns the leaf command that was selected in the CLI call.
	SelectedCommand() LeafCommand

	// FindShortFlag looks up a target Flag instance by its short-form character.
	//
	// If no such flag exists on this node, nil will be returned.
	FindShortFlag(c byte) Flag

	// FindShortFlagRecursive looks up a target Flag instance by its short-form
	// character recursively on the selected LeafCommand node and any intermediate
	// nodes back to this TreeCommand.
	//
	// If no such flag exists in the tree hierarchy, nil will be returned.
	//
	// If this method is called before the CLI input has been parsed into this
	// TreeCommand
	//
	// The behavior of this method is not affected by the InheritParentFlags
	// TreeCommandOption setting.
	FindShortFlagRecursive(c byte) Flag

	// FindLongFlag looks up a target Flag instance by its long-form name.
	//
	// If no such flag exists on this node, the ancestor nodes will be checked
	// parent by parent.
	//
	// To search this node and all parent nodes, use FindLongFlagRecursive
	FindLongFlag(name string) Flag

	// FindLongFlagRecursive looks up a target Flag instance by its long-form
	// name recursively on this node and all ancestor nodes.
	//
	// If no such flag exists in the tree hierarchy, nil will be returned.
	//
	// The behavior of this method is not affected by the InheritParentFlags
	// TreeCommandOption setting.
	FindLongFlagRecursive(name string) Flag

	Options() TreeCommandOptions
}

// A TreeCommandBuilder is a builder type used to construct a TreeCommand
// instance.
//
// A TreeCommand is a command that consists of branching subcommands.  Examples
// of such commands include the `go` command, `docker`, or `kubectl`.
//
// To use the Docker command example we have a command tree that includes the
// following:
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
type TreeCommandBuilder interface {
	ParentNodeBuilder

	// WithDescription sets the description value that will be used for the built
	// cli command or subcommand.
	//
	// Descriptions are used when rendering help text.
	WithDescription(desc string) TreeCommandBuilder

	HasDescription() bool

	Description() string

	// WithHelpDisabled disables the automatic `-h` and `--help` flags for
	// rendering help text.
	WithHelpDisabled() TreeCommandBuilder

	IsHelpDisabled() bool

	// WithFlagGroup appends the given FlagGroupBuilder to this CLI component
	// builder.
	WithFlagGroup(group FlagGroupBuilder) TreeCommandBuilder

	WithFlagGroups(groups ...FlagGroupBuilder) TreeCommandBuilder

	HasFlagGroups(includeDefault bool) bool

	FlagGroups(includeDefault bool) []FlagGroupBuilder

	// WithFlag attaches the given FlagBuilder to the default FlagGroupBuilder
	// instance attached to this CLI component builder.
	WithFlag(flag FlagBuilder) TreeCommandBuilder

	WithFlags(flags ...FlagBuilder) TreeCommandBuilder

	HasFlags() bool

	WithCallback(callback CommandCallback[TreeCommand]) TreeCommandBuilder

	Callback() CommandCallback[TreeCommand]

	HasCallback() bool

	// WithBranch appends the given BranchCommandBuilder instance to this
	// TreeCommandBuilder instance.
	//
	// Branch commands added to this TreeCommandBuilder are placed in the default
	// CommandGroupBuilder.
	WithBranch(branch BranchCommandBuilder) TreeCommandBuilder

	// WithBranches appends the given BranchCommandBuilder instances to this
	// TreeCommandBuilder instance.
	//
	// Branch commands added to this TreeCommandBuilder are placed in the default
	// CommandGroupBuilder.
	WithBranches(branches ...BranchCommandBuilder) TreeCommandBuilder

	// WithLeaf appends the given LeafCommandBuilder instance to this
	// TreeCommandBuilder instance.
	//
	// Leaf commands added to this TreeCommandBuilder are placed in the default
	// CommandGroupBuilder.
	WithLeaf(leaf LeafCommandBuilder) TreeCommandBuilder

	// WithLeaves appends the given LeafCommandBuilder instances to this
	// TreeCommandBuilder instance.
	//
	// Leaf commands added to this TreeCommandBuilder are placed in the default
	// CommandGroupBuilder.
	WithLeaves(leaves ...LeafCommandBuilder) TreeCommandBuilder

	// WithCommandGroup appends the given CommandGroupBuilder instance to this
	// TreeCommandBuilder instance.
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
	WithCommandGroup(group CommandGroupBuilder) TreeCommandBuilder

	WithCommandGroups(groups ...CommandGroupBuilder) TreeCommandBuilder

	WithIncompleteHandler(handler IncompleteCommandHandler[TreeCommand]) TreeCommandBuilder
	IncompleteHandler() IncompleteCommandHandler[TreeCommand]

	WithOptions(opts TreeCommandOptions) TreeCommandBuilder

	Options() TreeCommandOptions
}
