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

	// FindShortFlag looks up a target Flag instance by its short-form character.
	//
	// If no such flag exists on this node, nil will be returned.
	//
	// To search this node and all parent nodes, use FindShortFlagRecursive
	FindShortFlag(c byte) Flag

	// FindShortFlagRecursive looks up a target Flag instance by its short-form
	// character recursively on this node and all ancestor nodes.
	//
	// If no such flag exists in the tree hierarchy, nil will be returned.
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
