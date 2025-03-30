package tree

import (
	"errors"

	"github.com/foxcapades/argonaut/v3/internal/argo/command/common"
	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/internal/argo/opts"
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/internal/xerr"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func Build(builder argo.TreeCommandBuilder, options argo.Options) (argo.TreeCommand, error) {
	opts.FixOptions(&options)

	errs := xerr.NewMultiError()
	tree := new(Tree)

	if !builder.HasSubcommands() {
		errs.AppendError(errors.New("command tree has no subcommands"))
	}

	common.TryAddHelpFlags(builder, MakeRenderTreeHelpCallback(tree, options), options)

	tree.flagGroups = flag.BuildGroups(builder.FlagGroups(true), options, errs)
	tree.commandGroups = BuildCommandGroups(builder.CommandGroups(true), options, tree, errs)

	commandGroups := make([]argo.CommandGroup, 0, len(builder.CommandGroups(true)))
	UniqueCommandNames(builder.CommandGroups(true), errs)
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
		MakeDefaultOnIncompleteHandler[argo.TreeCommand](options),
	)

	return tree, nil
}
