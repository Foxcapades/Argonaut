package argo

// A ParentNode is a Node instance that may contain child Node
// instances.
type ParentNode[T any] interface {
	Node[T]

	HasSubcommands() bool

	// CommandGroups returns the CommandGroup instances attached to this
	// ParentNode node.
	CommandGroups(includeDefault bool) []CommandGroup

	HasCommandGroups(includeDefault bool) bool

	// FindChild searches this ParentNode's CommandGroup instances for a
	// subcommand that matches the given string.
	//
	// A subcommand may match on either its name or one of its aliases.
	FindChild(name string) ChildNode[any]

	HasSelectedChild() bool

	SelectChild(name string) bool

	SelectedChild() ChildNode[any]

	HasIncompleteHandler() bool

	IncompleteHandler() IncompleteCommandHandler[T]
}

type ParentBuilder[T any] interface {
	NodeBuilder[T]

	WithBranch(branch BranchBuilder) T
	WithBranches(branches ...BranchBuilder) T

	WithLeaf(leaf LeafBuilder) T
	WithLeaves(leaves LeafBuilder) T

	HasSubcommands() bool
	Subcommands() []ChildBuilder[any]

	WithIncompleteHandler(handler IncompleteCommandHandler[T]) T
	HasIncompleteHandler() bool
	IncompleteHandler() IncompleteCommandHandler[T]
}
