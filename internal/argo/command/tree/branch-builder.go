package tree

// WARNING:
//   This is a generated file!  Edits here will be lost!

import (
	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func NewBranchBuilder(name string) argo.BranchCommandBuilder {
	return &CommandBranchBuilder{
		name: name,
	}
}

type CommandBranchBuilder struct {
	comGroups    []argo.CommandGroupBuilder
	incompleteFn argo.IncompleteCommandHandler[argo.BranchCommand]
	name         string
	aliases      []string
	parent       argo.ParentNodeBuilder
	disableHelp  bool
	description  string
	flagGroups   []argo.FlagGroupBuilder
	callback     argo.CommandCallback[argo.BranchCommand]
}

func (i *CommandBranchBuilder) WithCommandGroup(group argo.CommandGroupBuilder) argo.BranchCommandBuilder {
	i.comGroups = append(i.comGroups, group)
	return i
}

func (i *CommandBranchBuilder) WithCommandGroups(groups ...argo.CommandGroupBuilder) argo.BranchCommandBuilder {
	i.comGroups = append(i.comGroups, groups...)
	return i
}

func (i *CommandBranchBuilder) HasCommandGroups(includeDefault bool) bool {
	return (includeDefault && len(i.comGroups) > 0) ||
		(i.hasDefaultCommandGroup() && len(i.comGroups) > 1) ||
		len(i.comGroups) > 0
}

func (i *CommandBranchBuilder) CommandGroups(includeDefault bool) []argo.CommandGroupBuilder {
	if !includeDefault && i.hasDefaultCommandGroup() {
		return i.comGroups[1:]
	}

	return i.comGroups
}

func (i *CommandBranchBuilder) hasDefaultCommandGroup() bool {
	return len(i.comGroups) > 0 && i.comGroups[0].Name() == DefaultCommandGroupName
}

func (i *CommandBranchBuilder) WithBranch(branch argo.BranchCommandBuilder) argo.BranchCommandBuilder {
	i.comGroups = EnsureDefaultCommandGroup(i.comGroups)
	i.comGroups[0].WithBranch(branch)
	return i
}

func (i *CommandBranchBuilder) WithBranches(branches ...argo.BranchCommandBuilder) argo.BranchCommandBuilder {
	i.comGroups = EnsureDefaultCommandGroup(i.comGroups)
	i.comGroups[0].WithBranches(branches...)
	return i
}

func (i *CommandBranchBuilder) WithLeaf(leaf argo.LeafCommandBuilder) argo.BranchCommandBuilder {
	i.comGroups = EnsureDefaultCommandGroup(i.comGroups)
	i.comGroups[0].WithLeaf(leaf)
	return i
}

func (i *CommandBranchBuilder) WithLeaves(leaves ...argo.LeafCommandBuilder) argo.BranchCommandBuilder {
	i.comGroups = EnsureDefaultCommandGroup(i.comGroups)
	i.comGroups[0].WithLeaves(leaves...)
	return i
}

func (i *CommandBranchBuilder) HasSubcommands() bool {
	for _, group := range i.comGroups {
		if group.HasSubcommands() {
			return true
		}
	}

	return false
}

func (i *CommandBranchBuilder) WithIncompleteHandler(handler argo.IncompleteCommandHandler[argo.BranchCommand]) argo.BranchCommandBuilder {
	i.incompleteFn = handler
	return i
}

func (i *CommandBranchBuilder) HasIncompleteHandler() bool {
	return i.incompleteFn != nil
}

func (i *CommandBranchBuilder) IncompleteHandler() argo.IncompleteCommandHandler[argo.BranchCommand] {
	return i.incompleteFn
}
func (i *CommandBranchBuilder) Name() string {
	return i.name
}

func (i *CommandBranchBuilder) WithAlias(alias string) argo.BranchCommandBuilder {
	i.aliases = append(i.aliases, alias)
	return i
}

func (i *CommandBranchBuilder) WithAliases(aliases ...string) argo.BranchCommandBuilder {
	i.aliases = append(i.aliases, aliases...)
	return i
}

func (i *CommandBranchBuilder) Aliases() []string {
	return i.aliases
}

func (i *CommandBranchBuilder) HasAliases() bool {
	return len(i.aliases) > 0
}

func (i *CommandBranchBuilder) SetParentNode(parent argo.ParentNodeBuilder) argo.BranchCommandBuilder {
	i.parent = parent
	return i
}

func (i *CommandBranchBuilder) ParentNode() argo.ParentNodeBuilder {
	return i.parent
}
func (i *CommandBranchBuilder) WithDescription(desc string) argo.BranchCommandBuilder {
	i.description = desc
	return i
}

func (i *CommandBranchBuilder) HasDescription() bool {
	return len(i.description) > 0
}

func (i *CommandBranchBuilder) Description() string {
	return i.description
}

func (i *CommandBranchBuilder) WithHelpDisabled() argo.BranchCommandBuilder {
	i.disableHelp = true
	return i
}

func (i *CommandBranchBuilder) IsHelpDisabled() bool {
	return i.disableHelp
}

func (i *CommandBranchBuilder) WithFlagGroup(group argo.FlagGroupBuilder) argo.BranchCommandBuilder {
	i.flagGroups = append(i.flagGroups, group)
	return i
}

func (i *CommandBranchBuilder) WithFlagGroups(groups ...argo.FlagGroupBuilder) argo.BranchCommandBuilder {
	i.flagGroups = append(i.flagGroups, groups...)
	return i
}

func (i *CommandBranchBuilder) HasFlagGroups(includeDefault bool) bool {
	return (includeDefault && len(i.flagGroups) > 0) ||
		(i.hasDefaultFlagGroup() && len(i.flagGroups) > 1) ||
		len(i.flagGroups) > 0
}

func (i *CommandBranchBuilder) FlagGroups(includeDefault bool) []argo.FlagGroupBuilder {
	if !includeDefault && i.hasDefaultFlagGroup() {
		return i.flagGroups[1:]
	}

	return i.flagGroups
}

func (i *CommandBranchBuilder) hasDefaultFlagGroup() bool {
	return len(i.flagGroups) > 0 && flag.IsDefaultGroup(i.flagGroups[0])
}

func (i *CommandBranchBuilder) WithFlag(fb argo.FlagBuilder) argo.BranchCommandBuilder {
	i.flagGroups = flag.EnsureDefaultGroup(i.flagGroups)
	i.flagGroups[0].WithFlag(fb)
	return i
}

func (i *CommandBranchBuilder) WithFlags(flags ...argo.FlagBuilder) argo.BranchCommandBuilder {
	i.flagGroups = flag.EnsureDefaultGroup(i.flagGroups)
	i.flagGroups[0].WithFlags(flags...)
	return i
}

func (i *CommandBranchBuilder) HasFlags() bool {
	for _, group := range i.flagGroups {
		if group != nil && group.HasFlags() {
			return true
		}
	}

	return false
}

func (i *CommandBranchBuilder) WithCallback(callback argo.CommandCallback[argo.BranchCommand]) argo.BranchCommandBuilder {
	i.callback = callback
	return i
}

func (i *CommandBranchBuilder) HasCallback() bool {
	return i.callback != nil
}

func (i *CommandBranchBuilder) Callback() argo.CommandCallback[argo.BranchCommand] {
	return i.callback
}
