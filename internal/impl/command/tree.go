package command

import (
	"os"
	"path/filepath"

	"github.com/Foxcapades/Argonaut/internal/consts"
	"github.com/Foxcapades/Argonaut/pkg/argo"
)

type tree struct {
	description   string
	disableHelp   bool
	flagGroups    []argo.FlagGroup
	commandGroups []argo.CommandGroup
	selected      argo.CommandLeaf
	callback      argo.CommandTreeCallback
}

func (_ tree) Name() string {
	return filepath.Base(os.Args[0])
}

func (_ tree) Parent() argo.CommandNode {
	return nil
}

func (_ tree) HasParent() bool {
	return false
}

func (t tree) Description() string {
	return t.description
}

func (t tree) HasDescription() bool {
	return len(t.description) > 0
}

func (t tree) FlagGroups() []argo.FlagGroup {
	return t.flagGroups
}

func (t *tree) HasFlagGroups() bool {
	return len(t.flagGroups) > 0
}

func (t tree) CommandGroups() []argo.CommandGroup {
	return t.commandGroups
}

func (t tree) HasCallback() bool {
	return t.callback != nil
}

func (t tree) RunCallback() {
	if t.callback != nil {
		t.callback(&t)
	}
}

func (t tree) HasCustomCommandGroups() bool {
	return len(t.commandGroups) > 1 || t.commandGroups[0].Name() != consts.DefaultGroupName
}

func (t tree) SelectedCommand() argo.CommandLeaf {
	return t.selected
}

func (t *tree) SelectCommand(leaf argo.CommandLeaf) {
	t.selected = leaf
}

func (t tree) IsHelpDisabled() bool {
	return t.disableHelp
}

func (t tree) FindChild(name string) argo.CommandNode {
	for _, group := range t.commandGroups {
		if child := group.FindChild(name); child != nil {
			return child
		}
	}

	return nil
}

func (t tree) FindShortFlag(b byte) argo.Flag {
	for _, group := range t.flagGroups {
		if flag := group.FindShortFlag(b); flag != nil {
			return flag
		}
	}

	return nil
}

func (t tree) FindLongFlag(name string) argo.Flag {
	for _, group := range t.FlagGroups() {
		if flag := group.FindLongFlag(name); flag != nil {
			return flag
		}
	}

	return nil
}

func (t tree) TryFlag(ref argo.FlagRef) (bool, error) {
	for _, flagGroup := range t.FlagGroups() {
		if ok, err := flagGroup.TryFlag(ref); ok || err != nil {
			return ok, err
		}
	}

	return false, nil
}
