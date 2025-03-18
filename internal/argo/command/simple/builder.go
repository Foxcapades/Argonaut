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
	errs := xerr.NewMultiError()
	com := new(command)

	com.warnings = ctx

	if !b.disableHelp {
		group := b.flagGroups[0]

		if len(b.flagGroups) > 1 || b.flagGroups[0].Size() > 5 {
			group = flag.NewGroupBuilder("Help Flags")
			b.flagGroups = append(b.flagGroups, group)
		}

		useLongH := true
		useShortH := true

		for _, group := range b.flagGroups {
			for _, flag := range group.Flags() {
				if flag.ShortForm() == 'h' {
					useShortH = false
				}
				if flag.LongForm() == "help" {
					useLongH = false
				}
				if !(useShortH || useLongH) {
					break
				}
			}
		}

		if useShortH || useLongH {
			group.WithFlag(makeCommandHelpFlag(useShortH, useLongH, com))
		}

	}

	com.flagGroups = make([]argo.FlagGroup, 0, len(b.flagGroups))
	uniqueFlagNames(b.flagGroups, errs)
	for _, builder := range b.flagGroups {
		if builder.HasFlags() {
			if group, err := builder.Build(ctx); err != nil {
				errs.AppendError(err)
			} else {
				com.flagGroups = append(com.flagGroups, group)
			}
		}
	}

	forceRequiredUntil := 0
	for i, builder := range b.arguments {
		if builder.IsRequired() {
			forceRequiredUntil = i
		}
	}

	com.arguments = make([]argo.Argument, 0, len(b.arguments))
	for i, builder := range b.arguments {

		if i < forceRequiredUntil && !builder.IsRequired() {
			builder.Require()
			ctx.AppendWarning(fmt.Sprintf("argument %d was not marked as required, but preceded required argument %d", i+1, forceRequiredUntil+1))
		}

		if arg, err := builder.Build(ctx); err != nil {
			errs.AppendError(err)
		} else {
			com.arguments = append(com.arguments, arg)
		}
	}

	if len(errs.Errors()) > 0 {
		return nil, errs
	}

	com.description = b.description
	com.unmappedLabel = b.unmapLabel
	com.callback = b.callback

	return com, nil
}

func makeCommandHelpFlag(short, long bool, com argo.Command) argo.FlagBuilder {
	out := flag.NewBuilder().
		MarkAsHelpFlag().
		WithCallback(func(flag argo.Flag) {
			utils.Must(comRenderer{}.RenderHelp(com, os.Stdout))
			os.Exit(0)
		}).
		WithDescription("Prints this help text.")

	if short {
		out.WithShortForm('h')
	}

	if long {
		out.WithLongForm("help")
	}

	return out
}
