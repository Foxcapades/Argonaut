package command

import (
	"github.com/Foxcapades/Argonaut/internal/consts"
	"github.com/Foxcapades/Argonaut/internal/impl/command/comutil"
	"github.com/Foxcapades/Argonaut/internal/impl/flag"
	"github.com/Foxcapades/Argonaut/internal/impl/flag/flagutil"
	"github.com/Foxcapades/Argonaut/internal/interpret"
	"github.com/Foxcapades/Argonaut/internal/util"
	"github.com/Foxcapades/Argonaut/internal/xerr"
	"github.com/Foxcapades/Argonaut/pkg/argo"
)

func TreeBuilder() argo.CommandTreeBuilder {
	return &treeBuilder{
		commandGroups: []argo.CommandGroupBuilder{GroupBuilder(consts.DefaultGroupName)},
		flagGroups:    []argo.FlagGroupBuilder{flag.GroupBuilder(consts.DefaultGroupName)},
	}
}

type treeBuilder struct {
	desc          string
	help          bool
	commandGroups []argo.CommandGroupBuilder
	flagGroups    []argo.FlagGroupBuilder
}

func (t *treeBuilder) WithDescription(desc string) argo.CommandTreeBuilder {
	t.desc = desc
	return t
}

func (t *treeBuilder) WithHelpDisabled() argo.CommandTreeBuilder {
	t.help = true
	return t
}

func (t *treeBuilder) WithBranch(branch argo.CommandBranchBuilder) argo.CommandTreeBuilder {
	t.commandGroups[0].WithBranch(branch)
	return t
}

func (t *treeBuilder) WithLeaf(leaf argo.CommandLeafBuilder) argo.CommandTreeBuilder {
	t.commandGroups[0].WithLeaf(leaf)
	return t
}

func (t *treeBuilder) WithCommandGroup(group argo.CommandGroupBuilder) argo.CommandTreeBuilder {
	t.commandGroups = append(t.commandGroups, group)
	return t
}

func (t *treeBuilder) WithFlag(flag argo.FlagBuilder) argo.CommandTreeBuilder {
	t.flagGroups[0].WithFlag(flag)
	return t
}

func (t *treeBuilder) WithFlagGroup(flagGroup argo.FlagGroupBuilder) argo.CommandTreeBuilder {
	t.flagGroups = append(t.flagGroups, flagGroup)
	return t
}

func (t treeBuilder) Parse(args []string) (argo.CommandTree, error) {
	ct, err := t.Build()
	if err != nil {
		return nil, err
	}

	err = interpret.CommandTreeInterpreter(args, ct).Run()
	if err != nil {
		return nil, err
	}

	return ct, nil
}

func (t treeBuilder) MustParse(args []string) argo.CommandTree {
	ct := util.MustReturn(t.Build())
	util.Must(interpret.CommandTreeInterpreter(args, ct).Run())
	return ct
}

func (t treeBuilder) Build() (argo.CommandTree, error) {
	errs := xerr.NewMultiError()

	flagGroups := make([]argo.FlagGroup, 0, len(t.flagGroups))
	flagutil.UniqueFlagNames(t.flagGroups, errs)
	for _, builder := range t.flagGroups {
		if builder.HasFlags() {
			if group, err := builder.Build(); err != nil {
				errs.AppendError(err)
			} else {
				flagGroups = append(flagGroups, group)
			}
		}
	}

	commandGroups := make([]argo.CommandGroup, 0, len(t.commandGroups))
	comutil.MassUniqueNames(t.commandGroups, errs)
	for _, builder := range t.commandGroups {
		if builder.HasSubcommands() {
			if group, err := builder.Build(); err != nil {
				errs.AppendError(err)
			} else {
				commandGroups = append(commandGroups, group)
			}
		}
	}

	if len(errs.Errors()) > 0 {
		return nil, errs
	}

	return &tree{
		description:   t.desc,
		flagGroups:    flagGroups,
		commandGroups: commandGroups,
	}, nil
}
