package argo

// WARNING:
//   This is a generated file!  Edits here will be lost!

type ChildNode interface {
	// Name returns the name of the subcommand.
	Name() string

	// Parent returns the parent Node for the current Node.
	//
	// If the current Node does not have a parent (meaning it is the
	// CommandTree instance) this method will return nil.
	Parent() ParentNode

	// HasAliases indicates whether this ChildNode has alias strings that may
	// be used to reference this ChildNode instead of the ChildNode's
	// assigned name.
	HasAliases() bool

	// Aliases returns the aliases attached to this ChildNode.
	Aliases() []string

	// Matches tests whether the branch name or any of its aliases match the given
	// string.
	Matches(name string) bool
}

type ChildNodeBuilder interface {
	// Name returns the name attached to the target command node.
	Name() string

	// HasAliases indicates whether the command node has any aliases assigned.
	HasAliases() bool

	// Aliases returns a slice of all the aliases attached to this subcommand.
	//
	// If no aliases have been attached to this subcommand, the return value will
	// be nil.
	Aliases() []string
}
