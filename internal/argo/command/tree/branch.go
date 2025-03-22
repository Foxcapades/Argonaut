package tree

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
  return len(c.description) > 0
}

func (i *Branch) Description() string {
  return c.description
}

func (i *Branch) HasFlagGroups(includeDefault bool) bool {

}

func (i *Branch) FlagGroups(includeDefault bool) []argo.FlagGroup {

}

func (i *Branch) HasCallback() bool {
  return c.callback != nil
}

func (i *Branch) Callback() argo.CommandCallback[argo.BranchCommand] {
  return c.callback
}

func (i *Branch) FindShortFlag(b byte) argo.Flag {
  for _, group := range c.flagGroups {
    if flag := group.FindShortFlag(b); flag != nil {
      return flag
    }
  }

  return nil
}

func (i *Branch) FindLongFlag(name string) argo.Flag {
  for _, group := range c.flagGroups {
    if flag := group.FindLongFlag(name); flag != nil {
      return flag
    }
  }

  return nil
}

func (i *Branch) IsHelpDisabled() bool {
  return c.disableHelp
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

func (i *Branch) HasCommandGroups(includeDefault bool) bool {
	return len(i.groups) > 1 || (includeDefault && i.hasDefaultCommandGroup())
}

func (i *Branch) CommandGroups(includeDefault bool) []argo.CommandGroup {
	if !i.HasCommandGroups(includeDefault) {
		return nil
	}

	if includeDefault && i.hasDefaultCommandGroup() {
		return i.groups
	}

	return i.groups[1:]
}

func (i *Branch) hasDefaultCommandGroup() bool {
	return i.groups[0] != nil
}

func (i *Branch) FindChild(name string) argo.ChildNode {
	for _, group := range i.groups {
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
	for _, group := range i.groups {
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

func (i *Branch) IncompleteHandler() argo.IncompleteCommandHandler[argo.CommandTree] {
	return i.incompleteFn
}
