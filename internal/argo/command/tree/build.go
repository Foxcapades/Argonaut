package tree

import (
	"errors"

	"github.com/foxcapades/argonaut/v3/internal/argo/command/common"
	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/internal/render"
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/internal/xerr"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func Build(builder argo.CommandTreeBuilder) (argo.CommandTree, error) {
	errs := xerr.NewMultiError()
	tree := new(commandTree)

	if !builder.HasSubcommands() {
		errs.AppendError(errors.New("command tree has no subcommands"))
	}

	var flagGroups []argo.FlagGroupBuilder
	if !builder.IsHelpDisabled() {
		flagGroups = common.ConfigureHelpFlags[argo.CommandTree](builder, render.CommandTreeHelpRenderer(), tree)
	} else {
		flagGroups = builder.FlagGroups(true)
	}

	tree.flagGroups = make([]argo.FlagGroup, 0, len(flagGroups))
	flag.UniqueFlagNames(flagGroups, errs)
	for _, groupBuilder := range flagGroups {
		if groupBuilder.HasFlags() {
			if group, err := flag.BuildGroup(groupBuilder); err != nil {
				errs.AppendError(err)
			} else {
				tree.flagGroups = append(tree.flagGroups, group)
			}
		}
	}

	commandGroups := make([]argo.CommandGroup, 0, len(builder.CommandGroups()))
	massUniqueCommandNames(builder.CommandGroups(), errs)
	for _, builder := range builder.commandGroups {
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
	tree.description = builder.desc
	tree.flagGroups = flagGroups
	tree.commandGroups = commandGroups
	tree.callback = builder.callback
	tree.onIncompleteHandler = utils.IfElse(builder.onIncompleteHandler == nil, defaultOnIncompleteHandler, builder.onIncompleteHandler)

	return tree, nil
}
