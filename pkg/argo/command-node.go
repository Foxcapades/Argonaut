package argo

// A CommandNode is a base type common to all elements in a CommandTree.
type CommandNode interface {

	// Parent returns the parent CommandNode for the current CommandNode.
	//
	// If the current CommandNode does not have a parent (meaning it is the
	// CommandTree instance) this method will return nil.
	Parent() CommandNode

	// HasParent indicates whether this CommandNode has a parent node.
	//
	// This means that the current CommandNode instance is either a branch or a
	// leaf node.
	HasParent() bool

	// Name returns the name of the command or subcommand.
	//
	// For the CommandTree node, this method will return the name of the CLI
	// command that was called.
	//
	// For branch and leaf nodes, it will return the assigned name of that node.
	Name() string

	// Description returns the description value assigned to this node.
	//
	// Description values are used when rendering help text.
	Description() string

	// HasDescription indicates whether this CommandNode has a description value
	// set.
	HasDescription() bool

	// FlagGroups returns the flag groups assigned to this CommandNode.
	//
	// This method will only return flag groups that had flags assigned to them,
	// the rest of the flag groups will have been filtered out when the node was
	// built.
	FlagGroups() []FlagGroup

	// HasFlagGroups indicates whether this CommandNode has at least one populated
	// flag group.
	HasFlagGroups() bool

	// FindShortFlag looks up a target Flag instance by its short-form character.
	//
	// If no such flag exists on this CommandNode or any of its parents, this
	// method will return nil.
	FindShortFlag(c byte) Flag

	// FindLongFlag looks up a target Flag instance by its long-form name.
	//
	// If no such flag exists on this CommandNode or any of its parents, this
	// method will return nil.
	FindLongFlag(name string) Flag

	Warnings() []string

	AppendWarning(warning string)
}
