package tree

import (
	"errors"

	"github.com/foxcapades/argonaut/v3/internal/argo/command/common"
	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/internal/text"
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/internal/xerr"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func BuildBranch(
	branch argo.BranchCommandBuilder,
	options argo.Options,
	parent argo.ParentNode,
) (argo.BranchCommand, error) {
	errs := xerr.NewMultiError()

	// Ensure name is not blank
	xerr.AppendIfPresent(ValidateNodeName(branch.Name()), errs)

	// Ensure aliases are not blank
	for _, alias := range branch.Aliases() {
		if text.IsBlank(alias) {
			errs.AppendError(errors.New("command branch aliases must not be blank"))
		}
	}

	// Create the out instance ahead of time so that we can set it as the parent
	// on the command groups we build.
	out := &Branch{
		name:        branch.Name(),
		description: branch.Description(),
		parent:      parent,
	}

	// If auto-help is not disabled, then...
	if !branch.IsHelpDisabled() {
		common.TryAddHelpFlags(branch, MakeRenderBranchHelpCallback(out, options), options)
	}

	out.flagGroups = flag.BuildGroups(branch.FlagGroups(true), options, errs)
	out.commandGroups = BuildCommandGroups(branch.CommandGroups(true), options, out, errs)

	if len(errs.Errors()) > 0 {
		return nil, errs
	}

	out.aliases = branch.Aliases()
	out.callback = branch.Callback()
	out.incompleteFn = utils.IfElse(
		branch.HasIncompleteHandler(),
		branch.IncompleteHandler(),
		MakeDefaultOnIncompleteHandler[argo.BranchCommand](options),
	)

	return out, nil
}
