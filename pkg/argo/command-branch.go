package argo

import (
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
	CommandChild

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

// Find Short Flag /////////////////////////////////////////////////////////////

func (c commandBranch) FindShortFlag(b byte) Flag {
	for _, group := range c.FlagGroups() {
		if flag := group.FindShortFlag(b); flag != nil {
			return flag
		}
	}

	return c.parent.FindShortFlag(b)
}

// Find Long Flag //////////////////////////////////////////////////////////////

func (c commandBranch) FindLongFlag(name string) Flag {
	for _, group := range c.FlagGroups() {
		if flag := group.FindLongFlag(name); flag != nil {
			return flag
		}
	}

	return c.parent.FindLongFlag(name)
}

func (c commandBranch) onIncomplete() {
	must(comBranchRenderer{}.RenderHelp(c, os.Stdout))
	os.Exit(1)
}
