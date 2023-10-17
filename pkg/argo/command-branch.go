package argo

import (
	"fmt"
	"os"
)

// CommandBranch represents a subcommand under a CommandTree that is an
// intermediate node between the tree root and an executable CommandLeaf.
//
// CommandBranches enable the organization of subcommands into categories.
//
// Example command tree:
//     docker
//      |- compose
//      |   |- build
//      |   |- down
//      |   |- ...
//      |- container
//      |   |- exec
//      |   |- ls
//      |   |- ...
//      |- ...
type CommandBranch interface {
	CommandNode
	CommandParent

	// Name returns the name of this CommandBranch.
	Name() string

	// Aliases returns the list of aliases assigned to this CommandBranch.
	Aliases() []string

	// HasAliases indicates whether this CommandBranch has one or more aliases
	// attached.
	HasAliases() bool

	// Matches tests whether the branch name or any of its aliases match the given
	// string.
	Matches(name string) bool

	executeCallback()

	hasCallback() bool
}

type CommandBranchCallback = func(branch CommandBranch)

type commandBranch struct {
	name          string
	desc          string
	parent        CommandNode
	aliases       []string
	flagGroups    []FlagGroup
	commandGroups []CommandGroup
	callback      CommandBranchCallback
}

// Find Child //////////////////////////////////////////////////////////////////

func (c commandBranch) FindChild(name string) CommandNode {
	for _, group := range c.commandGroups {
		if com := group.FindChild(name); com != nil {
			return com
		}
	}

	return nil
}

// On Hit //////////////////////////////////////////////////////////////////////

func (c commandBranch) executeCallback() {
	if c.callback != nil {
		c.callback(c)
	}
}

func (c commandBranch) hasCallback() bool {
	return c.callback != nil
}

// Parent //////////////////////////////////////////////////////////////////////

func (c commandBranch) Parent() CommandNode {
	return c.parent
}

func (c commandBranch) HasParent() bool {
	return c.parent != nil
}

// Description /////////////////////////////////////////////////////////////////

func (c commandBranch) Description() string {
	return c.desc
}

func (c commandBranch) HasDescription() bool {
	return len(c.desc) > 0
}

// Flag Groups /////////////////////////////////////////////////////////////////

func (c commandBranch) FlagGroups() []FlagGroup {
	return c.flagGroups
}

func (c commandBranch) HasFlagGroups() bool {
	return len(c.flagGroups) > 0
}

// Name ////////////////////////////////////////////////////////////////////////

func (c commandBranch) Name() string {
	return c.name
}

// Aliases /////////////////////////////////////////////////////////////////////

func (c commandBranch) Aliases() []string {
	return c.aliases
}

func (c commandBranch) HasAliases() bool {
	return len(c.aliases) > 0
}

// Matches /////////////////////////////////////////////////////////////////////

func (c commandBranch) Matches(name string) bool {
	if c.name == name {
		return true
	}

	for _, alias := range c.aliases {
		if alias == name {
			return true
		}
	}

	return false
}

// Command Groups //////////////////////////////////////////////////////////////

func (c commandBranch) CommandGroups() []CommandGroup {
	return c.commandGroups
}

func (c commandBranch) HasCustomCommandGroups() bool {
	return len(c.commandGroups) > 1 || c.commandGroups[0].Name() != defaultGroupName
}

// Find Short Flag /////////////////////////////////////////////////////////////

func (c commandBranch) FindShortFlag(b byte) Flag {
	var current CommandNode = c

	for current != nil {
		for _, group := range current.FlagGroups() {
			if flag := group.FindShortFlag(b); flag != nil {
				return flag
			}
		}

		current = current.Parent()
	}

	return nil
}

// Find Long Flag //////////////////////////////////////////////////////////////

func (c commandBranch) FindLongFlag(name string) Flag {
	var current CommandNode = c

	for current != nil {
		for _, group := range current.FlagGroups() {
			if flag := group.FindLongFlag(name); flag != nil {
				return flag
			}
		}

		current = current.Parent()
	}

	return nil
}

func (c commandBranch) onIncomplete() {
	fmt.Println(renderCommandBranch(c))
	os.Exit(1)
}
