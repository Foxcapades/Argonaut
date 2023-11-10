package argo

// A CommandChild is a CommandNode that is the child of another CommandNode.
type CommandChild interface {
	CommandNode

	// HasAliases indicates whether this CommandChild has alias strings that may
	// be used to reference this CommandChild instead of the CommandChild's
	// assigned name.
	HasAliases() bool

	// Aliases returns the aliases attached to this CommandChild.
	Aliases() []string

	// Matches tests whether the branch name or any of its aliases match the given
	// string.
	Matches(name string) bool
}
