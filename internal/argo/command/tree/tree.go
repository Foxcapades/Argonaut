package tree

// WARNING:
//   This is a generated file!  Edits here will be lost!

import (
	"os"
	"path/filepath"

	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

type Tree struct {
	disableHelp   bool
	description   string
	flagGroups    []argo.FlagGroup
	callback      argo.CommandCallback[argo.TreeCommand]
	commandGroups []argo.CommandGroup
	selectedChild argo.ChildNode
	incompleteFn  argo.IncompleteCommandHandler[argo.TreeCommand]
	options       argo.TreeCommandOptions
}

func (_ *Tree) Name() string {
	return filepath.Base(os.Args[0])
}

func (i *Tree) SelectedCommand() argo.LeafCommand {
	child := i.selectedChild

	for child != nil {
		if leaf, ok := child.(argo.LeafCommand); ok {
			return leaf
		}

		child = child.(argo.ParentNode).SelectedChild()
	}

	return nil
}

func (i *Tree) HasCommandGroups() bool {
	return len(i.commandGroups) > 0
}

func (i *Tree) CommandGroups() []argo.CommandGroup {
	return i.commandGroups
}

func (i *Tree) hasDefaultCommandGroup() bool {
	return i.commandGroups[0] != nil
}

func (i *Tree) HasSubcommands() bool {
	for _, group := range i.commandGroups {
		if group.HasSubcommands() {
			return true
		}
	}

	return false
}

func (i *Tree) FindChild(name string) argo.ChildNode {
	for _, group := range i.commandGroups {
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
	for _, group := range i.commandGroups {
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

func (i *Tree) IncompleteHandler() argo.IncompleteCommandHandler[argo.TreeCommand] {
	return i.incompleteFn
}

func (i *Tree) HasDescription() bool {
	return len(i.description) > 0
}

func (i *Tree) Description() string {
	return i.description
}

func (i *Tree) HasFlagGroups() bool {
	return len(i.flagGroups) > 0
}

func (i *Tree) FlagGroups() []argo.FlagGroup {
	return i.flagGroups
}

func (i *Tree) HasCallback() bool {
	return i.callback != nil
}

func (i *Tree) Callback() argo.CommandCallback[argo.TreeCommand] {
	return i.callback
}

func (i *Tree) IsHelpDisabled() bool {
	return i.disableHelp
}

func (i *Tree) Options() argo.TreeCommandOptions {
	return i.options
}

func (i *Tree) FindShortFlag(b byte) argo.Flag {
	for _, group := range i.flagGroups {
		if flag := group.FindShortFlag(b); flag != nil {
			return flag
		}
	}

	return nil
}

func (i *Tree) FindLongFlag(name string) argo.Flag {
	for _, group := range i.flagGroups {
		if flag := group.FindLongFlag(name); flag != nil {
			return flag
		}
	}

	return nil
}

func (i *Tree) FindShortFlagRecursive(c byte) argo.Flag {
	if selected := i.SelectedCommand(); selected != nil {
		return selected.FindShortFlagRecursive(c)
	}

	return nil
}

func (i *Tree) FindLongFlagRecursive(name string) argo.Flag {
	if selected := i.SelectedCommand(); selected != nil {
		return selected.FindLongFlagRecursive(name)
	}

	return nil
}
