package command

import (
	"errors"
	"fmt"

	"github.com/Foxcapades/Argonaut/internal/chars"
	"github.com/Foxcapades/Argonaut/internal/impl/flag"
	"github.com/Foxcapades/Argonaut/internal/impl/flag/flagutil"
	"github.com/Foxcapades/Argonaut/internal/xerr"
	"github.com/Foxcapades/Argonaut/pkg/argo"
)

func NewLeafBuilder(name string) argo.CommandLeafBuilder {
	return &leafBuilder{
		name:       name,
		flagGroups: []argo.FlagGroupBuilder{flag.GroupBuilder("default")},
	}
}

type leafBuilder struct {
	parent      argo.CommandNode
	name        string
	description string
	aliases     []string
	arguments   []argo.ArgumentBuilder
	flagGroups  []argo.FlagGroupBuilder
}

func (l *leafBuilder) GetName() string { return l.name }

func (l *leafBuilder) GetAliases() []string { return l.aliases }

func (l *leafBuilder) Parent(node argo.CommandNode) { l.parent = node }

func (l *leafBuilder) WithDescription(desc string) argo.CommandLeafBuilder {
	l.description = desc
	return l
}

func (l *leafBuilder) WithAliases(aliases ...string) argo.CommandLeafBuilder {
	l.aliases = append(l.aliases, aliases...)
	return l
}

func (l *leafBuilder) WithArgument(argument argo.ArgumentBuilder) argo.CommandLeafBuilder {
	l.arguments = append(l.arguments, argument)
	return l
}

func (l *leafBuilder) WithFlagGroup(flagGroup argo.FlagGroupBuilder) argo.CommandLeafBuilder {
	l.flagGroups = append(l.flagGroups, flagGroup)
	return l
}

func (l *leafBuilder) WithFlag(flag argo.FlagBuilder) argo.CommandLeafBuilder {
	l.flagGroups[0].WithFlag(flag)
	return l
}

func (l *leafBuilder) Build() (argo.CommandLeaf, error) {
	errs := xerr.NewMultiError()

	// Ensure the group name is not blank
	if chars.IsBlank(l.name) {
		errs.AppendError(errors.New("command leaf names must not be blank"))
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

	arguments := make([]argo.Argument, 0, len(l.arguments))
	for _, builder := range l.arguments {
		if arg, err := builder.Build(); err != nil {
			errs.AppendError(err)
		} else {
			arguments = append(arguments, arg)
		}
	}

	flagutil.UniqueFlagNames(l.flagGroups, errs)
	flagGroups := make([]argo.FlagGroup, 0, len(l.flagGroups))
	for _, builder := range l.flagGroups {
		if !builder.HasFlags() {
			errs.AppendError(fmt.Errorf("flag group %s has no flags", builder.GetName()))
			continue
		}

		if fg, err := builder.Build(); err != nil {
			errs.AppendError(err)
		} else {
			flagGroups = append(flagGroups, fg)
		}
	}

	if len(errs.Errors()) > 0 {
		return nil, errs
	}

	return &commandLeaf{
		name:    l.name,
		desc:    l.description,
		aliases: l.aliases,
		parent:  l.parent,
		flags:   flagGroups,
		args:    arguments,
	}, nil
}
