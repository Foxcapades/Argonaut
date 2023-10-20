package argo

type CommandChild interface {
	CommandNode

	HasAliases() bool

	Aliases() []string

	// Matches tests whether the branch name or any of its aliases match the given
	// string.
	Matches(name string) bool
}
