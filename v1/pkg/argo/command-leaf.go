package argo

type CommandLeaf interface {
	CommandNode
	Command

	HasOnHit() bool

	CallOnHit()

	// Aliases returns the aliases for this CommandLeaf.
	Aliases() []string

	// HasAliases indicates whether this CommandLeaf has aliases assigned.
	HasAliases() bool

	// Matches tests whether this CommandLeaf name or aliases match the given
	// string value.
	Matches(name string) bool
}

type CommandLeafCallback = func(leaf CommandLeaf)
