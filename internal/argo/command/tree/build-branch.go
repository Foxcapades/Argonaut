package tree

import (
	"errors"
	"os"

	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/internal/chars"
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/internal/xerr"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

// MaxDefaultFlagGroupSizeForMetaGroup defines the maximum number of flags that
// a command's default flag group may have before automatic flags such as the
// help flags are split off into their own group.
const MaxDefaultFlagGroupSizeForMetaGroup = 5

func BuildBranch(branch argo.BranchCommandBuilder, parent argo.ParentNode) (argo.BranchCommand, error) {
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
		tryAddHelpFlags(branch)
	}

	out.flagGroups = processFlagGroups(branch.FlagGroups(true), errs)
	out.commandGroups = processCommandGroups(branch.CommandGroups(), errs)

	if len(errs.Errors()) > 0 {
		return nil, errs
	}

	out.aliases = branch.Aliases()
	out.callback = branch.Callback()
	out.incompleteFn = branch.IncompleteHandler()

	return out, nil
}

func tryAddHelpFlags(c argo.BranchCommandBuilder) {
	var targetGroup argo.FlagGroupBuilder
	var groups []argo.FlagGroupBuilder
	hasShortH := false
	hasLongH := false

	if c.HasFlagGroups(true) {
		groups = c.FlagGroups(true)

		// If there is no default flag group, or the default flag group already
		// has too many flags, create a meta group.
		if !flag.IsDefaultGroup(groups[0]) || groups[0].Size() > MaxDefaultFlagGroupSizeForMetaGroup {
			targetGroup = flag.NewGroupBuilder(flag.MetaFlagGroupName)
			c.WithFlagGroup(targetGroup)
		} else {
			targetGroup = groups[0]
		}
	} else {
		targetGroup = flag.NewDefaultGroupBuilder()
		c.WithFlagGroup(targetGroup)
	}

OUTER:
	for _, group := range groups {
		for _, flag := range group.Flags() {
			if flag.HasShortForm() && flag.ShortForm() == 'h' {
				hasShortH = true
			}

			if flag.HasLongForm() && flag.LongForm() == "help" {
				hasLongH = true
			}

			if hasShortH && hasLongH {
				break OUTER
			}
		}
	}

	if !hasLongH || !hasShortH {
		flag := flag.NewBuilder().
			WithDescription("Prints this help text.").
			WithCallback(func(f argo.Flag) {
				utils.Must(comBranchRenderer{}.RenderHelp(out, os.Stdout))
				os.Exit(0)
			})

		if !hasLongH {
			flag.WithLongForm("help")
		}

		if !hasShortH {
			flag.WithShortForm('h')
		}

		targetGroup.WithFlag(flag)
	}
}

func processFlagGroups(builders []argo.FlagGroupBuilder, errs argo.MultiError) []argo.FlagGroup {
	flagGroups := make([]argo.FlagGroup, 0, len(builders))

	flag.UniqueFlagNames(builders, errs)

	for _, builder := range builders {
		if builder.HasFlags() {
			if group, err := builder.Build(ctx); err != nil {
				errs.AppendError(err)
			} else {
				flagGroups = append(flagGroups, group)
			}
		}
	}

	return flagGroups
}

func processCommandGroups(builders []argo.CommandGroupBuilder, errs argo.MultiError) []argo.CommandGroup {
	commandGroups := make([]argo.CommandGroup, 0, len(builders))

	massUniqueCommandNames(builders, errs)

	for _, build := range builders {
		if build.HasSubcommands() {
			if group, err := build.Build(ctx); err != nil {
				errs.AppendError(err)
			} else {
				commandGroups = append(commandGroups, group)
			}
		}
	}

	return commandGroups
}
