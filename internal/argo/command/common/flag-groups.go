package common

import (
	"os"

	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/internal/render"
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func EnsureDefaultFlagGroup(groups []argo.FlagGroupBuilder) []argo.FlagGroupBuilder {
	if len(groups) == 0 {
		return append(make([]argo.FlagGroupBuilder, 0, 2), flag.NewDefaultGroupBuilder())
	}

	if flag.IsDefaultGroup(groups[0]) {
		return append(append(make([]argo.FlagGroupBuilder, 0, len(groups)+2), flag.NewDefaultGroupBuilder()), groups...)
	}

	return groups
}

type FlagGroupContainer interface {
	HasFlagGroups(includeDefault bool) bool
	FlagGroups(includeDefault bool) []argo.FlagGroupBuilder
}

func ConfigureHelpFlags[T any](builder FlagGroupContainer, renderer render.HelpRenderer[T], com T) []argo.FlagGroupBuilder {
	var group argo.FlagGroupBuilder
	var groups []argo.FlagGroupBuilder

	useLongH := true
	useShortH := true

	if !builder.HasFlagGroups(true) {
		group = flag.NewGroupBuilder(flag.DefaultFlagGroupName)
		groups = append(groups, group)
	} else {
		groups = builder.FlagGroups(true)

		if len(groups) > 1 || groups[0].Size() > 5 {
			group = flag.NewGroupBuilder("Help Flags")
			groups = append(groups, group)
		} else {
			group = groups[0]
		}

		for _, group := range groups {
			for _, flag := range group.Flags() {
				if flag.ShortForm() == 'h' {
					useShortH = false
				}
				if flag.LongForm() == "help" {
					useLongH = false
				}
				if !(useShortH || useLongH) {
					break
				}
			}
		}
	}

	if useShortH || useLongH {
		group.WithFlag(makeHelpFlag(useShortH, useLongH, renderer, com))
	}

	return groups
}

func makeHelpFlag[T any](short, long bool, renderer render.HelpRenderer[T], com T) argo.FlagBuilder {
	out := flag.NewBuilder().
		MarkAsHelpFlag().
		WithCallback(func(flag argo.Flag) {
			utils.Must(renderer.RenderHelp(com, os.Stdout))
			os.Exit(0)
		}).
		WithDescription("Prints this help text.")

	if short {
		out.WithShortForm('h')
	}

	if long {
		out.WithLongForm("help")
	}

	return out
}
