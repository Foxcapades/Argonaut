package tree

import (
	"errors"

	"github.com/foxcapades/argonaut/v3/internal/argo/command/common"
	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/internal/xerr"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func Build(builder argo.TreeCommandBuilder) (argo.TreeCommand, error) {
	options := FixOptions(builder.Options())
	comOpts := WrapOptions(&options)

	errs := xerr.NewMultiError()

	if !builder.HasSubcommands() {
		errs.AppendError(errors.New("command tree has no subcommands"))
	}

	tree := &Tree{
		disableHelp: builder.IsHelpDisabled(),
		description: builder.Description(),
		callback:    builder.Callback(),
		options:     options,
		incompleteFn: utils.IfElse(
			builder.HasIncompleteHandler(),
			builder.IncompleteHandler(),
			MakeDefaultOnIncompleteHandler[argo.TreeCommand](comOpts),
		),

		// flagGroups:    nil, // filled below
		// commandGroups: nil, // filled below

		// selectedChild: nil, // filled on parse
		// selectedLeaf:  nil, // filled on parse
	}

	common.TryAddHelpFlags(builder, MakeRenderTreeHelpCallback(tree, comOpts), comOpts)

	tree.flagGroups = flag.BuildGroups(builder.FlagGroups(true), errs)
	tree.commandGroups = BuildCommandGroups(builder.CommandGroups(true), comOpts, tree, errs)

	if len(errs.Errors()) > 0 {
		return nil, errs
	}

	return tree, nil
}
