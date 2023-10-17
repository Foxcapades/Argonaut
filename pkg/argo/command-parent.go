package argo

// A CommandParent is a CommandNode instance that may contain child CommandNode
// instances.
type CommandParent interface {
	// CommandGroups returns the CommandGroup instances attached to this
	// CommandParent node.
	CommandGroups() []CommandGroup

	// FindChild searches this CommandParent's CommandGroup instances for a
	// subcommand that matches the given string.
	//
	// A subcommand may match on either its name or one of its aliases.
	FindChild(name string) CommandNode

	onIncomplete()
}
