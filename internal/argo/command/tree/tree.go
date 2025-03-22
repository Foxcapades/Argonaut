package tree

import (
	"os"
	"path/filepath"

	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

type Tree struct {
	disableHelp bool
  description string
  flagGroups  []argo.FlagGroup
  callback    argo.CommandCallback[argo.CommandTree]
	commandGroups []argo.CommandGroup
	selectedChild argo.ChildNode
	incompleteFn argo.IncompleteCommandHandler[argo.CommandTree]
	selectedLeaf argo.LeafCommand
}

func (_ *Tree) Name() string {
	return filepath.Base(os.Args[0])
}

func (i *Tree) SelectedCommand() argo.LeafCommand {
	return i.selectedLeaf
}

func (i *Tree) HasCommandGroups(includeDefault bool) bool {
	return len(i.groups) > 1 || (includeDefault && i.hasDefaultCommandGroup())
}

func (i *Tree) CommandGroups(includeDefault bool) []argo.CommandGroup {
	if !i.HasCommandGroups(includeDefault) {
		return nil
	}

	if includeDefault && i.hasDefaultCommandGroup() {
		return i.groups
	}

	return i.groups[1:]
}

func (i *Tree) hasDefaultCommandGroup() bool {
	return i.groups[0] != nil
}

func (i *Tree) FindChild(name string) argo.ChildNode {
	for _, group := range i.groups {
		if child := group.FindChild(name); child != nil {
			return child
		}
	}

	return nil
}

func (i *Tree) HasSelectedChild() bool {
	return i.selectedChild != nil
}

func (i *Tree) SelectChild(name string) bool {
	for _, group := range i.groups {
		if child := group.FindChild(name); child != nil {
			i.selectedChild = child
			return true
		}
	}

	return false
}

func (i *Tree) SelectedChild() argo.ChildNode {
	return i.selectedChild
}

func (i *Tree) HasIncompleteHandler() bool {
	return i.incompleteFn != nil
}

func (i *Tree) IncompleteHandler() argo.IncompleteCommandHandler[argo.CommandTree] {
	return i.incompleteFn
}

func (i *Tree) HasDescription() bool {
  return len(c.description) > 0
}

func (i *Tree) Description() string {
  return c.description
}

func (i *Tree) HasFlagGroups(includeDefault bool) bool {

}

func (i *Tree) FlagGroups(includeDefault bool) []argo.FlagGroup {

}

func (i *Tree) HasCallback() bool {
  return c.callback != nil
}

func (i *Tree) Callback() argo.CommandCallback[argo.CommandTree] {
  return c.callback
}

func (i *Tree) FindShortFlag(b byte) argo.Flag {
  for _, group := range c.flagGroups {
    if flag := group.FindShortFlag(b); flag != nil {
      return flag
    }
  }

  return nil
}

func (i *Tree) FindLongFlag(name string) argo.Flag {
  for _, group := range c.flagGroups {
    if flag := group.FindLongFlag(name); flag != nil {
      return flag
    }
  }

  return nil
}

func (i *Tree) IsHelpDisabled() bool {
  return c.disableHelp
}
