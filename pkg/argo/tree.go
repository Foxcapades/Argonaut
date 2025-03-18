package argo

// CommandTree represents the root of a tree of subcommands.
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
type CommandTree interface {
	ParentNode[CommandTree]

	// SelectedCommand returns the leaf command that was selected in the CLI call.
	SelectedCommand() CommandLeaf
}

type CommandTreeCallback = func(com CommandTree)

// IncompleteCommandHandler defines a function type that may be used as a callback
// for when a command leaf is not reached when parsing a command tree structure.
type IncompleteCommandHandler[T any] = func(command T)

// A CommandTreeBuilder is a builder type used to construct a CommandTree
// instance.
//
// A CommandTree is a command that consists of branching subcommands.  Examples
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
type CommandTreeBuilder interface {
	ParentBuilder[CommandTreeBuilder]

	// // WithDescription sets a description value for the root of this command tree.
	// //
	// // Descriptions are used when rendering help text.
	// WithDescription(desc string) CommandTreeBuilder

	// WithHelpDisabled disables the automatic `-h` and `--help` flags for
	// rendering help text.
	WithHelpDisabled() CommandTreeBuilder

	// WithBranch appends the given branch builder to be built with this command
	// tree.
	//
	// The built branch will be available as a subcommand directly under the root
	// command call.
	// WithBranch(branch BranchBuilder) CommandTreeBuilder

	// WithLeaf appends the given leaf builder to be built with this command tree.
	//
	// The built leaf will be available as a subcommand directly under the root
	// command call.
	// WithLeaf(leaf LeafBuilder) CommandTreeBuilder

	// WithCommandGroup appends the given command group builder to be built with
	// this command tree.
	//
	// Command groups are used for organizing subcommands into named groups that
	// are primarily used for rendering help text.
	WithCommandGroup(group CommandGroupBuilder) CommandTreeBuilder

	Build(warnings *WarningContext) (CommandTree, error)

	// Parse builds the command tree and attempts to parse the given CLI arguments
	// into that command tree's components.
	Parse(args []string) (CommandTree, error)

	// MustParse calls Parse and panics if an error is returned.
	MustParse(args []string) CommandTree
}
