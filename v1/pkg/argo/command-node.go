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

	// FindShortFlag searches this command node, then its parent command nodes
	// until a matching flag is found.
	FindShortFlag(c byte) Flag

	// FindLongFlag searches this command node, then its parent command nodes
	// until a matching flag is found.
	FindLongFlag(name string) Flag

	// TryFlag tries to fill a flag matching the given flag reference.
	//
	// To fill the target flag, this method checks the current command node's flag
	// groups for a matching flag.  If the target flag is found attached to this
	// command node it will be filled.  If it was not found on this command node,
	// this method will move up the command tree an attempt to find and fill a
	// matching flag on the parent node, moving up the parents as needed until
	// reaching the CommandTree instance.  If no matching flag could be found,
	// this method returns false.
	//
	// If a matching flag _is_ found, this method will mark that flag as being hit
	// and, if an argument value is provided, will attempt to fill that argument
	// with the provided value.  If the argument could not be parsed from the
	// given value, an error will be returned.
	TryFlag(ref FlagRef) (bool, error)
}
