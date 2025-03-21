package argo

// A ParentNode is a Node instance that may contain child Node
// instances.
type ParentNode[T any] interface {
	Node[T]

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

	WithBranch(branch BranchCommandBuilder) T
	WithBranches(branches ...BranchCommandBuilder) T

	WithLeaf(leaf LeafCommandBuilder) T
	WithLeaves(leaves ...LeafCommandBuilder) T

	// WithCommandGroup appends the given command group builder to be built with
	// this command tree.
	//
	// Command groups are used for organizing subcommands into named groups that
	// are primarily used for rendering help text.
	WithCommandGroup(group CommandGroupBuilder) T

	WithCommandGroups(groups ...CommandGroupBuilder) T

	HasCommandGroups() bool

	CommandGroups() []CommandGroupBuilder

	HasSubcommands() bool

	WithIncompleteHandler(handler IncompleteCommandHandler[T]) T
	HasIncompleteHandler() bool
	IncompleteHandler() IncompleteCommandHandler[T]
}
