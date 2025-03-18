package tree

import (
	"errors"
	"fmt"
	"os"

	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/internal/chars"
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func NewLeafBuilder(name string) argo.LeafBuilder {
	return &leafBuilder{
		name:       name,
		flagGroups: []argo.FlagGroupBuilder{nil},
	}
}

type leafBuilder struct {
	parentNode  argo.Node
	disableHelp bool
	name        string
	description string
	umapLabel   string
	aliases     []string
	arguments   []argo.ArgumentBuilder
	flagGroups  []argo.FlagGroupBuilder
	callback    argo.CommandLeafCallback
}

// PUBLIC API //////////////////////////////////////////////////////////////////////////////////////////////////////////

//

func (l *leafBuilder) Name() string {
	return l.name
}

//

func (l *leafBuilder) WithDescription(desc string) argo.LeafBuilder {
	l.description = desc
	return l
}

func (l *leafBuilder) HasDescription() bool {
	return len(l.description) > 0
}

func (l *leafBuilder) Description() string {
	return l.description
}

//

func (l *leafBuilder) WithAlias(alias string) argo.LeafBuilder {
	l.aliases = append(l.aliases, alias)
	return l
}

func (l *leafBuilder) WithAliases(aliases ...string) argo.LeafBuilder {
	l.aliases = append(l.aliases, aliases...)
	return l
}

func (l *leafBuilder) HasAliases() bool {
	return len(l.aliases) > 0
}

func (l *leafBuilder) Aliases() []string {
	return l.aliases
}

//

func (l *leafBuilder) WithFlagGroup(flagGroup argo.FlagGroupBuilder) argo.LeafBuilder {
	l.flagGroups = append(l.flagGroups, flagGroup)
	return l
}

func (l *leafBuilder) WithFlagGroups(flagGroups ...argo.FlagGroupBuilder) argo.LeafBuilder {
	l.flagGroups = append(l.flagGroups, flagGroups...)
	return l
}

func (l *leafBuilder) HasFlagGroups(includeDefault bool) bool {
	var min int

	if includeDefault && l.hasDefaultFlagGroup() {
		min = 0
	} else {
		min = 1
	}

	return len(l.flagGroups) > min
}

func (l *leafBuilder) FlagGroups(includeDefault bool) []argo.FlagGroupBuilder {
	if includeDefault && l.hasDefaultFlagGroup() {
		return append([]argo.FlagGroupBuilder(nil), l.flagGroups...)
	}

	return append([]argo.FlagGroupBuilder(nil), l.flagGroups[1:]...)
}

func (l *leafBuilder) hasDefaultFlagGroup() bool {
	return l.flagGroups[0] != nil
}

func (l *leafBuilder) WithFlag(flag argo.FlagBuilder) argo.LeafBuilder {
	l.flagGroups[0].WithFlag(flag)
	return l
}

func (l *leafBuilder) WithFlags(flags ...argo.FlagBuilder) argo.LeafBuilder {
	for _, flag := range flags {
		l.flagGroups[0].WithFlag(flag)
	}
	return l
}

func (l *leafBuilder) HasFlags() bool {
	for _, group := range l.flagGroups {
		if group != nil && group.Size() > 0 {
			return true
		}
	}

	return false
}

//

func (l *leafBuilder) WithArgument(argument argo.ArgumentBuilder) argo.LeafBuilder {
	l.arguments = append(l.arguments, argument)
	return l
}

func (l *leafBuilder) WithUnmappedLabel(label string) argo.LeafBuilder {
	l.umapLabel = label
	return l
}

func (l *leafBuilder) WithCallback(cb argo.CommandLeafCallback) argo.LeafBuilder {
	l.callback = cb
	return l
}

func (l *leafBuilder) WithHelpDisabled() argo.LeafBuilder {
	l.disableHelp = true
	return l
}

// INTERNALS ///////////////////////////////////////////////////////////////////////////////////////////////////////////

func (l *leafBuilder) getName() string {
	return l.name
}

func (l *leafBuilder) parent(node argo.Node) {
	l.parentNode = node
}

func (l *leafBuilder) getAliases() []string {
	return l.aliases
}

func (l *leafBuilder) Build(ctx *argo.WarningContext) (argo.CommandLeaf, error) {
	errs := argo.NewMultiError()

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

	if l.parentNode == nil {
		panic("illegal state: attempted to build a command leaf with no parent set")
	}

	leaf := new(commandLeaf)
	leaf.warnings = ctx

	forceRequiredUntil := 0
	for i, builder := range l.arguments {
		if builder.isRequired() {
			forceRequiredUntil = i
		}
	}

	leaf.args = make([]argo.Argument, 0, len(l.arguments))
	for i, builder := range l.arguments {
		if i < forceRequiredUntil && !builder.isRequired() {
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
	leaf.flags = make([]argo.FlagGroup, 0, len(l.flagGroups))
	for _, builder := range l.flagGroups {
		if builder.HasFlags() {
			if fg, err := builder.Build(ctx); err != nil {
				errs.AppendError(err)
			} else {
				leaf.flags = append(leaf.flags, fg)
			}
		}
	}

	if !l.disableHelp {
		useShortH := true
		useLongH := true

		for _, group := range leaf.flags {
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
			if len(leaf.flags) == 0 || leaf.flags[0].Name() != chars.DefaultGroupName || leaf.flags[0].size() > 5 {
				group, err := flag.NewGroupBuilder("Help Flags").
					WithFlag(makeLeafHelp(useShortH, useLongH, leaf)).
					Build(ctx)

				if err != nil {
					errs.AppendError(err)
				} else {
					leaf.flags = append(leaf.flags, group)
				}
			} else {
				flag, err := makeLeafHelp(useShortH, useLongH, leaf).Build(ctx)

				if err != nil {
					errs.AppendError(err)
				} else {
					group := leaf.flags[0].(*flagGroup) // FIXME: don't assume internal access!!!
					group.flags = append(group.flags, flag)
				}
			}
		}
	}

	if len(errs.Errors()) > 0 {
		return nil, errs
	}

	leaf.name = l.name
	leaf.desc = l.description
	leaf.aliases = l.aliases
	leaf.parent = l.parentNode
	leaf.callback = l.callback
	leaf.uLabel = l.umapLabel

	return leaf, nil
}

func makeLeafHelp(short, long bool, leaf argo.CommandLeaf) argo.FlagBuilder {
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
