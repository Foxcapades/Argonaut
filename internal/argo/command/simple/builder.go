package command

import (
	"fmt"
	"os"

	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/internal/chars"
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/internal/xerr"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func NewBuilder() argo.CommandBuilder {
	return &commandBuilder{
		flagGroups: []argo.FlagGroupBuilder{flag.NewGroupBuilder(chars.DefaultGroupName)},
	}
}

type commandBuilder struct {
	description string
	unmapLabel  string
	flagGroups  []argo.FlagGroupBuilder
	arguments   []argo.ArgumentBuilder
	disableHelp bool
	callback    argo.CommandCallback
}

func (b *commandBuilder) WithDescription(desc string) argo.CommandBuilder {
	b.description = desc
	return b
}

func (b *commandBuilder) WithHelpDisabled() argo.CommandBuilder {
	b.disableHelp = true
	return b
}

func (b *commandBuilder) WithFlagGroup(group argo.FlagGroupBuilder) argo.CommandBuilder {
	b.flagGroups = append(b.flagGroups, group)
	return b
}

func (b *commandBuilder) WithUnmappedLabel(label string) argo.CommandBuilder {
	b.unmapLabel = label
	return b
}

func (b *commandBuilder) WithFlag(flag argo.FlagBuilder) argo.CommandBuilder {
	b.flagGroups[0].WithFlag(flag)
	return b
}

func (b *commandBuilder) WithArgument(arg argo.ArgumentBuilder) argo.CommandBuilder {
	b.arguments = append(b.arguments, arg)
	return b
}

func (b *commandBuilder) WithCallback(cb argo.CommandCallback) argo.CommandBuilder {
	b.callback = cb
	return b
}

func (b *commandBuilder) Parse(args []string) (argo.Command, error) {
	ctx := new(argo.WarningContext)
	if cmd, err := b.Build(ctx); err != nil {
		return nil, err
	} else {
		if err = newCommandInterpreter(args, cmd).Run(); err != nil {
			return nil, err
		}

		return cmd, nil
	}
}

func (b *commandBuilder) MustParse(args []string) argo.Command {
	return utils.MustReturn(b.Parse(args))
}

func (b *commandBuilder) Build(ctx *argo.WarningContext) (argo.Command, error) {

}
