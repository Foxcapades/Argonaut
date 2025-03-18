package tree

import (
	"errors"
	"os"

	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/internal/chars"
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func NewBuilder() argo.CommandTreeBuilder {
	return &commandTreeBuilder{
		commandGroups: []argo.CommandGroupBuilder{NewCommandGroupBuilder(chars.DefaultGroupName)},
		flagGroups:    []argo.FlagGroupBuilder{flag.NewGroupBuilder(chars.DefaultGroupName)},
	}
}

type commandTreeBuilder struct {
	desc          string
	helpDisabled  bool
	commandGroups []argo.CommandGroupBuilder
	flagGroups    []argo.FlagGroupBuilder
	callback      argo.CommandTreeCallback

	onIncompleteHandler argo.IncompleteCommandHandler
}

func (t *commandTreeBuilder) WithDescription(desc string) argo.CommandTreeBuilder {
	t.desc = desc
	return t
}

func (t *commandTreeBuilder) WithCallback(cb argo.CommandTreeCallback) argo.CommandTreeBuilder {
	t.callback = cb
	return t
}

func (t *commandTreeBuilder) WithHelpDisabled() argo.CommandTreeBuilder {
	t.helpDisabled = true
	return t
}

func (t *commandTreeBuilder) WithBranch(branch argo.BranchBuilder) argo.CommandTreeBuilder {
	t.commandGroups[0].WithBranch(branch)
	return t
}

func (t *commandTreeBuilder) WithLeaf(leaf argo.LeafBuilder) argo.CommandTreeBuilder {
	t.commandGroups[0].WithLeaf(leaf)
	return t
}

func (t *commandTreeBuilder) WithCommandGroup(group argo.CommandGroupBuilder) argo.CommandTreeBuilder {
	t.commandGroups = append(t.commandGroups, group)
	return t
}

func (t *commandTreeBuilder) WithFlag(flag argo.FlagBuilder) argo.CommandTreeBuilder {
	t.flagGroups[0].WithFlag(flag)
	return t
}

func (t *commandTreeBuilder) WithFlagGroup(flagGroup argo.FlagGroupBuilder) argo.CommandTreeBuilder {
	t.flagGroups = append(t.flagGroups, flagGroup)
	return t
}

func (t *commandTreeBuilder) OnIncomplete(handler argo.IncompleteCommandHandler) argo.CommandTreeBuilder {
	t.onIncompleteHandler = handler
	return t
}

func (t *commandTreeBuilder) Parse(args []string) (argo.CommandTree, error) {
	ctx := new(argo.WarningContext)
	ct, err := t.Build(ctx)
	if err != nil {
		return nil, err
	}

	err = newCommandTreeInterpreter(args, ct).Run()
	if err != nil {
		return nil, err
	}

	return ct, nil
}

func (t *commandTreeBuilder) MustParse(args []string) argo.CommandTree {
	ctx := new(argo.WarningContext)
	ct := utils.MustReturn(t.Build(ctx))
	utils.Must(newCommandTreeInterpreter(args, ct).Run())
	return ct
}

func (t *commandTreeBuilder) hasSubCommands() bool {
	for _, group := range t.commandGroups {
		if group.hasSubcommands() {
			return true
		}
	}

	return false
}

func (t *commandTreeBuilder) Build(warnings *argo.WarningContext) (argo.CommandTree, error) {
	errs := argo.NewMultiError()

	tree := new(commandTree)

	if !t.hasSubCommands() {
		errs.AppendError(errors.New("command tree has no subcommands"))
	}

	if !t.helpDisabled {
		group := t.flagGroups[0]

		if len(t.flagGroups) > 1 || t.flagGroups[0].size() > 5 {
			group = cli_flag.NewFlagGroupBuilder("Help Flags")
			t.flagGroups = append(t.flagGroups, group)
		}

		useLongH := true
		useShortH := true

		for _, group := range t.flagGroups {
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
			group.WithFlag(makeCommandTreeHelpFlag(useLongH, useLongH, tree))
		}
	}

	flagGroups := make([]argo.FlagGroup, 0, len(t.flagGroups))
	uniqueFlagNames(t.flagGroups, errs)
	for _, builder := range t.flagGroups {
		if builder.HasFlags() {
			if group, err := builder.Build(warnings); err != nil {
				errs.AppendError(err)
			} else {
				flagGroups = append(flagGroups, group)
			}
		}
	}

	commandGroups := make([]argo.CommandGroup, 0, len(t.commandGroups))
	massUniqueCommandNames(t.commandGroups, errs)
	for _, builder := range t.commandGroups {
		builder.parent(tree)
		if builder.hasSubcommands() {
			if group, err := builder.Build(warnings); err != nil {
				errs.AppendError(err)
			} else {
				commandGroups = append(commandGroups, group)
			}
		}
	}

	if len(errs.Errors()) > 0 {
		return nil, errs
	}

	tree.warnings = warnings
	tree.description = t.desc
	tree.flagGroups = flagGroups
	tree.commandGroups = commandGroups
	tree.callback = t.callback
	tree.onIncompleteHandler = utils.IfElse(t.onIncompleteHandler == nil, defaultOnIncompleteHandler, t.onIncompleteHandler)

	return tree, nil
}

func makeCommandTreeHelpFlag(short, long bool, tree argo.CommandTree) argo.FlagBuilder {
	out := flag.NewBuilder().
		MarkAsHelpFlag().
		WithCallback(func(flag argo.Flag) {
			utils.Must(comTreeRenderer{}.RenderHelp(tree, os.Stdout))
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
