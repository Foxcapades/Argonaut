package tree

import (
	"errors"

	"github.com/foxcapades/argonaut/v3/internal/argo/command/common"
	"github.com/foxcapades/argonaut/v3/internal/chars"
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/internal/xerr"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func BuildBranch(branch argo.BranchCommandBuilder, options argo.Options, parent argo.ParentNode) (argo.BranchCommand, error) {
	errs := xerr.NewMultiError()

	// Ensure name is not blank
	if err := chars.ValidateCommandNodeName(branch.Name()); err != nil {
		errs.AppendError(err)
	}

	// Ensure aliases are not blank
	for _, alias := range branch.Aliases() {
		if chars.IsBlank(alias) {
			errs.AppendError(errors.New("command branch aliases must not be blank"))
		}
	}

	// Ensure a parent is set
	if branch.ParentNode() == nil {
		panic("illegal state: attempted to build a command branch with no parent set")
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

	out.flagGroups = processFlagGroups(branch.FlagGroups(true), errs)
	out.commandGroups = processCommandGroups(branch.CommandGroups(true), options, out, errs)

	if len(errs.Errors()) > 0 {
		return nil, errs
	}

	out.aliases = branch.Aliases()
	out.callback = branch.Callback()
	out.incompleteFn = utils.IfElse(
		branch.HasIncompleteHandler(),
		branch.IncompleteHandler(),
		makeDefaultOnIncompleteHandler[argo.BranchCommand](options),
	)

	return out, nil
}
