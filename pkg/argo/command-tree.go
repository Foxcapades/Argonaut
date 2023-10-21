package argo

import (
	"os"
	"path/filepath"

	"github.com/Foxcapades/Argonaut/internal/util"
)

// CommandTree represents the root of a tree of subcommands.
//
// The command tree consists of branch and leaf nodes.  The branch nodes can be
// thought of as categories for containing sub-branches and/or leaves.  Leaf
// nodes are the actual callable command implementations.
//
// All levels of the command tree accept flags, with sub-node flags taking
// priority over parent node flags on flag collision.  Leaf nodes, however, are
// the only nodes that accept positional arguments, or passthroughs.
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
type CommandTree interface {
	CommandNode
	CommandParent

	// SelectedCommand returns the leaf command that was selected in the CLI call.
	SelectedCommand() CommandLeaf

	selectCommand(leaf CommandLeaf)

	hasCallback() bool

	executeCallback()
}

type CommandTreeCallback = func(com CommandTree)

type commandTree struct {
	description   string
	disableHelp   bool
	flagGroups    []FlagGroup
	commandGroups []CommandGroup
	selected      CommandLeaf
	callback      CommandTreeCallback
	warnings      *WarningContext
}

func (_ commandTree) Name() string {
	return filepath.Base(os.Args[0])
}

func (_ commandTree) Parent() CommandNode {
	return nil
}

func (_ commandTree) HasParent() bool {
	return false
}

func (t commandTree) Description() string {
	return t.description
}

func (t commandTree) HasDescription() bool {
	return len(t.description) > 0
}

func (t commandTree) FlagGroups() []FlagGroup {
	return t.flagGroups
}

func (t *commandTree) HasFlagGroups() bool {
	return len(t.flagGroups) > 0
}

func (t commandTree) CommandGroups() []CommandGroup {
	return t.commandGroups
}

func (t commandTree) hasCallback() bool {
	return t.callback != nil
}

func (t commandTree) executeCallback() {
	if t.callback != nil {
		t.callback(&t)
	}
}

func (t commandTree) SelectedCommand() CommandLeaf {
	return t.selected
}

func (t *commandTree) selectCommand(leaf CommandLeaf) {
	t.selected = leaf
}

func (t commandTree) FindChild(name string) CommandChild {
	for _, group := range t.commandGroups {
		if child := group.FindChild(name); child != nil {
			return child
		}
	}

	return nil
}

func (t commandTree) FindShortFlag(b byte) Flag {
	for _, group := range t.flagGroups {
		if flag := group.FindShortFlag(b); flag != nil {
			return flag
		}
	}

	return nil
}

func (t commandTree) FindLongFlag(name string) Flag {
	for _, group := range t.FlagGroups() {
		if flag := group.FindLongFlag(name); flag != nil {
			return flag
		}
	}

	return nil
}

func (t commandTree) onIncomplete() {
	util.Must(comTreeRenderer{}.RenderHelp(&t, os.Stdout))
	os.Exit(1)
}

func (t *commandTree) Warnings() []string {
	return t.warnings.GetWarnings()
}

func (t *commandTree) AppendWarning(warning string) {
	t.warnings.appendWarning(warning)
}
