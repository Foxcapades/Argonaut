package argo

import (
	"fmt"
	"os"
)

type CommandTreeBuilder interface {
	WithDescription(desc string) CommandTreeBuilder

	WithCallback(cb CommandTreeCallback) CommandTreeBuilder

	WithHelpDisabled() CommandTreeBuilder

	WithBranch(branch CommandBranchBuilder) CommandTreeBuilder

	WithLeaf(leaf CommandLeafBuilder) CommandTreeBuilder

	WithCommandGroup(group CommandGroupBuilder) CommandTreeBuilder

	WithFlag(flag FlagBuilder) CommandTreeBuilder

	WithFlagGroup(flagGroup FlagGroupBuilder) CommandTreeBuilder

	Parse(args []string) (CommandTree, error)

	MustParse(args []string) CommandTree
}

func NewCommandTreeBuilder() CommandTreeBuilder {
	return &commandTreeBuilder{
		commandGroups: []CommandGroupBuilder{NewCommandGroupBuilder(defaultGroupName)},
		flagGroups:    []FlagGroupBuilder{NewFlagGroupBuilder(defaultGroupName)},
	}
}

type commandTreeBuilder struct {
	desc          string
	helpDisabled  bool
	commandGroups []CommandGroupBuilder
	flagGroups    []FlagGroupBuilder
	callback      CommandTreeCallback
}

func (t *commandTreeBuilder) WithDescription(desc string) CommandTreeBuilder {
	t.desc = desc
	return t
}

func (t *commandTreeBuilder) WithCallback(cb CommandTreeCallback) CommandTreeBuilder {
	t.callback = cb
	return t
}

func (t *commandTreeBuilder) WithHelpDisabled() CommandTreeBuilder {
	t.helpDisabled = true
	return t
}

func (t *commandTreeBuilder) WithBranch(branch CommandBranchBuilder) CommandTreeBuilder {
	t.commandGroups[0].WithBranch(branch)
	return t
}

func (t *commandTreeBuilder) WithLeaf(leaf CommandLeafBuilder) CommandTreeBuilder {
	t.commandGroups[0].WithLeaf(leaf)
	return t
}

func (t *commandTreeBuilder) WithCommandGroup(group CommandGroupBuilder) CommandTreeBuilder {
	t.commandGroups = append(t.commandGroups, group)
	return t
}

func (t *commandTreeBuilder) WithFlag(flag FlagBuilder) CommandTreeBuilder {
	t.flagGroups[0].WithFlag(flag)
	return t
}

func (t *commandTreeBuilder) WithFlagGroup(flagGroup FlagGroupBuilder) CommandTreeBuilder {
	t.flagGroups = append(t.flagGroups, flagGroup)
	return t
}

func (t commandTreeBuilder) Parse(args []string) (CommandTree, error) {
	ct, err := t.build()
	if err != nil {
		return nil, err
	}

	err = CommandTreeInterpreter(args, ct).Run()
	if err != nil {
		return nil, err
	}

	return ct, nil
}

func (t commandTreeBuilder) MustParse(args []string) CommandTree {
	ct := mustReturn(t.build())
	must(CommandTreeInterpreter(args, ct).Run())
	return ct
}

func (t *commandTreeBuilder) build() (CommandTree, error) {
	errs := newMultiError()

	tree := new(commandTree)

	if !t.helpDisabled {
		group := t.flagGroups[0]

		if len(t.flagGroups) > 1 || t.flagGroups[0].size() > 5 {
			group = NewFlagGroupBuilder("Meta Flags")
			t.flagGroups = append(t.flagGroups, group)
		}

		useLongH := true
		useShortH := true

		for _, group := range t.flagGroups {
			for _, flag := range group.getFlags() {
				if flag.getShortForm() == 'h' {
					useShortH = false
				}
				if flag.getLongForm() == "help" {
					useLongH = false
				}
				if !(useShortH || useLongH) {
					break
				}
			}
		}

		if useShortH || useLongH {
			group.WithFlag(makeCommandTreeHelpFlag(useLongH, useLongH, tree))
		}
	}

	flagGroups := make([]FlagGroup, 0, len(t.flagGroups))
	uniqueFlagNames(t.flagGroups, errs)
	for _, builder := range t.flagGroups {
		if builder.hasFlags() {
			if group, err := builder.build(); err != nil {
				errs.AppendError(err)
			} else {
				flagGroups = append(flagGroups, group)
			}
		}
	}

	commandGroups := make([]CommandGroup, 0, len(t.commandGroups))
	massUniqueCommandNames(t.commandGroups, errs)
	for _, builder := range t.commandGroups {
		builder.Parent(tree)
		if builder.hasSubcommands() {
			if group, err := builder.build(); err != nil {
				errs.AppendError(err)
			} else {
				commandGroups = append(commandGroups, group)
			}
		}
	}

	if len(errs.Errors()) > 0 {
		return nil, errs
	}

	tree.description = t.desc
	tree.flagGroups = flagGroups
	tree.commandGroups = commandGroups
	tree.callback = t.callback

	return tree, nil
}

func makeCommandTreeHelpFlag(short, long bool, tree CommandTree) FlagBuilder {
	out := NewFlagBuilder().
		setIsHelpFlag().
		WithCallback(func(flag Flag) {
			fmt.Println(renderCommandTree(tree))
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
