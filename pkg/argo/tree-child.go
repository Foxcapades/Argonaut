package argo

// A ChildNode is a Node that is the child of another Node.
type ChildNode[T any] interface {
	Node[T]

	// Parent returns the parent Node for the current Node.
	//
	// If the current Node does not have a parent (meaning it is the
	// CommandTree instance) this method will return nil.
	Parent() ParentNode[any]

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

type ChildBuilder[T any] interface {
	NodeBuilder[T]

	// Name returns the name attached to the target subcommand.
	Name() string

	// WithAlias assigns the given alias to the target subcommand.
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
	//   argo.Leaf("list").WithAlias("ls")
	WithAlias(alias string) T

	// WithAliases assigns the given aliases to the target subcommand.
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
	WithAliases(aliases ...string) T

	// HasAliases indicates whether the target node has any aliases assigned.
	HasAliases() bool

	// Aliases returns a slice of all the aliases attached to this subcommand.
	//
	// If no aliases have been attached to this subcommand, the return value will
	// be nil.
	Aliases() []string

	// ParentNode returns the tree command root or branch command parent of the
	// target subcommand.
	ParentNode() ParentBuilder[any]

	SetParentNode(parent ParentBuilder[any]) T
}
