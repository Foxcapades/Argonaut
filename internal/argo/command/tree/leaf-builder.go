package tree

// WARNING:
//   This is a generated file!  Edits here will be lost!

import (
	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func NewLeafBuilder(name string) argo.LeafCommandBuilder {
	return &LeafCommandBuilder{
		name: name,
	}
}

type LeafCommandBuilder struct {
	name    string
	aliases []string
	parent  argo.ParentNodeBuilder

	disableHelp bool
	description string
	flagGroups  []argo.FlagGroupBuilder
	callback    argo.CommandCallback[argo.LeafCommand]

	arguments     []argo.ArgumentBuilder
	unmapped      []string
	unmappedLabel string
}

func (i *LeafCommandBuilder) Name() string {
	return i.name
}

func (i *LeafCommandBuilder) WithAlias(alias string) argo.LeafCommandBuilder {
	i.aliases = append(i.aliases, alias)
	return i
}

func (i *LeafCommandBuilder) WithAliases(aliases ...string) argo.LeafCommandBuilder {
	i.aliases = append(i.aliases, aliases...)
	return i
}

func (i *LeafCommandBuilder) Aliases() []string {
	return i.aliases
}

func (i *LeafCommandBuilder) HasAliases() bool {
	return len(i.aliases) > 0
}

func (i *LeafCommandBuilder) SetParentNode(parent argo.ParentNodeBuilder) argo.LeafCommandBuilder {
	i.parent = parent
	return i
}

func (i *LeafCommandBuilder) ParentNode() argo.ParentNodeBuilder {
	return i.parent
}

func (i *LeafCommandBuilder) WithDescription(desc string) argo.LeafCommandBuilder {
	i.description = desc
	return i
}

func (i *LeafCommandBuilder) HasDescription() bool {
	return len(i.description) > 0
}

func (i *LeafCommandBuilder) Description() string {
	return i.description
}

func (i *LeafCommandBuilder) WithHelpDisabled() argo.LeafCommandBuilder {
	i.disableHelp = true
	return i
}

func (i *LeafCommandBuilder) IsHelpDisabled() bool {
	return i.disableHelp
}

func (i *LeafCommandBuilder) WithFlagGroup(group argo.FlagGroupBuilder) argo.LeafCommandBuilder {
	i.flagGroups = append(i.flagGroups, group)
	return i
}

func (i *LeafCommandBuilder) WithFlagGroups(groups ...argo.FlagGroupBuilder) argo.LeafCommandBuilder {
	i.flagGroups = append(i.flagGroups, groups...)
	return i
}

func (i *LeafCommandBuilder) HasFlagGroups(includeDefault bool) bool {
	return (includeDefault && len(i.flagGroups) > 0) ||
		(i.hasDefaultFlagGroup() && len(i.flagGroups) > 1) ||
		len(i.flagGroups) > 0
}

func (i *LeafCommandBuilder) FlagGroups(includeDefault bool) []argo.FlagGroupBuilder {
	if !includeDefault && i.hasDefaultFlagGroup() {
		return i.flagGroups[1:]
	}

	return i.flagGroups
}

func (i *LeafCommandBuilder) hasDefaultFlagGroup() bool {
	return i.flagGroups[0].Size() > 0 && flag.IsDefaultGroup(i.flagGroups[0])
}

func (i *LeafCommandBuilder) WithFlag(fb argo.FlagBuilder) argo.LeafCommandBuilder {
	i.flagGroups = flag.EnsureDefaultGroup(i.flagGroups)
	i.flagGroups[0].WithFlag(fb)
	return i
}

func (i *LeafCommandBuilder) WithFlags(flags ...argo.FlagBuilder) argo.LeafCommandBuilder {
	i.flagGroups = flag.EnsureDefaultGroup(i.flagGroups)
	i.flagGroups[0].WithFlags(flags...)
	return i
}

func (i *LeafCommandBuilder) HasFlags() bool {
	for _, group := range i.flagGroups {
		if group != nil && group.HasFlags() {
			return true
		}
	}

	return false
}

func (i *LeafCommandBuilder) WithCallback(callback argo.CommandCallback[argo.LeafCommand]) argo.LeafCommandBuilder {
	i.callback = callback
	return i
}

func (i *LeafCommandBuilder) HasCallback() bool {
	return i.callback != nil
}

func (i *LeafCommandBuilder) Callback() argo.CommandCallback[argo.LeafCommand] {
	return i.callback
}

func (i *LeafCommandBuilder) WithArgument(argument argo.ArgumentBuilder) argo.LeafCommandBuilder {
	i.arguments = append(i.arguments, argument)
	return i
}

func (i *LeafCommandBuilder) WithArguments(arguments ...argo.ArgumentBuilder) argo.LeafCommandBuilder {
	i.arguments = append(i.arguments, arguments...)
	return i
}

func (i *LeafCommandBuilder) HasArguments() bool {
	return len(i.arguments) > 0
}

func (i *LeafCommandBuilder) Arguments() []argo.ArgumentBuilder {
	return i.arguments
}

func (i *LeafCommandBuilder) HasUnmappedInputs() bool {
	return len(i.unmapped) > 0
}

func (i *LeafCommandBuilder) UnmappedInputs() []string {
	return i.unmapped
}

func (i *LeafCommandBuilder) AppendUnmappedInput(input string) {
	i.unmapped = append(i.unmapped, input)
}

func (i *LeafCommandBuilder) WithUnmappedInputLabel(label string) argo.LeafCommandBuilder {
	i.unmappedLabel = label
	return i
}

func (i *LeafCommandBuilder) HasUnmappedInputLabel() bool {
	return len(i.unmappedLabel) > 0
}

func (i *LeafCommandBuilder) UnmappedInputLabel() string {
	return i.unmappedLabel
}
