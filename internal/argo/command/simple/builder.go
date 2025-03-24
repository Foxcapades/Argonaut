package command

import (
	"github.com/foxcapades/argonaut/v3/internal/argo/command/common"
	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func NewBuilder() argo.CommandBuilder {
	return &commandBuilder{}
}

type commandBuilder struct {
	description string
	unmapLabel  string
	flagGroups  []argo.FlagGroupBuilder
	arguments   []argo.ArgumentBuilder
	disableHelp bool
	callback    argo.CommandCallback[argo.Command]
}

//

func (b *commandBuilder) WithDescription(desc string) argo.CommandBuilder {
	b.description = desc
	return b
}

func (b *commandBuilder) HasDescription() bool {
	return len(b.description) > 0
}

func (b *commandBuilder) Description() string {
	return b.description
}

//

func (b *commandBuilder) WithHelpDisabled() argo.CommandBuilder {
	b.disableHelp = true
	return b
}

func (b *commandBuilder) IsHelpDisabled() bool {
	return b.disableHelp
}

//

func (b *commandBuilder) WithFlagGroup(group argo.FlagGroupBuilder) argo.CommandBuilder {
	b.flagGroups = append(b.flagGroups, group)
	return b
}

func (b *commandBuilder) WithFlagGroups(groups ...argo.FlagGroupBuilder) argo.CommandBuilder {
	b.flagGroups = append(b.flagGroups, groups...)
	return b
}

func (b *commandBuilder) HasFlagGroups(includeDefault bool) bool {
	return (includeDefault && len(b.flagGroups) > 0) ||
		(b.hasDefaultFlagGroup() && len(b.flagGroups) > 1) ||
		len(b.flagGroups) > 0
}

func (b *commandBuilder) FlagGroups(includeDefault bool) []argo.FlagGroupBuilder {
	if includeDefault || !b.hasDefaultFlagGroup() || len(b.flagGroups) == 0 {
		return b.flagGroups
	}

	return b.flagGroups[1:]
}

func (b *commandBuilder) hasDefaultFlagGroup() bool {
	return len(b.flagGroups) > 0 && flag.IsDefaultGroup(b.flagGroups[0])
}

//

func (b *commandBuilder) WithFlag(flag argo.FlagBuilder) argo.CommandBuilder {
	b.flagGroups = common.EnsureDefaultFlagGroup(b.flagGroups)
	b.flagGroups[0].WithFlag(flag)
	return b
}

func (b *commandBuilder) WithFlags(flags ...argo.FlagBuilder) argo.CommandBuilder {
	b.flagGroups = common.EnsureDefaultFlagGroup(b.flagGroups)
	b.flagGroups[0].WithFlags(flags...)
	return b
}

func (b *commandBuilder) HasFlags() bool {
	for _, group := range b.flagGroups {
		if group.HasFlags() {
			return true
		}
	}

	return false
}

//

func (b *commandBuilder) WithArgument(arg argo.ArgumentBuilder) argo.CommandBuilder {
	b.arguments = append(b.arguments, arg)
	return b
}

func (b *commandBuilder) WithArguments(args ...argo.ArgumentBuilder) argo.CommandBuilder {
	b.arguments = append(b.arguments, args...)
	return b
}

func (b *commandBuilder) HasArguments() bool {
	return len(b.arguments) > 0
}

func (b *commandBuilder) Arguments() []argo.ArgumentBuilder {
	return b.arguments
}

//

func (b *commandBuilder) WithUnmappedInputLabel(label string) argo.CommandBuilder {
	b.unmapLabel = label
	return b
}

func (b *commandBuilder) HasUnmappedInputLabel() bool {
	return len(b.unmapLabel) > 0
}

func (b *commandBuilder) UnmappedInputLabel() string {
	return b.unmapLabel
}

//

func (b *commandBuilder) WithCallback(cb argo.CommandCallback[argo.Command]) argo.CommandBuilder {
	b.callback = cb
	return b
}

func (b *commandBuilder) HasCallback() bool {
	return b.callback != nil
}

func (b *commandBuilder) Callback() argo.CommandCallback[argo.Command] {
	return b.callback
}
