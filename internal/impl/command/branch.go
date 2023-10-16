package command

import (
	"github.com/Foxcapades/Argonaut/internal/consts"
	"github.com/Foxcapades/Argonaut/pkg/argo"
)

// implements argo.CommandBranch
type commandBranch struct {
	name          string
	desc          string
	parent        argo.CommandNode
	aliases       []string
	flagGroups    []argo.FlagGroup
	commandGroups []argo.CommandGroup
	onHit         argo.BranchHitCallback
}

// Find Child //////////////////////////////////////////////////////////////////

func (c commandBranch) FindChild(name string) argo.CommandNode {
	for _, group := range c.commandGroups {
		if com := group.FindChild(name); com != nil {
			return com
		}
	}

	return nil
}

// On Hit //////////////////////////////////////////////////////////////////////

func (c commandBranch) CallOnHit() {
	if c.onHit != nil {
		c.onHit(c)
	}
}

func (c commandBranch) HasOnHitCallback() bool {
	return c.onHit != nil
}

// Parent //////////////////////////////////////////////////////////////////////

func (c commandBranch) Parent() argo.CommandNode {
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

func (c commandBranch) HasCustomCommandGroups() bool {
	return len(c.commandGroups) > 1 || c.commandGroups[0].Name() != consts.DefaultGroupName
}

// Find Short Flag /////////////////////////////////////////////////////////////

func (c commandBranch) FindShortFlag(b byte) argo.Flag {
	var current argo.CommandNode = c

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

func (c commandBranch) FindLongFlag(name string) argo.Flag {
	var current argo.CommandNode = c

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

// Try Flag ////////////////////////////////////////////////////////////////////

func (c commandBranch) TryFlag(ref argo.FlagRef) (bool, error) {
	var current argo.CommandNode = c

	for current != nil {
		for _, group := range current.FlagGroups() {
			if ok, err := group.TryFlag(ref); ok || err != nil {
				return ok, err
			}
		}

		current = current.Parent()
	}

	return false, nil
}
