package tree

import (
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

type commandBranch struct {
	name          string
	desc          string
	parent        argo.ParentNode[any]
	aliases       []string
	flagGroups    []argo.FlagGroup
	commandGroups []argo.CommandGroup
	callback      argo.CommandBranchCallback
	warnings      *argo.WarningContext

	onIncompleteHandler argo.IncompleteCommandHandler[argo.Branch]
}

func (c commandBranch) FindChild(name string) argo.ChildNode[any] {
	for _, group := range c.commandGroups {
		if com := group.FindChild(name); com != nil {
			return com
		}
	}

	return nil
}

func (c commandBranch) HasCallback() bool {
	return c.callback != nil
}

// Parent //////////////////////////////////////////////////////////////////////

func (c commandBranch) Parent() argo.ParentNode[any] {
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

func (c commandBranch) FlagGroups() []argo.FlagGroup {
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

func (c commandBranch) CommandGroups() []argo.CommandGroup {
	return c.commandGroups
}

// Find Short Flag /////////////////////////////////////////////////////////////

func (c commandBranch) FindShortFlag(b byte) argo.Flag {
	for _, group := range c.FlagGroups() {
		if flag := group.FindShortFlag(b); flag != nil {
			return flag
		}
	}

	return c.parent.FindShortFlag(b)
}

func (c commandBranch) FindLongFlag(name string) argo.Flag {
	for _, group := range c.FlagGroups() {
		if flag := group.FindLongFlag(name); flag != nil {
			return flag
		}
	}

	return c.parent.FindLongFlag(name)
}

func (c commandBranch) Warnings() []string {
	return c.warnings.GetWarnings()
}

func (c commandBranch) AppendWarning(warning string) {
	c.warnings.AppendWarning(warning)
}

func (c commandBranch) onIncomplete(node argo.ParentNode) {
	if c.onIncompleteHandler != nil {
		c.onIncompleteHandler(node)
	} else {
		c.parent.IncompleteHandler(node)
	}
}
