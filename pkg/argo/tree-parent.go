package argo

// WARNING:
//   This is a generated file!  Edits here will be lost!

type ParentNode interface {
	// CommandGroups returns the CommandGroup instances attached to this
	// ParentNode node.
	CommandGroups() []CommandGroup

	HasCommandGroups() bool

	HasSubcommands() bool

	// FindChild searches this ParentNode's CommandGroup instances for a
	// subcommand that matches the given string.
	//
	// A subcommand may match on either its name or one of its aliases.
	FindChild(name string) ChildNode

	HasSelectedChild() bool

	SelectChild(name string) bool

	SelectedChild() ChildNode

	HasIncompleteHandler() bool
}

type ParentNodeBuilder interface {
	HasCommandGroups(includeDefault bool) bool

	CommandGroups(includeDefault bool) []CommandGroupBuilder

	HasSubcommands() bool

	HasIncompleteHandler() bool
}
