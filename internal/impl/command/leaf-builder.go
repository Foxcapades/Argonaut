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
	umapLabel   string
	aliases     []string
	arguments   []argo.ArgumentBuilder
	flagGroups  []argo.FlagGroupBuilder
	callback    argo.CommandLeafCallback
}

func (l *leafBuilder) GetName() string {
	return l.name
}

func (l *leafBuilder) Parent(node argo.CommandNode) {
	l.parent = node
}

// Aliases /////////////////////////////////////////////////////////////////////

func (l *leafBuilder) WithAliases(aliases ...string) argo.CommandLeafBuilder {
	l.aliases = append(l.aliases, aliases...)
	return l
}

func (l leafBuilder) GetAliases() []string {
	return l.aliases
}

func (l leafBuilder) HasAliases() bool {
	return len(l.aliases) > 0
}

// Description /////////////////////////////////////////////////////////////////

func (l *leafBuilder) WithDescription(desc string) argo.CommandLeafBuilder {
	l.description = desc
	return l
}

func (l leafBuilder) HasDescription() bool {
	return len(l.description) > 0
}

func (l leafBuilder) GetDescription() string {
	return l.description
}

// Arguments ///////////////////////////////////////////////////////////////////

func (l *leafBuilder) WithArgument(argument argo.ArgumentBuilder) argo.CommandLeafBuilder {
	l.arguments = append(l.arguments, argument)
	return l
}

func (l leafBuilder) HasArguments() bool {
	return len(l.arguments) > 0
}

func (l leafBuilder) GetArguments() []argo.ArgumentBuilder {
	return l.arguments
}

// Flag Groups /////////////////////////////////////////////////////////////////

func (l *leafBuilder) WithFlagGroup(flagGroup argo.FlagGroupBuilder) argo.CommandLeafBuilder {
	l.flagGroups = append(l.flagGroups, flagGroup)
	return l
}

func (l leafBuilder) HasCustomFlagGroups() bool {
	return len(l.flagGroups) > 1
}

func (l leafBuilder) GetFlagGroups() []argo.FlagGroupBuilder {
	return l.flagGroups
}

func (l leafBuilder) GetCustomFlagGroups() []argo.FlagGroupBuilder {
	return l.flagGroups[1:]
}

// Flags ///////////////////////////////////////////////////////////////////////

func (l *leafBuilder) WithFlag(flag argo.FlagBuilder) argo.CommandLeafBuilder {
	l.flagGroups[0].WithFlag(flag)
	return l
}

func (l leafBuilder) HasFlags() bool {
	for _, g := range l.flagGroups {
		if g.HasFlags() {
			return true
		}
	}

	return false
}

// Unmapped ////////////////////////////////////////////////////////////////////

func (l *leafBuilder) WithUnmappedLabel(label string) argo.CommandLeafBuilder {
	l.umapLabel = label
	return l
}

func (l leafBuilder) HasUnmappedLabel() bool {
	return len(l.umapLabel) > 0
}

func (l leafBuilder) GetUnmappedLabel() string {
	return l.umapLabel
}

// Callback ////////////////////////////////////////////////////////////////////

func (l *leafBuilder) WithCallback(cb argo.CommandLeafCallback) argo.CommandLeafBuilder {
	l.callback = cb
	return l
}

func (l leafBuilder) HasCallback() bool {
	return l.callback != nil
}

func (l leafBuilder) GetCallback() argo.CommandLeafCallback {
	return l.callback
}

// Build ///////////////////////////////////////////////////////////////////////

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
		name:     l.name,
		desc:     l.description,
		aliases:  l.aliases,
		parent:   l.parent,
		flags:    flagGroups,
		args:     arguments,
		callback: l.callback,
	}, nil
}
