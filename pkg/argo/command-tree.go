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
//     docker
//      |- compose
//      |   |- build
//      |   |- down
//      |   |- ...
//      |- container
//      |   |- exec
//      |   |- ls
//      |   |- ...
//      |- ...
type CommandTree interface {
	CommandNode
	CommandParent

	// SelectedCommand returns the leaf command that was selected in the CLI call.
	SelectedCommand() CommandLeaf

	// SelectCommand selects the given leaf command.
	SelectCommand(leaf CommandLeaf)

	IsHelpDisabled() bool
}
