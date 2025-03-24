package tree

import (
	"errors"

	"github.com/foxcapades/argonaut/v3/internal/argo/command/common"
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/internal/xerr"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func Build(builder argo.TreeCommandBuilder, options argo.Options) (argo.TreeCommand, error) {
	errs := xerr.NewMultiError()
	tree := new(Tree)

	if !builder.HasSubcommands() {
		errs.AppendError(errors.New("command tree has no subcommands"))
	}

	var flagGroups []argo.FlagGroupBuilder
	if !builder.IsHelpDisabled() {
		common.TryAddHelpFlags(builder, MakeRenderTreeHelpCallback(tree, options), options)
		flagGroups = common.ConfigureHelpFlags[argo.TreeCommand](builder, RenderHelp, tree)
	} else {
		flagGroups = builder.FlagGroups(true)
	}

	tree.flagGroups = processFlagGroups(flagGroups, errs)
	tree.commandGroups = processCommandGroups(builder.CommandGroups(true), options, tree, errs)

	commandGroups := make([]argo.CommandGroup, 0, len(builder.CommandGroups(true)))
	massUniqueCommandNames(builder.CommandGroups(true), errs)
	for _, builder := range builder.CommandGroups(true) {
		if builder.HasSubcommands() {
			if group, err := BuildGroup(builder, options, tree); err != nil {
				errs.AppendError(err)
			} else {
				commandGroups = append(commandGroups, group)
			}
		}
	}

	if len(errs.Errors()) > 0 {
		return nil, errs
	}

	tree.description = builder.Description()
	tree.commandGroups = commandGroups
	tree.callback = builder.Callback()
	tree.incompleteFn = utils.IfElse(
		builder.HasIncompleteHandler(),
		builder.IncompleteHandler(),
		makeDefaultOnIncompleteHandler[argo.TreeCommand](options),
	)

	return tree, nil
}
