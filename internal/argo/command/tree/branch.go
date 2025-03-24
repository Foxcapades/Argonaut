package tree

// WARNING:
//   This is a generated file!  Edits here will be lost!

import "github.com/foxcapades/argonaut/v3/pkg/argo"

type Branch struct {
  disableHelp bool
  description string
  flagGroups  []argo.FlagGroup
  callback    argo.CommandCallback[argo.BranchCommand]
  name    string
  parent  argo.ParentNode
  aliases []string
  commandGroups []argo.CommandGroup
	selectedChild argo.ChildNode
	incompleteFn argo.IncompleteCommandHandler[argo.BranchCommand]
}

func (i *Branch) HasDescription() bool {
  return len(i.description) > 0
}

func (i *Branch) Description() string {
  return i.description
}

func (i *Branch) HasFlagGroups() bool {
  return len(i.flagGroups) > 0
}

func (i *Branch) FlagGroups() []argo.FlagGroup {
  return i.flagGroups
}

func (i *Branch) HasCallback() bool {
  return i.callback != nil
}

func (i *Branch) Callback() argo.CommandCallback[argo.BranchCommand] {
  return i.callback
}

func (i *Branch) FindShortFlag(b byte) argo.Flag {
  for _, group := range i.flagGroups {
    if flag := group.FindShortFlag(b); flag != nil {
      return flag
    }
  }

  return nil
}

func (i *Branch) FindLongFlag(name string) argo.Flag {
  for _, group := range i.flagGroups {
    if flag := group.FindLongFlag(name); flag != nil {
      return flag
    }
  }

  return nil
}

func (i *Branch) IsHelpDisabled() bool {
  return i.disableHelp
}

func (i *Branch) Name() string {
  return i.name
}

func (i *Branch) Parent() argo.ParentNode {
  return i.parent
}

func (i *Branch) HasAliases() bool {
  return len(i.aliases) > 0
}

func (i *Branch) Aliases() []string {
  return i.aliases
}

func (i *Branch) Matches(name string) bool {
  if i.name == name {
    return true
  }

  for _, alias := range i.aliases {
    if alias == name {
      return true
    }
  }

  return false
}

func (i *Branch) HasCommandGroups() bool {
	return len(i.commandGroups) > 0
}

func (i *Branch) CommandGroups() []argo.CommandGroup {
	return i.commandGroups
}

func (i *Branch) hasDefaultCommandGroup() bool {
	return i.commandGroups[0] != nil
}

func (i *Branch) HasSubcommands() bool {
	for _, group := range i.commandGroups {
		if group.HasSubcommands() {
			return true
		}
	}

	return false
}

func (i *Branch) FindChild(name string) argo.ChildNode {
	for _, group := range i.commandGroups {
		if child := group.FindChild(name); child != nil {
			return child
		}
	}

	return nil
}

func (i *Branch) HasSelectedChild() bool {
	return i.selectedChild != nil
}

func (i *Branch) SelectChild(name string) bool {
	for _, group := range i.commandGroups {
		if child := group.FindChild(name); child != nil {
			i.selectedChild = child
			return true
		}
	}

	return false
}

func (i *Branch) SelectedChild() argo.ChildNode {
	return i.selectedChild
}

func (i *Branch) HasIncompleteHandler() bool {
	return i.incompleteFn != nil
}

func (i *Branch) IncompleteHandler() argo.IncompleteCommandHandler[argo.TreeCommand] {
	return i.incompleteFn
}
