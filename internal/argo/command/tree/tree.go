package tree

import (
	"os"
	"path/filepath"

	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

type commandTree struct {
	description   string
	disableHelp   bool
	flagGroups    []argo.FlagGroup
	commandGroups []argo.CommandGroup
	warnings      *argo.WarningContext

	selectedLeaf  argo.CommandLeaf
	selectedChild argo.ChildNode[any]

	callback argo.CommandNodeCallback[argo.CommandTree]

	onIncompleteHandler argo.IncompleteCommandHandler[argo.CommandTree]
}

// region Node Impl

func (_ *commandTree) Name() string {
	return filepath.Base(os.Args[0])
}

//

func (t *commandTree) HasDescription() bool {
	return len(t.description) > 0
}

func (t *commandTree) Description() string {
	return t.description
}

//

func (t *commandTree) HasFlagGroups() bool {
	return len(t.flagGroups) > 0
}

func (t *commandTree) FlagGroups() []argo.FlagGroup {
	return t.flagGroups
}

//

func (t *commandTree) FindShortFlag(b byte) argo.Flag {
	for _, group := range t.flagGroups {
		if flag := group.FindShortFlag(b); flag != nil {
			return flag
		}
	}

	return nil
}

func (t *commandTree) FindLongFlag(name string) argo.Flag {
	for _, group := range t.FlagGroups() {
		if flag := group.FindLongFlag(name); flag != nil {
			return flag
		}
	}

	return nil
}

//

func (t *commandTree) HasCallback() bool {
	return t.callback != nil
}

func (t *commandTree) Callback() argo.CommandNodeCallback[argo.CommandTree] {
	return t.callback
}

//

func (t *commandTree) Warnings() []string {
	return t.warnings.GetWarnings()
}

func (t *commandTree) AppendWarning(warning string) {
	t.warnings.AppendWarning(warning)
}

// endregion Node Impl

// region ParentNode Impl

func (t *commandTree) HasSubcommands() bool {
	return true
}

//

func (t *commandTree) HasCommandGroups(includeDefault bool) bool {
	return len(t.commandGroups) > 1 || (includeDefault && t.hasDefaultCommandGroup())
}

func (t *commandTree) CommandGroups(includeDefault bool) []argo.CommandGroup {
	return t.commandGroups
}

func (t *commandTree) hasDefaultCommandGroup() bool {
	return t.commandGroups[0] != nil
}

//

func (t *commandTree) FindChild(name string) argo.ChildNode[any] {
	for _, group := range t.commandGroups {
		if child := group.FindChild(name); child != nil {
			return child
		}
	}

	return nil
}

//

func (t *commandTree) HasSelectedChild() bool {
	return t.selectedLeaf != nil
}

func (t *commandTree) SelectChild(name string) bool {
	for _, group := range t.commandGroups {
		if child := group.FindChild(name); child != nil {
			t.selectedChild = child
			return true
		}
	}

	return false
}

func (t *commandTree) SelectedChild() argo.ChildNode[any] {
	var child argo.ChildNode[any] = t.selectedLeaf

	for child.Parent() != t {
		child = child.Parent().(argo.ChildNode[any])
	}

	return child
}

//

func (t *commandTree) HasIncompleteHandler() bool {
	return t.onIncompleteHandler != nil
}

func (t *commandTree) IncompleteHandler() argo.IncompleteCommandHandler[argo.CommandTree] {
	return t.onIncompleteHandler
}

// endregion ParentNode Impl

// region CommandTree Impl

func (t *commandTree) SelectedCommand() argo.CommandLeaf {
	return t.selectedLeaf
}

// endregion CommandTree Impl
