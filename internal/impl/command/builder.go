package command

import (
	"github.com/Foxcapades/Argonaut/internal/consts"
	"github.com/Foxcapades/Argonaut/internal/impl/flag"
	"github.com/Foxcapades/Argonaut/internal/interpret"
	"github.com/Foxcapades/Argonaut/internal/util"
	"github.com/Foxcapades/Argonaut/internal/xerr"
	"github.com/Foxcapades/Argonaut/pkg/argo"
)

func Builder() argo.CommandBuilder {
	return &builder{
		flagGroups: []argo.FlagGroupBuilder{flag.GroupBuilder(consts.DefaultGroupName)},
	}
}

type builder struct {
	description string
	unmapLabel  string
	flagGroups  []argo.FlagGroupBuilder
	arguments   []argo.ArgumentBuilder
}

func (b *builder) WithDescription(desc string) argo.CommandBuilder {
	b.description = desc
	return b
}

func (b builder) HasDescription() bool {
	return len(b.description) > 0
}

func (b builder) GetDescription() string {
	return b.description
}

func (b *builder) WithFlagGroup(group argo.FlagGroupBuilder) argo.CommandBuilder {
	b.flagGroups = append(b.flagGroups, group)
	return b
}

func (b builder) HasFlagGroups() bool {
	return len(b.flagGroups) > 1
}

func (b builder) GetFlagGroups() []argo.FlagGroupBuilder {
	return b.flagGroups
}

func (b *builder) WithUnmappedLabel(label string) argo.CommandBuilder {
	b.unmapLabel = label
	return b
}

func (b builder) HasUnmappedLabel() bool {
	return len(b.unmapLabel) > 0
}

func (b builder) GetUnmappedLabel() string {
	return b.unmapLabel
}

func (b *builder) WithFlag(flag argo.FlagBuilder) argo.CommandBuilder {
	b.flagGroups[0].WithFlag(flag)
	return b
}

func (b builder) HasFlags() bool {
	for _, group := range b.flagGroups {
		if group.HasFlags() {
			return true
		}
	}

	return false
}

func (b *builder) WithArgument(arg argo.ArgumentBuilder) argo.CommandBuilder {
	b.arguments = append(b.arguments, arg)
	return b
}

func (b builder) HasArguments() bool {
	return len(b.arguments) > 0
}

func (b builder) GetArguments() []argo.ArgumentBuilder {
	return b.arguments
}

func (b builder) Parse(args []string) (argo.Command, error) {
	if cmd, err := b.build(); err != nil {
		return nil, err
	} else {
		if err = interpret.CommandInterpreter(args, cmd).Run(); err != nil {
			return nil, err
		}

		return cmd, nil
	}
}

func (b builder) MustParse(args []string) argo.Command {
	return util.MustReturn(b.Parse(args))
}

func (b builder) build() (argo.Command, error) {
	errs := xerr.NewMultiError()

	flagGroups := make([]argo.FlagGroup, 0, len(b.flagGroups))
	for _, builder := range b.flagGroups {
		if builder.HasFlags() {
			if group, err := builder.Build(); err != nil {
				errs.AppendError(err)
			} else {
				flagGroups = append(flagGroups, group)
			}
		}
	}

	arguments := make([]argo.Argument, 0, len(b.arguments))
	for _, builder := range b.arguments {
		if arg, err := builder.Build(); err != nil {
			errs.AppendError(err)
		} else {
			arguments = append(arguments, arg)
		}
	}

	return &command{
		description:   b.description,
		flagGroups:    flagGroups,
		arguments:     arguments,
		unmappedLabel: b.unmapLabel,
	}, nil
}
