package common

import (
	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func TryAddHelpFlags[T any](c flag.GroupContainerBuilder[T], fn argo.FlagCallback, options argo.Options) {
	var targetGroup argo.FlagGroupBuilder
	var groups []argo.FlagGroupBuilder
	hasShortH := false
	hasLongH := false

	if c.HasFlagGroups(true) {
		groups = c.FlagGroups(true)

		// If there is no default flag group, or the default flag group already
		// has too many flags, create a meta group.
		if !flag.IsDefaultGroup(groups[0]) || groups[0].Size() > options.MaxDefaultFlagGroupSizeForMetaGroup {
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
	for _, g := range groups {
		for _, f := range g.Flags() {
			if f.HasShortForm() && f.ShortForm() == 'h' {
				hasShortH = true
			}

			if f.HasLongForm() && f.LongForm() == "help" {
				hasLongH = true
			}

			if hasShortH && hasLongH {
				break OUTER
			}
		}
	}

	if !hasLongH || !hasShortH {
		f := flag.NewBuilder().
			WithDescription("Prints this help text.").
			WithCallback(fn)

		if !hasLongH {
			f.WithLongForm("help")
		}

		if !hasShortH {
			f.WithShortForm('h')
		}

		targetGroup.WithFlag(f)
	}
}
