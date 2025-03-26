package common

import (
	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func TryAddHelpFlags[T any](c flag.GroupContainerBuilder[T], fn argo.FlagCallback, options argo.Options) {
	var targetGroup argo.FlagGroupBuilder
	var groups []argo.FlagGroupBuilder

	hasH := false
	hasHelp := false

	if c.HasFlagGroups(true) {
		groups = c.FlagGroups(true)

		// If there is no default flag group, or the default flag group already
		// has too many flags, create a meta group.
		if !flag.IsDefaultGroup(groups[0]) || groups[0].Size() > options.MaxDefaultFlagGroupSizeForMetaGroup {
			targetGroup = flag.NewDefaultGroupBuilder()
			c.WithFlagGroup(targetGroup)
		} else {
			targetGroup = groups[0]
		}

		if hasH, hasHelp = checkForHelpFlagConflicts(groups); hasH && hasHelp {
			return
		}
	} else {
		targetGroup = flag.NewDefaultGroupBuilder()
		c.WithFlagGroup(targetGroup)
	}

	f := flag.NewBuilder().
		WithDescription("Prints this help text.").
		WithCallback(fn)

	if !hasHelp {
		f.WithLongForm("help")
	}

	if !hasH {
		f.WithShortForm('h')
	}

	targetGroup.WithFlag(f)
}

func checkForHelpFlagConflicts(groups []argo.FlagGroupBuilder) (hasH, hasHelp bool) {
	for _, g := range groups {
		for _, f := range g.Flags() {
			if f.HasShortForm() && f.ShortForm() == 'h' {
				hasH = true
			}

			if f.HasLongForm() && f.LongForm() == "help" {
				hasHelp = true
			}

			if hasH && hasHelp {
				return
			}
		}
	}

	return
}
