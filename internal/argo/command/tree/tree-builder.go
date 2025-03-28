package tree

// WARNING:
//   This is a generated file!  Edits here will be lost!

import (
	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func NewBuilder() argo.TreeCommandBuilder {
	return new(TreeCommandBuilder)
}

type TreeCommandBuilder struct {
	comGroups    []argo.CommandGroupBuilder
	incompleteFn argo.IncompleteCommandHandler[argo.TreeCommand]
	disableHelp  bool
	description  string
	flagGroups   []argo.FlagGroupBuilder
	callback     argo.CommandCallback[argo.TreeCommand]
}

func (i *TreeCommandBuilder) WithCommandGroup(group argo.CommandGroupBuilder) argo.TreeCommandBuilder {
	i.comGroups = append(i.comGroups, group)
	return i
}

func (i *TreeCommandBuilder) WithCommandGroups(groups ...argo.CommandGroupBuilder) argo.TreeCommandBuilder {
	i.comGroups = append(i.comGroups, groups...)
	return i
}

func (i *TreeCommandBuilder) HasCommandGroups(includeDefault bool) bool {
	return (includeDefault && len(i.comGroups) > 0) ||
		(i.hasDefaultCommandGroup() && len(i.comGroups) > 1) ||
		len(i.comGroups) > 0
}

func (i *TreeCommandBuilder) CommandGroups(includeDefault bool) []argo.CommandGroupBuilder {
	if !includeDefault && i.hasDefaultCommandGroup() {
		return i.comGroups[1:]
	}

	return i.comGroups
}

func (i *TreeCommandBuilder) hasDefaultCommandGroup() bool {
	return len(i.comGroups) > 0 && i.comGroups[0].Name() == DefaultCommandGroupName
}

func (i *TreeCommandBuilder) WithBranch(branch argo.BranchCommandBuilder) argo.TreeCommandBuilder {
	i.comGroups = EnsureDefaultCommandGroup(i.comGroups)
	i.comGroups[0].WithBranch(branch)
	return i
}

func (i *TreeCommandBuilder) WithBranches(branches ...argo.BranchCommandBuilder) argo.TreeCommandBuilder {
	i.comGroups = EnsureDefaultCommandGroup(i.comGroups)
	i.comGroups[0].WithBranches(branches...)
	return i
}

func (i *TreeCommandBuilder) WithLeaf(leaf argo.LeafCommandBuilder) argo.TreeCommandBuilder {
	i.comGroups = EnsureDefaultCommandGroup(i.comGroups)
	i.comGroups[0].WithLeaf(leaf)
	return i
}

func (i *TreeCommandBuilder) WithLeaves(leaves ...argo.LeafCommandBuilder) argo.TreeCommandBuilder {
	i.comGroups = EnsureDefaultCommandGroup(i.comGroups)
	i.comGroups[0].WithLeaves(leaves...)
	return i
}

func (i *TreeCommandBuilder) HasSubcommands() bool {
	for _, group := range i.comGroups {
		if group.HasSubcommands() {
			return true
		}
	}

	return false
}

func (i *TreeCommandBuilder) WithIncompleteHandler(handler argo.IncompleteCommandHandler[argo.TreeCommand]) argo.TreeCommandBuilder {
	i.incompleteFn = handler
	return i
}

func (i *TreeCommandBuilder) HasIncompleteHandler() bool {
	return i.incompleteFn != nil
}

func (i *TreeCommandBuilder) IncompleteHandler() argo.IncompleteCommandHandler[argo.TreeCommand] {
	return i.incompleteFn
}
func (i *TreeCommandBuilder) WithDescription(desc string) argo.TreeCommandBuilder {
	i.description = desc
	return i
}

func (i *TreeCommandBuilder) HasDescription() bool {
	return len(i.description) > 0
}

func (i *TreeCommandBuilder) Description() string {
	return i.description
}

func (i *TreeCommandBuilder) WithHelpDisabled() argo.TreeCommandBuilder {
	i.disableHelp = true
	return i
}

func (i *TreeCommandBuilder) IsHelpDisabled() bool {
	return i.disableHelp
}

func (i *TreeCommandBuilder) WithFlagGroup(group argo.FlagGroupBuilder) argo.TreeCommandBuilder {
	i.flagGroups = append(i.flagGroups, group)
	return i
}

func (i *TreeCommandBuilder) WithFlagGroups(groups ...argo.FlagGroupBuilder) argo.TreeCommandBuilder {
	i.flagGroups = append(i.flagGroups, groups...)
	return i
}

func (i *TreeCommandBuilder) HasFlagGroups(includeDefault bool) bool {
	return (includeDefault && len(i.flagGroups) > 0) ||
		(i.hasDefaultFlagGroup() && len(i.flagGroups) > 1) ||
		len(i.flagGroups) > 0
}

func (i *TreeCommandBuilder) FlagGroups(includeDefault bool) []argo.FlagGroupBuilder {
	if !includeDefault && i.hasDefaultFlagGroup() {
		return i.flagGroups[1:]
	}

	return i.flagGroups
}

func (i *TreeCommandBuilder) hasDefaultFlagGroup() bool {
	return i.flagGroups[0].Size() > 0 && flag.IsDefaultGroup(i.flagGroups[0])
}

func (i *TreeCommandBuilder) WithFlag(fb argo.FlagBuilder) argo.TreeCommandBuilder {
	i.flagGroups = flag.EnsureDefaultGroup(i.flagGroups)
	i.flagGroups[0].WithFlag(fb)
	return i
}

func (i *TreeCommandBuilder) WithFlags(flags ...argo.FlagBuilder) argo.TreeCommandBuilder {
	i.flagGroups = flag.EnsureDefaultGroup(i.flagGroups)
	i.flagGroups[0].WithFlags(flags...)
	return i
}

func (i *TreeCommandBuilder) HasFlags() bool {
	for _, group := range i.flagGroups {
		if group != nil && group.HasFlags() {
			return true
		}
	}

	return false
}

func (i *TreeCommandBuilder) WithCallback(callback argo.CommandCallback[argo.TreeCommand]) argo.TreeCommandBuilder {
	i.callback = callback
	return i
}

func (i *TreeCommandBuilder) HasCallback() bool {
	return i.callback != nil
}

func (i *TreeCommandBuilder) Callback() argo.CommandCallback[argo.TreeCommand] {
	return i.callback
}
