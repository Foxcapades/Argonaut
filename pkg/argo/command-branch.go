package argo

// CommandBranch represents a subcommand under a CommandTree that is an
// intermediate node between the tree root and an executable CommandLeaf.
//
// CommandBranches enable the organization of subcommands into categories.
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
type CommandBranch interface {
	CommandNode
	CommandParent

	// Name returns the name of this CommandBranch.
	Name() string

	// Aliases returns the list of aliases assigned to this CommandBranch.
	Aliases() []string

	// HasAliases indicates whether this CommandBranch has one or more aliases
	// attached.
	HasAliases() bool

	// Matches tests whether the branch name or any of its aliases match the given
	// string.
	Matches(name string) bool

	RunCallback()

	HasCallback() bool
}

type CommandBranchCallback = func(branch CommandBranch)
