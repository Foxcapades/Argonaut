package tree

import (
	"errors"
	"fmt"
	"os"

	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/internal/chars"
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/internal/xerr"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func NewLeafBuilder(name string) argo.LeafCommandBuilder {
	out := new(leafBuilder)
	out.childBuilder = newChildBuilder[argo.LeafCommandBuilder, argo.LeafCommand](name, out)
	return out
}

type leafBuilder struct {
	childBuilder[argo.LeafCommandBuilder, argo.LeafCommand]

	umapLabel string
	arguments []argo.ArgumentBuilder
}

//

func (l *leafBuilder) WithArgument(argument argo.ArgumentBuilder) argo.LeafCommandBuilder {
	l.arguments = append(l.arguments, argument)
	return l
}

func (l *leafBuilder) WithArguments(arguments ...argo.ArgumentBuilder) argo.LeafCommandBuilder {
	l.arguments = append(l.arguments, arguments...)
	return l
}

func (l *leafBuilder) HasArguments() bool {
	return len(l.arguments) > 0
}

func (l *leafBuilder) Arguments() []argo.ArgumentBuilder {
	return l.arguments
}

//

func (l *leafBuilder) WithUnmappedInputLabel(label string) argo.LeafCommandBuilder {
	l.umapLabel = label
	return l
}

func (l *leafBuilder) HasUnmappedInputLabel() bool {
	return len(l.umapLabel) > 0
}

func (l *leafBuilder) UnmappedInputLabel() string {
	return l.umapLabel
}

//

func (l *leafBuilder) Build(ctx *argo.WarningContext) (argo.LeafCommand, error) {
	errs := xerr.NewMultiError()

	// Ensure the group name is not blank
	if err := chars.ValidateCommandNodeName(l.name); err != nil {
		errs.AppendError(err)
	}

	// Ensure the aliases are all not blank
	for _, alias := range l.aliases {
		if chars.IsBlank(alias) {
			errs.AppendError(errors.New("command leaf aliases must not be blank"))
		}
	}

	if l.parent == nil {
		panic("illegal state: attempted to build a command leaf with no parent set")
	}

	leaf := new(commandLeaf)
	leaf.warnings = ctx

	forceRequiredUntil := 0
	for i, builder := range l.arguments {
		if builder.IsRequired() {
			forceRequiredUntil = i
		}
	}

	leaf.args = make([]argo.Argument, 0, len(l.arguments))
	for i, builder := range l.arguments {
		if i < forceRequiredUntil && !builder.IsRequired() {
			builder.Require()
			ctx.AppendWarning(fmt.Sprintf("argument %d was not marked as required, but preceded required argument %d", i+1, forceRequiredUntil+1))
		}
		if arg, err := builder.Build(ctx); err != nil {
			errs.AppendError(err)
		} else {
			leaf.args = append(leaf.args, arg)
		}
	}

	uniqueFlagNames(l.flagGroups, errs)
	leaf.flagGroups = make([]argo.FlagGroup, 0, len(l.flagGroups))
	for _, builder := range l.flagGroups {
		if builder.HasFlags() {
			if fg, err := builder.Build(ctx); err != nil {
				errs.AppendError(err)
			} else {
				leaf.flagGroups = append(leaf.flagGroups, fg)
			}
		}
	}

	if !l.disableHelp {
		useShortH := true
		useLongH := true

		for _, group := range leaf.flagGroups {
			for _, flag := range group.Flags() {
				if flag.ShortForm() == 'h' {
					useShortH = false
				}
				if flag.LongForm() == "help" {
					useLongH = false
				}
			}
		}

		if useShortH || useLongH {
			if len(leaf.flagGroups) == 0 || leaf.flagGroups[0].Name() != chars.DefaultGroupName || leaf.flagGroups[0].Size() > 5 {
				group, err := flag.NewGroupBuilder("Help Flags").
					WithFlag(makeLeafHelp(useShortH, useLongH, leaf)).
					Build(ctx)

				if err != nil {
					errs.AppendError(err)
				} else {
					leaf.flagGroups = append(leaf.flagGroups, group)
				}
			} else {
				flag, err := makeLeafHelp(useShortH, useLongH, leaf).Build(ctx)

				if err != nil {
					errs.AppendError(err)
				} else {
					group := leaf.flagGroups[0].(*flagGroup) // FIXME: don't assume internal access!!!
					group.flags = append(group.flags, flag)
				}
			}
		}
	}

	if len(errs.Errors()) > 0 {
		return nil, errs
	}

	leaf.name = l.name
	leaf.description = l.description
	leaf.aliases = l.aliases
	leaf.parent = l.parentNode
	leaf.callback = l.callback
	leaf.unmappedLabel = l.umapLabel

	return leaf, nil
}

func makeLeafHelp(short, long bool, leaf argo.LeafCommand) argo.FlagBuilder {
	builder := flag.NewBuilder().
		MarkAsHelpFlag().
		WithCallback(func(flag argo.Flag) {
			utils.Must(comLeafRenderer{}.RenderHelp(leaf, os.Stdout))
			os.Exit(0)
		}).
		WithDescription("Prints this help text.")

	if short {
		builder.WithShortForm('h')
	}

	if long {
		builder.WithLongForm("help")
	}

	return builder
}
