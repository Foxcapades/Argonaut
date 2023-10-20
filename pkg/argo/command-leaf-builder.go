package argo

import (
	"errors"
	"os"
)

// CommandLeafBuilder defines a builder type that is used to construct
// CommandLeaf instances.
type CommandLeafBuilder interface {
	getName() string

	// WithAliases attaches the given aliases to this CommandLeafBuilder.
	//
	// Command aliases must be unique per level in a command tree.  This means
	// that for a given step in the tree, no alias may conflict with another
	// CommandNode's name or aliases.
	//
	// If a conflict is found, an error will be returned when attempting to build
	// the CommandLeaf.
	WithAliases(aliases ...string) CommandLeafBuilder

	// WithDescription sets a description for the CommandLeaf to be built.
	//
	// Descriptions are used when rendering help text.
	WithDescription(desc string) CommandLeafBuilder

	WithHelpDisabled() CommandLeafBuilder

	// WithUnmappedLabel provides a label for unmapped inputs.
	//
	// The unmapped label value is used when rendering the command usage line of
	// the auto-generated help text.  If a command expects an unknown number of
	// positional argument values, it is best to capture them as unmapped inputs
	// with a label.
	//
	// Example configuration:
	//     cli.CommandLeaf("my-leaf").
	//         WithUnmappedLabel("ITEMS...")
	//
	// Example usage line:
	//     Usage:
	//       my-leaf [ITEMS...]
	WithUnmappedLabel(label string) CommandLeafBuilder

	// WithCallback sets a callback on the CommandLeaf to be built that will be
	// executed when the command leaf is used, after any tree or branch callbacks
	// that were set on parent nodes.
	WithCallback(cb CommandLeafCallback) CommandLeafBuilder

	// WithArgument adds a positional argument to the CommandLeaf being built.
	WithArgument(argument ArgumentBuilder) CommandLeafBuilder

	// WithFlagGroup adds a new FlagGroup to this CommandLeaf being built.
	WithFlagGroup(flagGroup FlagGroupBuilder) CommandLeafBuilder

	// WithFlag adds the given flag to the default FlagGroup for the CommandLeaf
	// being built.
	WithFlag(flag FlagBuilder) CommandLeafBuilder

	getAliases() []string
	parent(node CommandNode)

	build() (CommandLeaf, error)
}

func NewCommandLeafBuilder(name string) CommandLeafBuilder {
	return &commandLeafBuilder{
		name:       name,
		flagGroups: []FlagGroupBuilder{NewFlagGroupBuilder(defaultGroupName)},
	}
}

type commandLeafBuilder struct {
	parentNode  CommandNode
	disableHelp bool
	name        string
	description string
	umapLabel   string
	aliases     []string
	arguments   []ArgumentBuilder
	flagGroups  []FlagGroupBuilder
	callback    CommandLeafCallback
}

// PUBLIC API //////////////////////////////////////////////////////////////////////////////////////////////////////////

func (l *commandLeafBuilder) WithAliases(aliases ...string) CommandLeafBuilder {
	l.aliases = append(l.aliases, aliases...)
	return l
}

func (l *commandLeafBuilder) WithDescription(desc string) CommandLeafBuilder {
	l.description = desc
	return l
}

func (l *commandLeafBuilder) WithArgument(argument ArgumentBuilder) CommandLeafBuilder {
	l.arguments = append(l.arguments, argument)
	return l
}

func (l *commandLeafBuilder) WithFlagGroup(flagGroup FlagGroupBuilder) CommandLeafBuilder {
	l.flagGroups = append(l.flagGroups, flagGroup)
	return l
}

func (l *commandLeafBuilder) WithFlag(flag FlagBuilder) CommandLeafBuilder {
	l.flagGroups[0].WithFlag(flag)
	return l
}

func (l *commandLeafBuilder) WithUnmappedLabel(label string) CommandLeafBuilder {
	l.umapLabel = label
	return l
}

func (l *commandLeafBuilder) WithCallback(cb CommandLeafCallback) CommandLeafBuilder {
	l.callback = cb
	return l
}

func (l *commandLeafBuilder) WithHelpDisabled() CommandLeafBuilder {
	l.disableHelp = true
	return l
}

// INTERNALS ///////////////////////////////////////////////////////////////////////////////////////////////////////////

func (l *commandLeafBuilder) getName() string {
	return l.name
}

func (l *commandLeafBuilder) hasFlags() bool {
	for _, group := range l.flagGroups {
		if group.hasFlags() {
			return true
		}
	}

	return false
}

func (l *commandLeafBuilder) hasCustomFlagGroups() bool {
	return len(l.flagGroups) > 1
}

func (l *commandLeafBuilder) parent(node CommandNode) {
	l.parentNode = node
}

func (l commandLeafBuilder) getAliases() []string {
	return l.aliases
}

func (l *commandLeafBuilder) build() (CommandLeaf, error) {
	errs := newMultiError()

	// Ensure the group name is not blank
	if err := validateCommandNodeName(l.name); err != nil {
		errs.AppendError(err)
	}

	// Ensure the aliases are all not blank
	for _, alias := range l.aliases {
		if isBlank(alias) {
			errs.AppendError(errors.New("command leaf aliases must not be blank"))
		}
	}

	if l.parentNode == nil {
		panic("illegal state: attempted to build a command leaf with no parent set")
	}

	leaf := new(commandLeaf)

	leaf.args = make([]Argument, 0, len(l.arguments))
	for _, builder := range l.arguments {
		if arg, err := builder.Build(); err != nil {
			errs.AppendError(err)
		} else {
			leaf.args = append(leaf.args, arg)
		}
	}

	uniqueFlagNames(l.flagGroups, errs)
	leaf.flags = make([]FlagGroup, 0, len(l.flagGroups))
	for _, builder := range l.flagGroups {
		if builder.hasFlags() {
			if fg, err := builder.build(); err != nil {
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
			if len(leaf.flags) == 0 || leaf.flags[0].Name() != defaultGroupName || leaf.flags[0].size() > 5 {
				group, err := NewFlagGroupBuilder("Meta Flags").
					WithFlag(makeLeafHelp(useShortH, useLongH, leaf)).
					build()

				if err != nil {
					errs.AppendError(err)
				} else {
					leaf.flags = append(leaf.flags, group)
				}
			} else {
				flag, err := makeLeafHelp(useShortH, useLongH, leaf).build()

				if err != nil {
					errs.AppendError(err)
				} else {
					group := leaf.flags[0].(*flagGroup)
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

func makeLeafHelp(short, long bool, leaf CommandLeaf) FlagBuilder {
	builder := NewFlagBuilder().
		setIsHelpFlag().
		WithCallback(func(flag Flag) {
			must(comLeafRenderer{}.RenderHelp(leaf, os.Stdout))
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
